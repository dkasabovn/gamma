package objectstore

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/labstack/gommon/log"
)

type fsStore struct {
	folderDir string
}

func newFsStore() Storage {
	root := os.Getenv("FS_STORE_ROOT")
	stat, err := os.Stat(root)
	if err != nil || !stat.IsDir() {
		panic("Could not open root folder")
	}
	return &fsStore{
		folderDir: root,
	}
}

func (f *fsStore) Get(ctx context.Context, key string) (Object, error) {
	b, err := os.ReadFile(path.Join(f.folderDir, key))
	if err != nil {
		return Object{}, err
	}
	return Object{
		Data: b,
		Key:  key,
	}, nil
}

func (f *fsStore) Put(ctx context.Context, key string, put Object) (string, error) {
	fullPath := path.Join(f.folderDir, key)
	log.Infof("fullPath: %s", fullPath)
	err := os.MkdirAll(filepath.Dir(fullPath), 0755)
	if err != nil {
		log.Errorf("Could not create key path: %v", err)
		return "", err
	}

	if err := ioutil.WriteFile(fullPath, put.Data, 0600); err != nil {
		log.Errorf("Could not write file: %v", err)
		return "", err
	}
	return fmt.Sprintf("file://%s", fullPath), nil
}
