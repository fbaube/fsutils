package fss

import (
	"io/fs"
	"sync"
)

// BaseFS is an abstract class, an incomplete class.
type BaseFS struct {
	inputFS     fs.FS
	rootAbsPath string
	namespace   string
	sync.Mutex
	isLocked bool
}

// Lock is func (*Mutex) Lock :: If the lock is already in use,
// the calling goroutine blocks until the mutex is available.
func (bfs *BaseFS) Lock() (success bool) {
	if bfs.isLocked {
		return false
	}
	bfs.isLocked = true
	bfs.Mutex.Lock()
	return true
}

// Unlock is func (*Mutex) Unlock :: It is a run-time error
// if m is not locked on entry to Unlock.
func (bfs *BaseFS) Unlock() {
	if !bfs.isLocked {
		panic("Unlock failed: is not locked, would throw RTE")
	}
	bfs.Mutex.Unlock()
	bfs.isLocked = false
}

// IsLocked is duh.
func (bfs *BaseFS) IsLocked() bool { return bfs.isLocked }
