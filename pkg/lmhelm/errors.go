package lmhelm

import "strings"

const (
	releaseExistsError string = "cannot re-use a name that is still in use"
)

// ErrorReleaseExists returns true if the helm error is the result of trying to reinstall and existing release
func ErrorReleaseExists(err error) bool {
	return strings.Contains(err.Error(), releaseExistsError)
}
