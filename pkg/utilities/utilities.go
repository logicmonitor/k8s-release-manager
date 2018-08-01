package utilities

import (
	"os"

	"github.com/logicmonitor/k8s-release-manager/pkg/backend"
)

// FileExists returns true if the file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		return false
	}
}

// RemoteFilePath returns the full appropriate backend file path based on the app's configuration
func RemoteFilePath(b backend.Backend, name string) string {
	if b.Config().StoragePath == b.PathSeparator() {
		return name
	}
	return b.Config().StoragePath + b.PathSeparator() + name
}
