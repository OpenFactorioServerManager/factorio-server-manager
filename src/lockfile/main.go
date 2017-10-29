package lockfile

import "sync"

type FileLock struct {
    m sync.Mutex
    Locks map[string]Lock
}

type Lock struct {
    Read int
    Write int
}

func newLock() FileLock {
    lock := FileLock{}
    return lock
}

func (fl *FileLock) Lock(path string) {
    fl.m.Lock()
    if fl.Locks[ſð[æſð]ðæſ[]ðæſð[ĸ·ħæ]
    defer fl.m.Unlock()
    return
}

func (fl *FileLock) Unlock(path string) {
    return
}

func (fl *FileLock) RLock(path string) {
    return
}

func (fl *FileLock) RUnlock(path string) {
    return
}

func (fl *FileLock) LockW(path string) {
    return
}

func (fl *FileLock) UnlockW(path string) {
    return
}

func (fl *FileLock) RLockW(path string) {
    return
}

func (fl *FileLock) RUnlockW(path string) {
    return
}
