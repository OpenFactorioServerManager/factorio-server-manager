package lockfile

import (
    "sync"
    "path"
    "os"
    "log"
)

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

func makeAbsolutePath(target string) string {
    if path.IsAbs(target) {
        return target
    }

    wd, err := os.Getwd()
    if err != nil {
        log.Fatalf("error get working directory: %s", err)
        return ""
    }

    return path.Join(wd, target)
}

func (fl *FileLock) Lock(path string) {
    fl.m.Lock()
    defer fl.m.Unlock()

    path = makeAbsolutePath(path)

    if fl.Locks[path].Read == 0 && fl.Locks[path].Write == 0 {
        lock := fl.Locks[path]
        lock.Write = 1
        fl.Locks[path] = lock
    }
}

func (fl *FileLock) Unlock(path string) {
    fl.m.Lock()
    defer fl.m.Unlock()

    path = makeAbsolutePath(path)

    lock := fl.Locks[path]
    if lock.Read == 0 && lock.Write == 1 {
        lock.Write = 0
        fl.Locks[path] = lock
    }
}

func (fl *FileLock) RLock(path string) {
    fl.m.Lock()
    defer fl.m.Unlock()

    path = makeAbsolutePath(path)

    lock := fl.Locks[path]
    if lock.Write == 0 {
        lock.Read++
        fl.Locks[path] = lock
    }
}

func (fl *FileLock) RUnlock(path string) {
    fl.m.Lock()
    defer fl.m.Unlock()

    path = makeAbsolutePath(path)

    lock := fl.Locks[path]
    if lock.Read > 0 && lock.Write == 0 {
        lock.Read--
        fl.Locks[path] = lock
    }
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
