package healthz

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

var failures int
var currentlyHealthy bool

const (
	maxFailures = 2
)

// IncrementFailure count
func IncrementFailure() {
	failures++
	Healthy()
}

// ResetFailure count
func ResetFailure() {
	failures = 0
	Healthy()
}

// Healthy returns the health of the service
func Healthy() bool {
	if !healthy() {
		if currentlyHealthy {
			log.Warnf("The service is now in an unhealthy state")
			currentlyHealthy = false
		}
		return false
	}

	if !currentlyHealthy {
		log.Infof("The service is now healthy")
		currentlyHealthy = true
	}
	return true
}

func healthy() bool {
	switch true {
	case failures >= maxFailures:
		return false
	default:
		return true
	}
}

// HandleFunc is an http handler function to expose health metrics.
func HandleFunc(w http.ResponseWriter, req *http.Request) {
	var code int
	var message string

	if Healthy() {
		code = http.StatusOK
		message = "ok"
	} else {
		code = http.StatusInternalServerError
		message = "unhealthy"
	}

	w.WriteHeader(code)
	_, err := w.Write([]byte(message))
	if err != nil {
		log.Errorf("Failed to write healthz: %v", err)
	}
}

func init() {
	failures = 0
	currentlyHealthy = true
}
