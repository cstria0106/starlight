/*
   Copyright The starlight Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

   file created by maverick in 2021
*/

package proxy

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/containerd/containerd/log"
	"github.com/gorilla/websocket"
	"github.com/mc256/starlight/fs"
	"github.com/mc256/starlight/merger"
	"github.com/mc256/starlight/util"
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

type transition struct {
	tagFrom []*util.ImageRef
	tagTo   []*util.ImageRef
}

type StarlightProxyServer struct {
	http.Server

	ctx      context.Context
	database *bolt.DB

	containerRegistry string

	builder *DeltaBundleBuilder
}

func (a *StarlightProxyServer) getDeltaImage(w http.ResponseWriter, req *http.Request, from string, to string) error {
	// Parse Image Reference
	t := &transition{
		tagFrom: make([]*util.ImageRef, 0),
		tagTo:   make([]*util.ImageRef, 0),
	}

	var err error
	if from != "_" {
		if t.tagFrom, err = util.NewImageRef(from); err != nil {
			return err
		}
	}
	if t.tagTo, err = util.NewImageRef(to); err != nil {
		return err
	}

	// Load Optimized Merged Image Collections
	var cTo, cFrom *Collection

	if cTo, err = LoadCollection(a.ctx, a.database, t.tagTo); err != nil {
		return err
	}
	if len(t.tagFrom) != 0 {
		if cFrom, err = LoadCollection(a.ctx, a.database, t.tagFrom); err != nil {
			return err
		}
		cTo.Minus(cFrom)
	}

	deltaBundle := cTo.ComposeDeltaBundle()

	buf := bytes.NewBuffer(make([]byte, 0))
	wg := &sync.WaitGroup{}

	// build header
	headerSize, contentLength, err := a.builder.WriteHeader(buf, deltaBundle, wg, false)
	if err != nil {
		log.G(a.ctx).WithField("err", err).Error("write header cache")
		return nil
	}

	header := http.Header{}
	header.Set("Starlight-Header-Size", fmt.Sprintf("%d", headerSize))
	header.Set("Starlight-Payload-Size", fmt.Sprintf("%d", contentLength))
	header.Set("Starlight-Version", util.Version)

	upgrader := websocket.Upgrader{WriteBufferSize: 1024, ReadBufferSize: 1024}
	conn, err := upgrader.Upgrade(w, req, header)
	if err != nil {
		log.G(a.ctx).WithField("err", err).Error("http upgrade error")
		return nil
	}

	defer conn.Close()

	var wc io.WriteCloser

	if wc, err = conn.NextWriter(websocket.BinaryMessage); err != nil {
		log.G(a.ctx).WithField("err", err).Error("create writer error")
	}

	defer wc.Close()

	// write header
	if n, err := io.CopyN(wc, buf, headerSize); err != nil || n != headerSize {
		log.G(a.ctx).WithField("err", err).Error("write header error")
		return nil
	}

	// write payload
	fileRequests := &FileRequests{}

	go func() {
		for {
			t, b, err := conn.ReadMessage()
			if err != nil {
				return
			}
			if t != websocket.BinaryMessage || len(b) != 12 {
				continue
			}
			fileRequests.Push(
				int(int32(binary.BigEndian.Uint32(b[:4]))),
				int64(binary.BigEndian.Uint64(b[4:12])),
				int64(binary.BigEndian.Uint32(b[12:])),
			)
		}
	}()

	if err = a.builder.WriteBody(fileRequests, wc, deltaBundle, wg); err != nil {
		log.G(a.ctx).WithField("err", err).Error("write body error")
		return nil
	}

	return nil
}

