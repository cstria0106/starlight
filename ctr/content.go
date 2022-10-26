/*
   file created by Junlin Chen in 2022

*/

package ctr

import (
	"context"
	"fmt"
	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/remotes"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

func Fetch(ctx context.Context, ingester content.Ingester, fetcher remotes.Fetcher, desc ocispec.Descriptor) error {
	log.G(ctx).Debug("fetch")

	cw, err := content.OpenWriter(ctx, ingester, content.WithRef(remotes.MakeRefKey(ctx, desc)), content.WithDescriptor(desc))
	if err != nil {
		if errdefs.IsAlreadyExists(err) {
			return nil
		}
		return err
	}
	defer cw.Close()

	ws, err := cw.Status()
	if err != nil {
		return err
	}

	if desc.Size == 0 {
		// most likely a poorly configured registry/web front end which responded with no
		// Content-Length header; unable (not to mention useless) to commit a 0-length entry
		// into the content store. Error out here otherwise the error sent back is confusing
		return fmt.Errorf("unable to fetch descriptor (%s) which reports content size of zero: %w", desc.Digest, errdefs.ErrInvalidArgument)
	}
	if ws.Offset == desc.Size {
		// If writer is already complete, commit and return
		err := cw.Commit(ctx, desc.Size, desc.Digest)
		if err != nil && !errdefs.IsAlreadyExists(err) {
			return fmt.Errorf("failed commit on ref %q: %w", ws.Ref, err)
		}
		return nil
	}

	rc, err := fetcher.Fetch(ctx, desc)
	if err != nil {
		return err
	}
	defer rc.Close()

	return content.Copy(ctx, cw, rc, desc.Size, desc.Digest)
}
