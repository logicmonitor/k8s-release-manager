package export

import (
	"fmt"
	"sync"

	"github.com/logicmonitor/k8s-release-manager/pkg/metrics"
	"github.com/logicmonitor/k8s-release-manager/pkg/release"
	log "github.com/sirupsen/logrus"
	rls "helm.sh/helm/v3/pkg/release"
)

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
		metrics.HelmError()
		metrics.JobError()
		return err
	}

	storedReleaseNames, err := m.storedReleases()
	if err != nil {
		metrics.StateError()
		metrics.JobError()
		return err
	}

	err = m.State.Update(currentReleases)
	if err != nil {
		metrics.StateError()
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
		metrics.JobCount()
		wg.Add(1)
		go func(r *rls.Release) {
			defer wg.Done()
			err := m.State.Releases.WriteRelease(r)
			if err != nil {
				metrics.SaveError()
				metrics.JobError()
				log.Warnf("%v", err)
			} else {
				metrics.SaveCount()
			}
		}(r)
	}
	wg.Wait()
}

func (m *Export) deleteReleases(current []*rls.Release, stored []string) {
	var wg sync.WaitGroup

	deletedReleases := deletedReleases(current, stored)
	for _, f := range deletedReleases {
		metrics.JobCount()
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			err := m.State.Releases.DeleteRelease(f)
			if err != nil {
				metrics.DeleteError()
				metrics.JobError()
				log.Warnf("%v", err)
			} else {
				metrics.DeleteCount()
			}
		}(f)
	}

	wg.Wait()
}
