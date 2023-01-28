/*
   file created by Junlin Chen in 2022

*/

package receive

import (
	"fmt"
	"path/filepath"
	"syscall"
	"unsafe"

	fuseFs "github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/mc256/starlight/client/fs"
	"github.com/mc256/starlight/util/common"
	"golang.org/x/sys/unix"
)

type ImageLayer struct {
	Size   int64 `json:"s"`
	Serial int64 `json:"f"`

	// Hash is the digest of the compressed layer
	Hash string `json:"h"`

	// path to the local storage
	Local string
}

func (il ImageLayer) String() string {
	return fmt.Sprintf("[%05d:%02d]%s-%d", il.Serial, -1, il.Hash, il.Size)
}

type Content struct {
	Signal chan interface{} `json:"-"`

	// ------------------------------------------
	// stack identify which layer should this content be placed, all the files will be referencing the content
	Stack int64 `json:"t"`

	// offset is non-zero if the file is in the delta bundle body
	Offset int64 `json:"o,omitempty"`

	// size is the size of the compressed content
	Size int64 `json:"s"`

	Chunks []*FileChunk `json:"c"`

	Digest string `json:"d"`
}

func (c *Content) GetBaseDir() string {
	return filepath.Join(c.Digest[7:8], c.Digest[8:10], c.Digest[10:12])
}

func (c *Content) GetPath() string {
	return filepath.Join(c.GetBaseDir(), c.Digest[12:])
}

type ReferencedFile struct {
	File

	// Stack in the existing image from bottom to top, 0-indexed
	Stack int64 `json:"S"`

	// if the file is available on the client then ReferenceFsId is non-zero,
	// expecting the file is available on the client and can be accessed using the File.Digest .
	// (This is the Serial index in the database `filesystem.id`)
	ReferenceFsId int64 `json:"R,omitempty"`

	// if the file is not available on the client but on other layers in the requested image，
	// then ReferenceFsId is zero and ReferenceStack is non-zero,
	// expecting the file content in the delta bundle body
	// (This is Stack not Serial, 0-indexed)
	ReferenceStack int64 `json:"T,omitempty"`

	// PayloadOrder is set if the file is not available on the client, then PayloadOrder is non-zero.
	// It indicates the order of the file in the delta bundle body (payload)
	PayloadOrder int `json:"O,omitempty"`

	// if Ready is nil or closed, means the file is ready
	Ready *chan interface{} `json:"-"`

	// if Waiting is not nil and closed, means the file is waited
	Waiting *chan struct{} `json:"-"`

	stable   fuseFs.StableAttr
	children []fs.ReceivedFile
}

// ------------------------------------------
// use in file system
//

func (r *ReferencedFile) GetChildren() []fs.ReceivedFile {
	return r.children
}

func (r *ReferencedFile) AppendChild(children fs.ReceivedFile) {
	if r.children == nil {
		r.children = make([]fs.ReceivedFile, 0)
	}
	r.children = append(r.children, children)
}

func (r *ReferencedFile) IsReady() bool {
	return r.Ready == nil
}

func (r *ReferencedFile) InitFuseStableAttr() {
	r.stable.Ino = uint64(uintptr(unsafe.Pointer(r)))
	r.stable.Gen = 0
	r.stable.Mode = modeOfEntry(r)
}

func (r *ReferencedFile) GetAttr(out *fuse.Attr) syscall.Errno {
	out.Ino = r.stable.Ino
	out.Size = uint64(r.Size)
	if r.IsDir() {
		out.Size = 4096
	} else if r.Type == "symlink" {
		out.Size = uint64(len(r.LinkName))
	}
	r.SetBlockSize(out)
	mtime := r.ModTime()
	out.SetTimes(&mtime, &mtime, &mtime)
	out.Mode = r.stable.Mode
	out.Owner = fuse.Owner{Uid: uint32(r.UID), Gid: uint32(r.GID)}
	out.Rdev = uint32(unix.Mkdev(uint32(r.DevMajor), uint32(r.DevMinor)))
	out.Nlink = uint32(r.NumLink)
	if out.Nlink == 0 {
		out.Nlink = 1 // zero "NumLink" means one.
	}
	return 0
}

func (r *ReferencedFile) GetXAttrs() map[string][]byte {
	if r.Xattrs == nil {
		return make(map[string][]byte)
	}
	return r.Xattrs
}

func (r *ReferencedFile) GetName() string {
	return r.Name
}

func (r *ReferencedFile) GetStableAttr() *fuseFs.StableAttr {
	return &r.stable
}

func (r *ReferencedFile) GetLinkName() string {
	return r.LinkName
}

func (r *ReferencedFile) GetBaseDir() string {
	return filepath.Join(r.Digest[7:8], r.Digest[8:10], r.Digest[10:12])
}

func (r *ReferencedFile) GetRealPath() string {
	return filepath.Join(r.GetBaseDir(), r.Digest[12:])
}

func (r *ReferencedFile) WaitForReady() {
	if r.Waiting != nil {
		close(*r.Waiting)
		r.Waiting = nil
	}
	<-*r.Ready
}

func (r *ReferencedFile) IsReferencingRequestedImage() (stack int64, yes bool) {
	if r.ReferenceFsId != 0 {
		return 0, false
	}

	// in the payload
	if r.ReferenceStack != 0 {
		// different layer
		return r.ReferenceStack, true
	}
	// same layer
	return r.Stack, true
}

func (r *ReferencedFile) IsReferencingLocalFilesystem() (serial int64, yes bool) {
	if r.ReferenceFsId != 0 {
		return r.ReferenceFsId, true
	}
	return 0, false
}

// InPayload returns true if the content of the file is in the delta bundle (payload)
func (r *ReferencedFile) InPayload() bool {
	return r.PayloadOrder > 0
}

type FileChunk struct {
	Offset         int64 `json:"o"`
	ChunkOffset    int64 `json:"c"`
	ChunkSize      int64 `json:"h"`
	CompressedSize int64 `json:"s"`
}

type File struct {
	common.TOCEntry
	Chunks []*FileChunk `json:"c,omitempty"`
	FsId   int64        `json:"-"`
}

type Image struct {
	Serial int64         `json:"s"`
	Layers []*ImageLayer `json:"l"`
}

func (i Image) String() string {
	return fmt.Sprintf("%d->%v", i.Serial, i.Layers)
}

type DeltaBundle struct {
	Source      *Image `json:"s"`
	Destination *Image `json:"d"`

	// contents and BodyLength are computed by Builder.computeDelta()
	Contents   []*Content `json:"c"`
	BodyLength int64      `json:"bl"`

	// RequestedFiles are all the files in the requested images
	// Use this to reconstruct the file system
	RequestedFiles []*ReferencedFile `json:"rf"`
}
