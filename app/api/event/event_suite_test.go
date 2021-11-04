package event_test

import (
	"gamma/app/system"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEvent(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Dir(strings.SplitAfter(filename, "backend/")[0])
	log.Println(dir)
	os.Chdir(dir)
	system.Initialize()

	RegisterFailHandler(Fail)
	RunSpecs(t, "Event Suite")
}
