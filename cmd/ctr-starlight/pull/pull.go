/*
   file created by Junlin Chen in 2022

*/

package pull

import (
	"errors"
	"fmt"
	"github.com/mc256/starlight/cmd/ctr-starlight/auth"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	var ref string
	if c.Args().Len() == 1 {
		ref = c.Args().First()
	} else {
		return errors.New("wrong arguments")
	}

	ns := c.String("namespace")
	if ns == "" {
		ns = "default"
	}

	socket := c.String("address")
	if socket == "" {
		socket = "/run/containerd/containerd.sock"
	}

	//from := c.String("from")
	fmt.Println(ref)
	fmt.Println(c.String("server"))
	// Prepare containerd

	// Check available images

	// Fetch Starlight image

	return nil
}

func Command() *cli.Command {
	cmd := cli.Command{
		Name:  "pull",
		Usage: "Launch background fetcher to load the delta image",
		Action: func(c *cli.Context) error {
			return Action(c)
		},
		Flags:     append(Flags, auth.StarlightProxyFlags...),
		ArgsUsage: "StarlightImage",
	}
	return &cmd
}
