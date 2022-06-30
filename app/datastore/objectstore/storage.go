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
		if system.ENVIRONMENT == "local" {
			storeInst = newFsStore()
		} else {
			storeInst = newCloudStore()
		}
	})
	return storeInst
}
