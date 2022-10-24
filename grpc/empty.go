/*
   file created by Junlin Chen in 2022

*/

package grpc

import (
	"context"
	"fmt"
	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/mount"
	"github.com/containerd/containerd/snapshots"
	"github.com/containerd/containerd/snapshots/storage"
	"github.com/google/uuid"
	starlightfs "github.com/mc256/starlight/fs"
	"github.com/mc256/starlight/util"
	"go.etcd.io/bbolt"
	"os"
	"path/filepath"
	"time"
)

type emptyInstance struct {
	parent *emptyInstance
	id     string
	mount  *mount.Mount
}

func NewEmptyInstance(parent *emptyInstance, id string) *emptyInstance {
	return &emptyInstance{
		parent: parent,
		id:     id,
		mount: &mount.Mount{
			Type:    "bind",
			Source:  filepath.Join("/tmp", uuid.New().String()),
			Options: []string{"rbind"},
		},
	}
}

type emptySnapshotter struct {
	gCtx context.Context

	ms *storage.MetaStore
	db *bbolt.DB

	layerStore *starlightfs.LayerStore
	receiver   map[string]*starlightfs.Receiver

	//imageReadersMux sync.Mutex
	fsMap   map[string]*starlightfs.FsInstance
	fsTrace bool

	cfg              *Configuration
	proxyConnections map[string]*ProxyConfig

	// Test Mounting Map
	buffer map[string]*emptyInstance
}

func NewEmptySnapshotter(ctx context.Context, cfg *Configuration) (snapshots.Snapshotter, error) {
	if err := os.MkdirAll(cfg.FileSystemRoot, 0700); err != nil {
		return nil, err
	}

	// containerd snapshot database
	ms, err := storage.NewMetaStore(cfg.Metadata + ".sn")
	if err != nil {
		return nil, err
	}

	// starlight metadata database
	db, err := util.OpenDatabase(ctx, cfg.Metadata+".sl")
	if err != nil {
		return nil, err
	}

	// root path for starlight fs
	layerStore, err := starlightfs.NewLayerStore(ctx, db, filepath.Join(cfg.FileSystemRoot, "sfs"))
	if err != nil {
		return nil, err
	}

	if err := os.Mkdir(filepath.Join(cfg.FileSystemRoot, "sfs"), 0700); err != nil && !os.IsExist(err) {
		return nil, err
	}
	return &emptySnapshotter{
		gCtx: ctx,

		ms: ms,
		db: db,

		layerStore: layerStore,

		receiver: make(map[string]*starlightfs.Receiver, 0),
		fsMap:    make(map[string]*starlightfs.FsInstance, 0),

		buffer: make(map[string]*emptyInstance, 0),
	}, nil
}

// Stat returns the info for an active or committed snapshot by name or
// key.
//
// Should be used for parent resolution, existence checks and to discern
// the kind of snapshot.
func (s *emptySnapshotter) Stat(ctx context.Context, key string) (snapshots.Info, error) {
	log.G(ctx).WithField("key", key).Info("stat")
	ctx, t, err := s.ms.TransactionContext(ctx, false)
	if err != nil {
		return snapshots.Info{}, err
	}
	defer t.Rollback()
	_, info, _, err := storage.GetInfo(ctx, key)
	if err != nil {
		return snapshots.Info{}, err
	}
	return info, nil
}

// Update updates the info for a snapshot.
//
// Only mutable properties of a snapshot may be updated.
func (s *emptySnapshotter) Update(ctx context.Context, info snapshots.Info, fieldpaths ...string) (snapshots.Info, error) {
	log.G(ctx).WithField("info", info).WithField("fieldPaths", fieldpaths).Info("update")
	ctx, t, err := s.ms.TransactionContext(ctx, true)
	if err != nil {
		return snapshots.Info{}, err
	}
	defer t.Rollback()
	info, err = storage.UpdateInfo(ctx, info, fieldpaths...)
	if err != nil {
		return snapshots.Info{}, err
	}

	if err := t.Commit(); err != nil {
		return snapshots.Info{}, err
	}
	return info, nil
}

// Usage returns the resource usage of an active or committed snapshot
// excluding the usage of parent snapshots.
//
// The running time of this call for active snapshots is dependent on
// implementation, but may be proportional to the size of the resource.
// Callers should take this into consideration. Implementations should
// attempt to honer context cancellation and avoid taking locks when making
// the calculation.
func (s *emptySnapshotter) Usage(ctx context.Context, key string) (snapshots.Usage, error) {
	log.G(ctx).WithField("key", key).Info("usage")
	ctx, t, err := s.ms.TransactionContext(ctx, false)
	if err != nil {
		return snapshots.Usage{}, err
	}
	defer t.Rollback()

	id, info, usage, err := storage.GetInfo(ctx, key)
	if err != nil {
		return snapshots.Usage{}, err
	}
	log.G(ctx).WithField("key", key).WithField("id", id).Info("usage")
	if info.Kind == snapshots.KindActive {
		usage = snapshots.Usage{
			Size:   100,
			Inodes: 100,
		}
	}
	return usage, nil
}

