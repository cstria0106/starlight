/*
   file created by Junlin Chen in 2022

*/

package ctr

import "github.com/containerd/containerd"

func NewContainerd(namespace, socket string) (client *containerd.Client, err error) {
	if client, err = containerd.New(socket, containerd.WithDefaultNamespace(namespace)); err != nil {
		return nil, err
	}

	return client, nil
}
