package lockfile

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileLock struct {
	m     sync.Mutex
	Locks map[string]Lock
}

type Lock struct {
	Read  int
	Write int
}

var ErrorLocking error = errors.New("error locking file")
var ErrorAlreadyLocked error = errors.New("file already locked by another routine")

func NewLock() FileLock {
	lock := FileLock{
		Locks: make(map[string]Lock),
	}
	return lock
}

func makeAbsolutePath(target string) string {
	if filepath.IsAbs(target) {
		return target
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("error get working directory: %s", err)
		return ""
	}

	return filepath.Join(wd, target)
}

func (fl *FileLock) Lock(filePath string) error {
	fl.m.Lock()
	defer fl.m.Unlock()

	filePath = makeAbsolutePath(filePath)

	if fl.Locks[filePath].Read == 0 && fl.Locks[filePath].Write == 0 {
		lock := fl.Locks[filePath]
		lock.Write = 1
		fl.Locks[filePath] = lock
	} else {
		return ErrorAlreadyLocked
	}
	return nil
}

func (fl *FileLock) Unlock(filePath string) error {
	fl.m.Lock()
	defer fl.m.Unlock()

	filePath = makeAbsolutePath(filePath)

	lock := fl.Locks[filePath]
	if lock.Read == 0 && lock.Write == 1 {
		lock.Write = 0
		fl.Locks[filePath] = lock
	} else {
		return ErrorAlreadyLocked
	}
	return nil
}

func (fl *FileLock) RLock(filePath string) error {
	fl.m.Lock()
	defer fl.m.Unlock()

	filePath = makeAbsolutePath(filePath)

	lock := fl.Locks[filePath]
	if lock.Write == 0 {
		lock.Read++
		fl.Locks[filePath] = lock
	} else {
		return ErrorAlreadyLocked
	}
	return nil
}

func (fl *FileLock) RUnlock(filePath string) error {
	fl.m.Lock()
	defer fl.m.Unlock()

	filePath = makeAbsolutePath(filePath)

	lock := fl.Locks[filePath]
	if lock.Read > 0 && lock.Write == 0 {
		lock.Read--
		fl.Locks[filePath] = lock
	} else {
		return ErrorAlreadyLocked
	}
	return nil
}

func (fl *FileLock) LockW(filePath string) {
	for {
		err := fl.Lock(filePath)
		if err == ErrorAlreadyLocked {
			time.Sleep(time.Second * 2)
			log.Println("file locked wait two seconds to access write-lock")
		}

		if err == nil {
			break
		}
	}
	return
}

func (fl *FileLock) RLockW(filePath string) {
	for {
		err := fl.RLock(filePath)

		if err == ErrorAlreadyLocked {
			time.Sleep(time.Second * 2)
			log.Println("file locked ... wait two seconds to try to access read-lock")
		}

		if err == nil {
			break
		}
	}
	return
}
