package backend

import (
	"io"

	"github.com/logicmonitor/k8s-release-manager/pkg/config"
)

// Backend is an interface that abstracts operations on a data store
type Backend interface {
	Config() *config.BackendConfig
	Delete(path string) error
	List(path string) ([]string, error)
	PathSeparator() string
	Read(path string) ([]byte, error)
	Write(path string, data io.Reader) error
}
