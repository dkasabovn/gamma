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
		if system.ENVIRONMENT == "prod" {
			storeInst = newCloudStore()
		} else {
			storeInst = newFsStore()
		}
	})
	return storeInst
}
