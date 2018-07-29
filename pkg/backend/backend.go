package backend

import "io"

// Backend is an interface that abstracts operations on a data store
type Backend interface {
	Read(path string) ([]byte, error)
	Write(path string, data io.Reader) error
	Delete(path string) error
	List(path string) ([]string, error)
	PathSeparator() string
}
