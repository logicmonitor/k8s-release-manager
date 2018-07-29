package manager

import (
	"sync"
	"time"

	"github.com/logicmonitor/k8s-release-manager/pkg/backend"
	"github.com/logicmonitor/k8s-release-manager/pkg/client"
	"github.com/logicmonitor/k8s-release-manager/pkg/config"
	"github.com/logicmonitor/k8s-release-manager/pkg/lmhelm"
	"github.com/logicmonitor/k8s-release-manager/pkg/release"
	"github.com/logicmonitor/k8s-release-manager/pkg/state"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	rls "k8s.io/helm/pkg/proto/hapi/release"
)

// Manager polls Tiller and exports releases
type Manager struct {
	Client     *client.Client
	Config     *config.Config
	HelmClient *lmhelm.Client
	State      *state.State
}

// New instantiates and returns a Manager and an error if any.
func New(rlsmgrconfig *config.Config, backend backend.Backend) (*Manager, error) {
	// Instantiate the Kubernetes in cluster config.
	restconfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// Instantiate the client.
	client, _, err := client.NewForConfig(restconfig)
	if err != nil {
		return nil, err
	}

	// initialize our LM helm wrapper struct
	helmClient := &lmhelm.Client{}
	err = helmClient.Init(rlsmgrconfig, restconfig)
	if err != nil {
		return nil, err
	}

	state := &state.State{
		Backend: backend,
		Config:  rlsmgrconfig,
	}

	c := &Manager{
		Client:     client,
		Config:     rlsmgrconfig,
		HelmClient: helmClient,
		State:      state,
	}
	return c, nil
}

// Run starts starts the release manager.
func (m *Manager) Run() error {
	if !m.Config.DebugMode {
		err := m.State.Init()
		if err != nil {
			return err
		}

		for {
			log.Debugf("Checking Tiller for installed releases")
			err := m.exportReleases()
			if err != nil {
				log.Errorf("%v", err)
			}
			time.Sleep(time.Duration(m.Config.PollingInterval) * time.Second)
		}
		return nil
	}
	log.Debugf("Debug mode enabled. Printing release info.")
	for {
		err := m.printReleases()
		if err != nil {
			log.Errorf("%v", err)
		}
		time.Sleep(time.Duration(m.Config.PollingInterval) * time.Second)
	}
	return nil
}

func (m *Manager) printReleases() error {
	currentReleases, err := m.currentReleases()
	if err != nil {
		return err
	}
	for _, r := range currentReleases {
		f, err := release.ToString(r, m.Config.VerboseMode)
		if err != nil {
			log.Warnf("%v", err)
		}
		log.Debugf(f)
	}
	return nil
}

func (m *Manager) exportReleases() error {
	currentReleases, err := m.currentReleases()
	if err != nil {
		return err
	}

	storedReleaseNames, err := m.State.StoredReleaseNames()
	if err != nil {
		return err
	}

	err = m.State.Update(currentReleases)
	if err != nil {
		log.Warnf("%v", err)
	}
	return m.updateReleases(currentReleases, storedReleaseNames)
}

func (m *Manager) updateReleases(current []*rls.Release, stored []string) error {
	var wg sync.WaitGroup

	updatedReleases := updatedReleases(current, stored)
	for _, r := range updatedReleases {
		wg.Add(1)
		go func(r *rls.Release) {
			defer wg.Done()
			err := m.State.WriteRelease(r)
			if err != nil {
				log.Warnf("%v", err)
			}
		}(r)
	}

	deletedReleases := deletedReleases(current, stored)
	for _, f := range deletedReleases {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			err := m.State.DeleteRelease(f)
			if err != nil {
				log.Warnf("%v", err)
			}
		}(f)
	}

	wg.Wait()
	return nil
}

func (m *Manager) currentReleases() ([]*rls.Release, error) {
	log.Debugf("Finding releases that exist locally.")
	return m.HelmClient.ListInstalledReleases()
}

// updated returns the list of current releases that have yet to be stored
func updatedReleases(current []*rls.Release, stored []string) (ret []*rls.Release) {
	log.Debugf("Generating list of updated releases.")
	for _, c := range current {
		exists := false
		for _, s := range stored {
			if s == release.Filename(c) {
				exists = true
				break
			}
		}
		if !exists {
			ret = append(ret, c)
		}
	}
	return
}

// deleted returns the filenames of stored releases that not longer exist
func deletedReleases(current []*rls.Release, stored []string) (ret []string) {
	log.Debugf("Generating list of deleted releases.")
	for _, s := range stored {
		exists := false
		for _, c := range current {
			if s == release.Filename(c) {
				exists = true
				break
			}
		}
		if !exists {
			ret = append(ret, s)
		}
	}
	return
}
