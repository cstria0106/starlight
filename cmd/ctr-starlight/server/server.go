/*
   file created by Junlin Chen in 2022

*/

package server

import (
	"context"
	"github.com/urfave/cli/v2"
)

func Action(ctx context.Context, c *cli.Context) (err error) {
	return nil
}

func Command() *cli.Command {
	ctx := context.Background()
	return &cli.Command{
		Name:  "server",
		Usage: "launch an API server for starlight",
		Action: func(c *cli.Context) error {
			return Action(ctx, c)
		},
		Flags:     []cli.Flag{},
		ArgsUsage: "",
	}
}
