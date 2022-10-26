/*
   file created by Junlin Chen in 2022

*/

package ctr

import (
	"context"
	"fmt"
	"github.com/containerd/containerd/images"
	"github.com/google/uuid"
	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	ctx := context.Background()

	client, err := NewContainerd("default", "/run/containerd/containerd.sock", "debug")
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Close()

	//cs := client.ContentStore()

	is := client.ImageService()
	var img images.Image
	img, err = is.Create(ctx, images.Image{
		Name: uuid.New().String()[:10],
		Labels: map[string]string{
			"test": "test",
		},
		Target: ocispec.Descriptor{
			MediaType:   "application/vnd.docker.distribution.manifest.v2+json",
			Digest:      digest.NewDigestFromHex("sha256", "e591efb0ce4780263ed1abbe32b905b63ed483ce34172cacacd18031a281fb28"),
			Size:        1000,
			URLs:        []string{"yuri.moe/starlight"},
			Annotations: nil,
			Platform:    nil,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(img)
}
