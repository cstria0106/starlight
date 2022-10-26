/*
   file created by Junlin Chen in 2022

*/

package grpc

import (
	"context"
	"github.com/containerd/containerd/api/services/content/v1"
	"github.com/containerd/containerd/log"
	"github.com/gogo/protobuf/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EmptyContent struct {
}

// Info returns information about a committed object.
//
// This call can be used for getting the size of content and checking for
// existence.
func (c *EmptyContent) Info(ctx context.Context, req *content.InfoRequest) (*content.InfoResponse, error) {
	log.G(ctx).Debugf("Info")
	return nil, status.Errorf(codes.Unimplemented, "method Info not implemented")
}

// Update updates content metadata.
//
// This call can be used to manage the mutable content labels. The
// immutable metadata such as digest, size, and committed at cannot
// be updated.
func (c *EmptyContent) Update(ctx context.Context, req *content.UpdateRequest) (*content.UpdateResponse, error) {
	log.G(ctx).Debugf("Update")
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}

// List streams the entire set of content as Info objects and closes the
// stream.
//
// Typically, this will yield a large response, chunked into messages.
// Clients should make provisions to ensure they can handle the entire data
// set.
func (c *EmptyContent) List(listReq *content.ListContentRequest, cls content.Content_ListServer) error {
	log.G(cls.Context()).Debugf("List")
	return status.Errorf(codes.Unimplemented, "method List not implemented")
}

// Delete will delete the referenced object.
func (c *EmptyContent) Delete(ctx context.Context, req *content.DeleteContentRequest) (*types.Empty, error) {
	log.G(ctx).Debugf("Delete")
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}

// Read allows one to read an object based on the offset into the content.
//
// The requested data may be returned in one or more messages.
func (c *EmptyContent) Read(req *content.ReadContentRequest, reqs content.Content_ReadServer) error {
	log.G(reqs.Context()).Debugf("Read")
	return status.Errorf(codes.Unimplemented, "method Read not implemented")
}

// Status returns the status for a single reference.
func (c *EmptyContent) Status(ctx context.Context, req *content.StatusRequest) (*content.StatusResponse, error) {
	log.G(ctx).Debugf("Status")
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}

// ListStatuses returns the status of ongoing object ingestions, started via
// Write.
//
// Only those matching the regular expression will be provided in the
// response. If the provided regular expression is empty, all ingestions
// will be provided.
func (c *EmptyContent) ListStatuses(ctx context.Context, req *content.ListStatusesRequest) (*content.ListStatusesResponse, error) {
	log.G(ctx).Debugf("ListStatuses")
	return nil, status.Errorf(codes.Unimplemented, "method ListStatuses not implemented")
}

// Write begins or resumes writes to a resource identified by a unique ref.
// Only one active stream may exist at a time for each ref.
//
// Once a write stream has started, it may only write to a single ref, thus
// once a stream is started, the ref may be omitted on subsequent writes.
//
// For any write transaction represented by a ref, only a single write may
// be made to a given offset. If overlapping writes occur, it is an error.
// Writes should be sequential and implementations may throw an error if
// this is required.
//
// If expected_digest is set and already part of the content store, the
// write will fail.
//
// When completed, the commit flag should be set to true. If expected size
// or digest is set, the content will be validated against those values.
func (c *EmptyContent) Write(cws content.Content_WriteServer) error {
	log.G(cws.Context()).Debugf("Write")
	return status.Errorf(codes.Unimplemented, "method Write not implemented")
}

// Abort cancels the ongoing write named in the request. Any resources
// associated with the write will be collected.
func (c *EmptyContent) Abort(ctx context.Context, req *content.AbortRequest) (*types.Empty, error) {
	log.G(ctx).Debugf("Abort")
	return nil, status.Errorf(codes.Unimplemented, "method Abort not implemented")
}

func NewEmptyContent() content.ContentServer {
	return &EmptyContent{}
}
