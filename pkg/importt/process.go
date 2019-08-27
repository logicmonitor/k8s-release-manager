package importt

import (
	"github.com/logicmonitor/k8s-release-manager/pkg/config"
	"github.com/logicmonitor/k8s-release-manager/pkg/constants"
	"github.com/logicmonitor/k8s-release-manager/pkg/release"
	log "github.com/sirupsen/logrus"
	rls "k8s.io/helm/pkg/proto/hapi/release"
)

func processReleases(releases []*rls.Release, config *config.ImportConfig) ([]*rls.Release, error) {
	releases = filterReleasesByNamespace(releases, config)
	releases, err := updateValues(releases, config)
	if err != nil {
		return nil, err
	}
	releases = updateNamespace(releases, config)
	return releases, nil
}

func filterReleasesByNamespace(releases []*rls.Release, config *config.ImportConfig) []*rls.Release {
	if config.Namespace != "" {
		return includeReleasesByNamespace(releases, config)
	}
	if len(config.ExcludeNamespaces) > 0 {
		return excludeReleasesByNamespace(releases, config)
	}
	return releases
}

func includeReleasesByNamespace(releases []*rls.Release, config *config.ImportConfig) []*rls.Release {
	var deploy []*rls.Release
	for _, r := range releases {
		if r.Namespace == config.Namespace {
			deploy = append(deploy, r)
		}
	}
	return deploy
}

func updateNamespace(releases []*rls.Release, config *config.ImportConfig) []*rls.Release {
	for _, r := range releases {
		// update the target namespace if option specified
		if config.Target != "" {
			r.Namespace = config.Target
		}
	}
	return releases
}

func excludeReleasesByNamespace(releases []*rls.Release, config *config.ImportConfig) []*rls.Release {
	var deploy []*rls.Release
	for _, r := range releases {
		bMatch := false
		for _, namespace := range config.ExcludeNamespaces {
			if r.Namespace == namespace {
				bMatch = true
				break
			}
		}
		if bMatch {
			continue
		}
		deploy = append(deploy, r)
	}
	return deploy
}

func updateValues(releases []*rls.Release, config *config.ImportConfig) ([]*rls.Release, error) {
	log.Debugf("Updating release values\n")
	var err error
	for _, r := range releases {
		for k, v := range config.Values {
			r, err = release.UpdateValue(r, k, v)
			if err != nil {
				return nil, err
			}
		}
	}
	return releases, nil
}

func updateManagerStoragePath(r *rls.Release, path string) (*rls.Release, error) {
	return release.UpdateValue(r, constants.ValueStoragePath, path)
}
