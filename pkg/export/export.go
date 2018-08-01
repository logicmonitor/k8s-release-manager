package export

import (
	"fmt"
	"sync"
	"time"

	"github.com/logicmonitor/k8s-release-manager/pkg/backend"
	"github.com/logicmonitor/k8s-release-manager/pkg/client"
	"github.com/logicmonitor/k8s-release-manager/pkg/config"
	"github.com/logicmonitor/k8s-release-manager/pkg/lmhelm"
	"github.com/logicmonitor/k8s-release-manager/pkg/release"
	"github.com/logicmonitor/k8s-release-manager/pkg/state"
	log "github.com/sirupsen/logrus"
	rls "k8s.io/helm/pkg/proto/hapi/release"
)

// Export polls Tiller and exports releases
type Export struct {
	Config     *config.Config
	HelmClient *lmhelm.Client
	State      *state.State
}

// New instantiates and returns a Export and an error if any.
func New(rlsmgrconfig *config.Config, backend backend.Backend) (*Export, error) {
	kubernetesClient, kubernetesConfig, err := client.KubernetesClient(rlsmgrconfig.ClusterConfig)
	if err != nil {
		return nil, err
	}

	helmClient := &lmhelm.Client{}
	err = helmClient.Init(rlsmgrconfig.Helm, kubernetesClient, kubernetesConfig)
	if err != nil {
		return nil, err
	}

	return &Export{
		Config:     rlsmgrconfig,
		HelmClient: helmClient,
		State: &state.State{
			Backend: backend,
			Config:  rlsmgrconfig,
		},
	}, nil
}

// Run the Export.
func (m *Export) Run() error {
	var run func() error
	if !m.Config.DryRun {
		run = m.exportReleases
	} else {
		run = m.printReleases
	}

	err := m.State.Init()
	if err != nil {
		return err
	}

	if m.Config.Export.ReleaseName != "" && !m.Config.DryRun {
		log.Infof("Cleaning old state")
		err := m.State.Remove
		if err != nil {
			log.Warnf("Error cleaning up old release manager state: %v", err)
		}
	}

	if !m.Config.Export.DaemonMode {
		return run()
	}

	for {
		log.Debugf("Checking Tiller for installed releases")
		err := run()
		if err != nil {
			log.Errorf("%v", err)
		}
		time.Sleep(time.Duration(m.Config.Export.PollingInterval) * time.Second)
	}
}

func (m *Export) printReleases() error {
	currentReleases, err := m.currentReleases()
	if err != nil {
		return err
	}
	for _, r := range currentReleases {
		fmt.Printf("%s\n", release.ToString(r, m.Config.VerboseMode))
	}
	return nil
}

func (m *Export) exportReleases() error {
	currentReleases, err := m.currentReleases()
	if err != nil {
		return err
	}

	storedReleaseNames, err := m.State.Releases.StoredReleaseNames()
	if err != nil {
		return err
	}

	err = m.State.Update(currentReleases)
	if err != nil {
		log.Warnf("%v", err)
	}
	return m.export(currentReleases, storedReleaseNames)
}

func (m *Export) export(current []*rls.Release, stored []string) error {
	var wg sync.WaitGroup

	wg.Add(2)
	go func(current []*rls.Release, stored []string) {
		defer wg.Done()
		m.updateReleases(current, stored)
	}(current, stored)

	go func(current []*rls.Release, stored []string) {
		defer wg.Done()
		m.deleteReleases(current, stored)
	}(current, stored)

	wg.Wait()
	return nil
}

func (m *Export) updateReleases(current []*rls.Release, stored []string) {
	var wg sync.WaitGroup

	updatedReleases := updatedReleases(current, stored)
	for _, r := range updatedReleases {
		wg.Add(1)
		go func(r *rls.Release) {
			defer wg.Done()
			err := m.State.Releases.WriteRelease(r)
			if err != nil {
				log.Warnf("%v", err)
			}
		}(r)
	}
	wg.Wait()
}

func (m *Export) deleteReleases(current []*rls.Release, stored []string) {
	var wg sync.WaitGroup

	deletedReleases := deletedReleases(current, stored)
	for _, f := range deletedReleases {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			err := m.State.Releases.DeleteRelease(f)
			if err != nil {
				log.Warnf("%v", err)
			}
		}(f)
	}

	wg.Wait()
}

func (m *Export) currentReleases() ([]*rls.Release, error) {
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
	return ret
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
	return ret
}
