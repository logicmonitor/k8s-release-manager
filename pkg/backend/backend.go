package backend

import (
	"io"

	"github.com/logicmonitor/k8s-release-manager/pkg/config"
)

// Backend is an interface that abstracts operations on a data store
type Backend interface {
	Config() *config.BackendConfig
	Delete(filename string) error
	List() ([]string, error)
	Read(filename string) ([]byte, error)
	Write(filename string, data io.Reader) error
}
