package state

import (
	"bytes"
	"encoding/json"
	"io"
)

// Info represents the state information that will be written to the backend
type Info struct {
	ReleaseFilename string
	ReleaseName     string
	ReleaseVersion  int32
}

// Serialize the state information for writing to disk
func (i *Info) Serialize() (io.Reader, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}
