package metrics

import (
	"expvar"
	"runtime"
	"sync"
)

var (
	c    *expvar.Map
	e    *expvar.Map
	once sync.Once
)

func init() {
	once.Do(func() {
		e = expvar.NewMap("errors")
		c = expvar.NewMap("jobs")
		c.Add("DeleteCount", 0)
		c.Add("FailedJobs", 0)
		c.Add("TotalJobs", 0)
		c.Add("UploadCount", 0)
		e.Add("DeleteErrors", 0)
		e.Add("HelmErrors", 0)
		e.Add("StateErrors", 0)
		e.Add("UploadErrors", 0)
	})
	expvar.Publish("goroutines", expvar.Func(goroutines))
}

// JobError increments the job error count by 1.
func JobError() {
	c.Add("FailedJobs", 1)
}

// JobCount increments the jobcount by 1.
func JobCount() {
	c.Add("TotalJobs", 1)
}

// UploadError increments the upload error count by 1.
func UploadError() {
	e.Add("UploadErrors", 1)
}

// DeleteError increments the delete error count by 1.
func DeleteError() {
	e.Add("DeleteErrors", 1)
}

// StateError increments the state error count by 1.
func StateError() {
	e.Add("UploadErrors", 1)
}

// HelmError increments the helm error count by 1.
func HelmError() {
	e.Add("UploadErrors", 1)
}

// UploadCount increments the upload count by 1.
func UploadCount() {
	c.Add("UploadCount", 1)
}

// DeleteCount increments the delete count by 1.
func DeleteCount() {
	c.Add("DeleteCount", 1)
}

func goroutines() interface{} {
	return runtime.NumGoroutine()
}
