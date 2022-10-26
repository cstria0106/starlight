/*
   file created by Junlin Chen in 2022

*/

package pull

import (
	"github.com/urfave/cli/v2"
)

var (
	Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "from",
			Usage:    "specify a particular container image that the (if not specified, the latest downloaded container image with the same 'image name' will be used)",
			Value:    "",
			Required: false,
		},
	}
)
