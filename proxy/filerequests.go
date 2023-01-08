package proxy

import "sync"

type fileRequestItem struct {
	Source                       int
	SourceOffset, CompressedSize int64
}

type FileRequests struct {
	lock  sync.RWMutex
	items []fileRequestItem
}

func (f *FileRequests) Pop() (bool, int, int64, int64) {
	f.lock.RLock()
	defer f.lock.RUnlock()
	if len(f.items) == 0 {
		return false, 0, 0, 0
	}

	e := f.items[0]
	f.items = f.items[1:]
	return true, e.Source, e.SourceOffset, e.CompressedSize
}

func (f *FileRequests) Push(source int, sourceOffset, compressedSize int64) {
	f.lock.Lock()
	f.items = append(f.items, fileRequestItem{source, sourceOffset, compressedSize})
	f.lock.Unlock()
}
