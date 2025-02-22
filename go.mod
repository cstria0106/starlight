module github.com/mc256/starlight

go 1.18

require (
	github.com/containerd/containerd v1.6.12
	github.com/containerd/continuity v0.3.0
	github.com/google/go-containerregistry v0.5.1
	github.com/google/uuid v1.3.0

	// resolved overlayfs OICTL issue in https://github.com/hanwen/go-fuse/pull/408
	// tested on Kernel 5.15.0-52-generic
	github.com/hanwen/go-fuse/v2 v2.1.1-0.20221003202731-4c25c9c1eece
	github.com/lib/pq v1.10.6
	github.com/opencontainers/go-digest v1.0.0
	github.com/opencontainers/image-spec v1.0.3-0.20220114050600-8b9d41f48198
	github.com/opencontainers/runtime-spec v1.0.3-0.20210326190908-1c3f411f0417
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1

	// newer version causes issue #24
	github.com/urfave/cli/v2 v2.3.0
	go.etcd.io/bbolt v1.3.6
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b
	golang.org/x/sync v0.0.0-20220601150217-0de741cfad7f
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f
	google.golang.org/grpc v1.47.0
)

require google.golang.org/protobuf v1.28.0

require (
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/Microsoft/hcsshim v0.9.5 // indirect
	github.com/containerd/cgroups v1.0.4 // indirect
	github.com/containerd/console v1.0.3 // indirect
	github.com/containerd/fifo v1.0.0 // indirect
	github.com/containerd/stargz-snapshotter/estargz v0.4.1 // indirect
	github.com/containerd/ttrpc v1.1.0 // indirect
	github.com/containerd/typeurl v1.0.2 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/docker/cli v20.10.17+incompatible // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v20.10.17+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.6.3 // indirect
	github.com/docker/go-events v0.0.0-20190806004212-e31b211e4f1c // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/klauspost/compress v1.15.6 // indirect
	github.com/moby/locker v1.0.1 // indirect
	github.com/moby/sys/mountinfo v0.6.2 // indirect
	github.com/moby/sys/signal v0.7.0 // indirect
	github.com/opencontainers/runc v1.1.2 // indirect
	github.com/opencontainers/selinux v1.10.1 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/urfave/cli v1.22.9 // indirect
	go.opencensus.io v0.23.0 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220602131408-e326c6e8e9c8 // indirect
)
