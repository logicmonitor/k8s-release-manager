package lmhelm

import "strings"

const (
	releaseExistsError string = "cannot re-use a name that is still in use"
)

func ErrorReleaseExists(err error) bool {
	return strings.Contains(err.Error(), releaseExistsError)
}
