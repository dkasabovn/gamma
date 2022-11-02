package objectstore

import (
	"gamma/app/system"
	"sync"
)

var (
	storeSync sync.Once
	storeInst Storage
)

func GetStorage() Storage {
	storeSync.Do(func() {
		if system.GetConfig().Environment != "test" {
			storeInst = newS3()
		} else {
			storeInst = newFsStore()
		}
	})
	return storeInst
}
