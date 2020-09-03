package export

import (
	"github.com/logicmonitor/k8s-release-manager/pkg/release"
	log "github.com/sirupsen/logrus"
	rls "helm.sh/helm/v3/pkg/release"
)

func (m *Export) currentReleases() ([]*rls.Release, error) {
	log.Debugf("Finding installed releases.")
	releases, err := m.HelmClient.ListInstalledReleases()
	if err != nil {
		return nil, err
	}

	if len(m.Config.Export.Namespaces) != 0 {
		releases = m.filterReleasesByNamespace(releases)
	}

	if m.Config.DebugMode {
		for _, r := range releases {
			log.Debugf("Found installed release %s", release.Filename(r))
		}
	}
	return releases, nil
}

func (m *Export) filterReleasesByNamespace(releases []*rls.Release) []*rls.Release {
	var results []*rls.Release
	for _, r := range releases {
		for _, ns := range m.Config.Export.Namespaces {
			if r.Namespace == ns {
				results = append(results, r)
				break
			}
		}
	}
	return results
}

func (m *Export) storedReleases() ([]string, error) {
	names, err := m.State.Releases.StoredReleaseNames()
	if m.Config.DebugMode && err == nil {
		for _, r := range names {
			log.Debugf("Found stored release %s", r)
		}
	}
	return names, err
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
			log.Debugf("Found release to save %s", release.Filename(c))
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
			log.Debugf("Found release to delete %s", s)
		}
	}
	return ret
}
