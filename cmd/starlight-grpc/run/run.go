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

package run

import (
	"github.com/containerd/containerd/log"
	sn "github.com/mc256/starlight/grpc"
	"github.com/mc256/starlight/util"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func RunAction(c *cli.Context) error {
	ctx := util.ConfigLoggerWithLevel(c.String("log-level"))
	log.G(ctx).WithFields(logrus.Fields{
		"version":   util.Version,
		"log-level": c.String("log-level"),
	}).Info("Starlight Snapshotter")

	protocol := "wss"
	if c.Bool("insecure") {
		protocol = "ws"
	}

	sn.NewSnapshotterGrpcService(
		ctx,
		c.String("socket"),
		protocol,
		c.String("server"),
		c.String("fs"),
		c.Bool("log-fs-trace"),
	)

	return nil
}

func RunCommand() *cli.Command {
	cmd := cli.Command{
		Name:  "run",
		Usage: "launch starlight gRPC snapshotter plugin",
		Action: func(c *cli.Context) error {
			return RunAction(c)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "socket",
				Value:    "/run/starlight-grpc/starlight-snapshotter.socket",
				Usage:    "gRPC socket address",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "server",
				Value:    "worker1",
				Usage:    "starlight proxy address",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "fs",
				Value:    "/var/lib/starlight-grpc",
				Usage:    "snapshotter file system path",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "log-level",
				Value:    "info",
				Usage:    "log level",
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "insecure",
				Usage:    "use plain ws connects to the remote server",
				Required: false,
			},
			&cli.BoolFlag{
				Name:        "log-fs-trace",
				Usage:       "collect file system traces",
				Value:       false,
				DefaultText: "false",
				Required:    false,
			},
		},
	}
	return &cmd
}
