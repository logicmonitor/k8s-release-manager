package export

import (
	"time"

	"github.com/logicmonitor/k8s-release-manager/pkg/config"
	"github.com/logicmonitor/k8s-release-manager/pkg/healthz"
	"github.com/logicmonitor/k8s-release-manager/pkg/lmhelm"
	"github.com/logicmonitor/k8s-release-manager/pkg/state"
	log "github.com/sirupsen/logrus"
)

// Export exports releases
type Export struct {
	Config     *config.Config
	HelmClient *lmhelm.Client
	State      *state.State
}

// New instantiates and returns a Export and an error if any.
func New(rlsmgrconfig *config.Config, state *state.State) (*Export, error) {

	helmClient := &lmhelm.Client{}

	err := helmClient.Init(rlsmgrconfig.ClusterConfig, rlsmgrconfig.OptionsConfig)
	if err != nil {
		return nil, err
	}

	return &Export{
		Config:     rlsmgrconfig,
		HelmClient: helmClient,
		State:      state,
	}, nil
}

// Run the Export.
func (m *Export) Run() error {
	if m.Config.Export.ReleaseName != "" {
		log.Infof("Cleaning old state")
		err := m.State.Remove()
		if err != nil {
			log.Warnf("Error cleaning up old release manager state: %v", err)
		}
	}

	// if not daemon mode, run once and exit
	if !m.Config.Export.DaemonMode {
		return m.strategy()()
	}
	return m.run()
}

func (m *Export) strategy() func() error {
	if m.Config.DryRun {
		return m.printReleases
	}
	return m.exportReleases
}

func (m *Export) run() error {
	// start stats server
	go m.serveStats()

	// daemon mode. run periodically forever
	for {
		log.Debugf("Checking for installed releases")
		err := m.strategy()()
		if err != nil {
			healthz.IncrementFailure()
			log.Errorf("%v", err)
		} else {
			healthz.ResetFailure()
		}
		time.Sleep(time.Duration(m.Config.Export.PollingInterval) * time.Second)
	}
}
