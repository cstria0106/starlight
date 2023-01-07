grpc/snapshotter.go
create http stream
```
func (o *snapshotter) pullImage ...

rc, headerSize, err := o.remote.FetchWithString(sourceImage, targetImage)
```

fs/receiver.go
manage file download
```
NewReceiver receives reader (the download stream)
```

fs/fs.go
implements starlight file system
declared in github.com/hanwen/go-fuse/fs/api.go