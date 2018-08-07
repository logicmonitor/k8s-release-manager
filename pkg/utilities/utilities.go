package utilities

import (
	"os"

	log "github.com/sirupsen/logrus"
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

// EnsureDirectory ensure that dir is a directory
func EnsureDirectory(dir string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		log.Debugf("Creating directory %s", dir)
		err = os.MkdirAll(dir, 0744)
		if err != nil {
			return err
		}
	}
	return nil
}
