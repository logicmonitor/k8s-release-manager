package export

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/logicmonitor/k8s-release-manager/pkg/healthz"
	log "github.com/sirupsen/logrus"
)

var m *Export

func (m *Export) serveStats() {
	// Health check.
	http.HandleFunc("/healthz", healthz.HandleFunc)
	http.HandleFunc("/releases", m.releasesFunc)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (m *Export) releasesFunc(w http.ResponseWriter, req *http.Request) { // nolint: unparam
	var message []byte
	code := http.StatusOK

	releases, err := m.State.Releases.StoredReleaseNames()
	if err != nil {
		code = http.StatusInternalServerError
		message = []byte(fmt.Sprintf("Error retrieving stored releases: %v", err))
		respond(w, code, message)
		return
	}

	message, err = json.Marshal(releases)
	if err != nil {
		code = http.StatusInternalServerError
		message = []byte(fmt.Sprintf("Error formatting response: %v", err))
	}
	respond(w, code, message)
	return
}

func respond(w http.ResponseWriter, responseCode int, responseBody []byte) {
	w.WriteHeader(responseCode)
	_, err := w.Write(responseBody)
	if err != nil {
		log.Errorf("Failed to write releases: %v", err)
	}
}
