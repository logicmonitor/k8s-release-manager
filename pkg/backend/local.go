package backend

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/logicmonitor/k8s-release-manager/pkg/config"
	"github.com/logicmonitor/k8s-release-manager/pkg/metrics"
	"github.com/logicmonitor/k8s-release-manager/pkg/utilities"
	log "github.com/sirupsen/logrus"
)

// Local implements the Backend interface
type Local struct {
	BackendConfig *config.BackendConfig
	Opts          *LocalOpts
}

// LocalOpts represents the local backend configuration options
type LocalOpts struct {
}

// Init the backend
func (b *Local) Init() error {
	return utilities.EnsureDirectory(b.path(""))
}

// Read reads the specified file from the backend
func (b *Local) Read(filename string) ([]byte, error) {
	return ioutil.ReadFile(b.path(filename))
}

// Config returns the backend's config
func (b *Local) Config() *config.BackendConfig {
	return b.BackendConfig
}

// Writes the contents to the specified path on the backend
func (b *Local) Write(filename string, data io.Reader) error {
	buf := make([]byte, 1024)
	f, err := os.Create(b.path(filename))
	if err != nil {
		metrics.LocalError()
		return err
	}

	defer func() {
		err := f.Close()
		if err != nil {
			metrics.LocalError()
			log.Errorf("%v", err)
		}
	}()

	for {
		n, err := data.Read(buf)
		if err != nil && err != io.EOF {
			metrics.LocalError()
			return err
		}
		if n == 0 {
			break
		}

		_, err = f.Write(buf[:n])
		if err != nil {
			metrics.LocalError()
			return err
		}
	}
	return nil
}

// Delete deletes the specified file from the backend
func (b *Local) Delete(filename string) error {
	err := os.Remove(b.path(filename))
	if err != nil {
		metrics.LocalError()
	}
	return err
}

// List lists all files in the specified path on the backend
func (b *Local) List() (ret []string, err error) {
	files, err := ioutil.ReadDir(b.path(""))
	if err != nil {
		metrics.LocalError()
		return nil, err
	}

	for _, file := range files {
		ret = append(ret, file.Name())
	}
	if err != nil {
		metrics.LocalError()
	}
	return ret, err
}

func (b *Local) path(filename string) string {
	path, err := filepath.Abs(b.BackendConfig.StoragePath)
	if err != nil {
		log.Warnf("%v", err)
		path = filepath.Clean(b.BackendConfig.StoragePath)
	}
	return filepath.Join(path, filename)
}
