package proxy

import (
	"sync"

	"github.com/mc256/starlight/util"
)

type FileRequests struct {
	lock  sync.RWMutex
	items []util.FileRequest
}

func (f *FileRequests) Pop() (bool, util.FileRequest) {
	f.lock.RLock()
	defer f.lock.RUnlock()
	if len(f.items) == 0 {
		return false, util.FileRequest{}
	}

	e := f.items[0]
	f.items = f.items[1:]
	return true, e
}

func (f *FileRequests) Push(fr util.FileRequest) {
	f.lock.Lock()
	f.items = append(f.items, fr)
	f.lock.Unlock()
}