func (a *StarlightProxyServer) getPrepared(w http.ResponseWriter, req *http.Request, image string) error {
	arr := strings.Split(strings.Trim(image, ""), ":")
	if len(arr) != 2 || arr[0] == "" || arr[1] == "" {
		return util.ErrWrongImageFormat
	}

	err := CacheToc(a.ctx, a.database, arr[0], arr[1], a.containerRegistry)
	if err != nil {
		return err
	}

	ob := merger.NewOverlayBuilder(a.ctx, a.database)
	if err = ob.AddImage(arr[0], arr[1]); err != nil {
		return err
	}
	if err = ob.SaveMergedImage(); err != nil {
		return err
	}

	header := w.Header()
	header.Set("Content-Type", "text/plain")
	header.Set("Starlight-Version", util.Version)
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "Cached TOC: %s\n", image)
	return nil
}

func (a *StarlightProxyServer) postReport(w http.ResponseWriter, req *http.Request) error {
	header := w.Header()
	header.Set("Content-Type", "text/plain")
	header.Set("Starlight-Version", util.Version)
	w.WriteHeader(http.StatusOK)

	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	tc, err := fs.NewTraceCollectionFromBuffer(buf)
	if err != nil {
		log.G(a.ctx).WithError(err).Info("cannot parse trace collection")
		return err
	}

	for _, grp := range tc.Groups {
		log.G(a.ctx).WithField("collection", grp.Images)
		fso, err := LoadCollection(a.ctx, a.database, grp.Images)
		if err != nil {
			return err
		}

		fso.AddOptimizeTrace(grp)

		if err := fso.SaveMergedApp(); err != nil {
			return err
		}
		_, _ = fmt.Fprintf(w, "Optimized: %s \n", grp.Images)
	}

	return nil
}

func (a *StarlightProxyServer) getDefault(w http.ResponseWriter, req *http.Request) {
	_, _ = fmt.Fprint(w, "Starlight Proxy OK!\n")
}

func (a *StarlightProxyServer) rootFunc(w http.ResponseWriter, req *http.Request) {
	params := strings.Split(strings.Trim(req.RequestURI, "/"), "/")
	remoteAddr := req.RemoteAddr

	if realIp := req.Header.Get("X-Real-IP"); realIp != "" {
		remoteAddr = realIp
	}
	log.G(a.ctx).WithFields(logrus.Fields{
		"remote": remoteAddr,
		"params": params,
	}).Info("request received")
	var err error
	switch {
	case len(params) == 4 && params[0] == "from" && params[2] == "to":
		err = a.getDeltaImage(w, req, strings.TrimSpace(params[1]), strings.TrimSpace(params[3]))
		break
	case len(params) == 2 && params[0] == "prepare":
		err = a.getPrepared(w, req, params[1])
		break
	case len(params) == 1 && params[0] == "report":
		err = a.postReport(w, req)
		break
	default:
		a.getDefault(w, req)
	}
	if err != nil {
		header := w.Header()
		header.Set("Content-Type", "text/plain")
		header.Set("Starlight-Version", util.Version)
		w.WriteHeader(http.StatusInternalServerError)

		_, _ = fmt.Fprintf(w, "Opoos! Something went wrong: \n\n%s\n", err)
	} else {
		log.G(a.ctx).WithFields(logrus.Fields{
			"remote": remoteAddr,
			"params": params,
		}).Info("request sent")
	}
}

func NewServer(registry, logLevel string, wg *sync.WaitGroup) *StarlightProxyServer {
	ctx := util.ConfigLoggerWithLevel(logLevel)

	log.G(ctx).WithFields(logrus.Fields{
		"registry":  registry,
		"log-level": logLevel,
	}).Info("Starlight Proxy")

	db, err := util.OpenDatabase(ctx, util.DataPath, util.ProxyDbName)
	if err != nil {
		log.G(ctx).WithError(err).Error("open database error")
		return nil
	}

	server := &StarlightProxyServer{
		Server: http.Server{
			Addr: ":8090",
		},
		database:          db,
		ctx:               ctx,
		containerRegistry: registry,
		builder:           NewBuilder(ctx, registry),
	}
	http.HandleFunc("/", server.rootFunc)

	go func() {
		defer wg.Done()
		defer server.database.Close()

		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.G(ctx).WithField("error", err).Error("server exit with error")
		}
	}()

	return server
}