// Mounts returns the mounts for the active snapshot transaction identified
// by key. Can be called on an read-write or readonly transaction. This is
// available only for active snapshots.
//
// This can be used to recover mounts after calling View or Prepare.
func (s *emptySnapshotter) Mounts(ctx context.Context, key string) ([]mount.Mount, error) {
	log.G(ctx).WithField("key", key).Info("mounts")
	m, has := s.buffer[key]
	if has == false {
		return nil, fmt.Errorf("cannot find mount for key %s", key)
	}
	return []mount.Mount{
		*m.mount,
	}, nil
}

// Prepare creates an active snapshot identified by key descending from the
// provided parent.  The returned mounts can be used to mount the snapshot
// to capture changes.
//
// If a parent is provided, after performing the mounts, the destination
// will start with the content of the parent. The parent must be a
// committed snapshot. Changes to the mounted destination will be captured
// in relation to the parent. The default parent, "", is an empty
// directory.
//
// The changes may be saved to a committed snapshot by calling Commit. When
// one is done with the transaction, Remove should be called on the key.
//
// Multiple calls to Prepare or View with the same key should fail.
func (s *emptySnapshotter) Prepare(ctx context.Context, key, parent string, opts ...snapshots.Opt) ([]mount.Mount, error) {
	log.G(ctx).
		WithField("key", key).
		WithField("parent", parent).
		Info("prepare")
	if parent == "" {
		in := NewEmptyInstance(nil, key)
		s.buffer[key] = in
		return []mount.Mount{
			*in.mount,
		}, nil
	}

	p, has := s.buffer[parent]
	if has == false {
		return nil, fmt.Errorf("cannot find parent %s", parent)
	}
	in := NewEmptyInstance(p, key)
	s.buffer[key] = in
	return []mount.Mount{
		*in.mount,
	}, nil
}

// View behaves identically to Prepare except the result may not be
// committed back to the snapshot snapshotter. View returns a readonly view on
// the parent, with the active snapshot being tracked by the given key.
//
// This method operates identically to Prepare, except that Mounts returned
// may have the readonly flag set. Any modifications to the underlying
// filesystem will be ignored. Implementations may perform this in a more
// efficient manner that differs from what would be attempted with
// `Prepare`.
//
// Commit may not be called on the provided key and will return an error.
// To collect the resources associated with key, Remove must be called with
// key as the argument.
func (s *emptySnapshotter) View(ctx context.Context, key, parent string, opts ...snapshots.Opt) ([]mount.Mount, error) {
	log.G(ctx).
		WithField("key", key).
		WithField("parent", parent).
		Info("view")
	if parent == "" {
		in := NewEmptyInstance(nil, key)
		s.buffer[key] = in
		return []mount.Mount{
			*in.mount,
		}, nil
	}

	p, has := s.buffer[parent]
	if has == false {
		return nil, fmt.Errorf("cannot find parent %s", parent)
	}
	in := NewEmptyInstance(p, key)
	s.buffer[key] = in
	return []mount.Mount{
		*in.mount,
	}, nil
}

// Commit captures the changes between key and its parent into a snapshot
// identified by name.  The name can then be used with the snapshotter's other
// methods to create subsequent snapshots.
//
// A committed snapshot will be created under name with the parent of the
// active snapshot.
//
// After commit, the snapshot identified by key is removed.
func (s *emptySnapshotter) Commit(ctx context.Context, name, key string, opts ...snapshots.Opt) error {
	log.G(ctx).
		WithField("key", key).
		WithField("name", name).
		Info("view")
	p, has := s.buffer[name]
	if has == false {
		return fmt.Errorf("cannot find name %s", name)
	}
	s.buffer[key] = p
	delete(s.buffer, name)
	return nil
}

// Remove the committed or active snapshot by the provided key.
//
// All resources associated with the key will be removed.
//
// If the snapshot is a parent of another snapshot, its children must be
// removed before proceeding.
func (s *emptySnapshotter) Remove(ctx context.Context, key string) error {
	log.G(ctx).
		WithField("key", key).
		Info("remove")
	delete(s.buffer, key)
	return nil
}

// Walk will call the provided function for each snapshot in the
// snapshotter which match the provided filters. If no filters are
// given all items will be walked.
// Filters:
//  name
//  parent
//  kind (active,view,committed)
//  labels.(label)
func (s *emptySnapshotter) Walk(ctx context.Context, fn snapshots.WalkFunc, filters ...string) error {
	log.G(ctx).
		WithField("filters", filters).
		Info("walk")
	for _, v := range s.buffer {
		p := ""
		if v.parent != nil {
			p = v.parent.id
		}
		if err := fn(ctx, snapshots.Info{
			Kind:    0,
			Name:    v.id,
			Parent:  p,
			Labels:  nil,
			Created: time.Time{},
			Updated: time.Time{},
		}); err != nil {
			return err
		}
	}
	return nil
}

// Close releases the internal resources.
//
// Close is expected to be called on the end of the lifecycle of the snapshotter,
// but not mandatory.
//
// Close returns nil when it is already closed.
func (s *emptySnapshotter) Close() error {
	return nil
}
