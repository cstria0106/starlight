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

package main

import (
	"context"
	"fmt"
	"github.com/containerd/containerd/log"
	"github.com/mc256/starlight/grpc"
	"github.com/mc256/starlight/util"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

func init() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println(c.App.Name, c.App.Version)
	}
}

func main() {
	app := New()
	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "starlight-grpc: \n%v\n", err)
		os.Exit(1)
	}
}

func New() *cli.App {
	app := cli.NewApp()
	cfg := grpc.LoadConfig(context.TODO())

	app.Name = "starlight-grpc"
	app.Version = util.Version
	app.Usage = `gRPC snapshotter plugin for faster container-based application deployment`
	app.Description = fmt.Sprintf("\n%s\n", app.Usage)

	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			DefaultText: "/etc/starlight/snapshotter_config.json",
			Aliases:     []string{"c"},
			EnvVars:     []string{"STARLIGHT_GRPC_CONFIG"},
			Usage:       "json configuration file. CLI parameter will override values in the config file if specified",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "log-level",
			DefaultText: cfg.LogLevel,
			Usage:       "Choose one log level (fatal, error, warning, info, debug, trace)",
			Required:    false,
		},
		// ----
		&cli.StringFlag{
			Name:        "metadata",
			DefaultText: cfg.Metadata,
			Aliases:     []string{"m"},
			Usage:       "path to store image metadata",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "socket",
			DefaultText: cfg.Socket,
			Usage:       "gRPC socket address",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "default",
			DefaultText: cfg.DefaultProxy,
			Aliases:     []string{"d"},
			Usage:       "name of the default proxy",
		},
		&cli.StringFlag{
			Name:        "fs-root",
			DefaultText: cfg.FileSystemRoot,
			Aliases:     []string{"fs"},
			Usage:       "path to store uncompress image layers",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "id",
			DefaultText: cfg.ClientId,
			Usage:       "identifier for the client",
			Required:    false,
		},
		// ----
		&cli.StringSliceFlag{
			Name:        "proxy",
			Aliases:     []string{"p"},
			Usage:       "proxy of the configuration use comma (',') to separate components, and use another tag for other proxies, name,protocol,address,username,password",
			Required:    false,
			DefaultText: "starlight-shared,https,starlight.yuri.moe,,",
		},
	}
	app.Action = func(c *cli.Context) error {
		return DefaultAction(c, cfg)
	}

	return app
}

func DefaultAction(context *cli.Context, cfg *grpc.Configuration) error {
	if context.Bool("version") == true {
		fmt.Printf("starlight-proxy v%s\n", util.Version)
		return nil
	}

	if l := context.String("log-level"); l != "" {
		cfg.LogLevel = l
	}
	c := util.ConfigLoggerWithLevel(cfg.LogLevel)

	log.G(c).WithFields(logrus.Fields{
		"version":   util.Version,
		"log-level": context.String("log-level"),
	}).Info("Starlight Snapshotter")

	grpc.NewSnapshotterGrpcService(c, cfg)

	return nil
}
