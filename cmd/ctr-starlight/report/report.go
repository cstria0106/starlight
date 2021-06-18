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

package report

import (
	"context"
	"github.com/mc256/starlight/fs"
	"github.com/mc256/starlight/grpc"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	ctx := context.Background()
	tc, err := fs.NewTraceCollection(ctx, c.String("path"))
	if err != nil {
		return err
	}
	protocol := "https"
	if c.Bool("plain-http") {
		protocol = "http"
	}
	proxy := grpc.NewStarlightProxy(ctx, protocol, c.String("server"))
	if err := proxy.Report(tc.ToJsonBuffer()); err != nil {
		return err
	}

	return nil
}

func Command() *cli.Command {
	cmd := cli.Command{
		Name:  "report",
		Usage: "Upload data collected by the optimizer back to Starlight Proxy to speed up other similar deployment",
		Action: func(c *cli.Context) error {
			return Action(c)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "path",
				Usage:       "path the the optimizer logs",
				Value:       "/tmp",
				DefaultText: "/tmp",
				Required:    false,
			},
			&cli.StringFlag{
				Name:     "server",
				Value:    "worker1",
				Usage:    "starlight proxy address",
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "plain-http",
				Usage:    "use plain http connects to the remote server",
				Required: false,
			},
		},
		ArgsUsage: "",
	}
	return &cmd
}
