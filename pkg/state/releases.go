package state

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/logicmonitor/k8s-release-manager/pkg/backend"
	"github.com/logicmonitor/k8s-release-manager/pkg/config"
	"github.com/logicmonitor/k8s-release-manager/pkg/constants"
	"github.com/logicmonitor/k8s-release-manager/pkg/release"
	log "github.com/sirupsen/logrus"
	rls "helm.sh/helm/v3/pkg/release"
)

// ReleaseState is a wrapper for interacting with the stored release info
type ReleaseState struct {
	Backend backend.Backend
	Config  *config.Config
}

// ReadRelease returns the remote release represented by the specified filename
func (rs *ReleaseState) ReadRelease(f string) (*rls.Release, error) {
	log.Debugf("Reading remote release %s", f)
	b, err := rs.Backend.Read(f)
	if err != nil {
		return nil, err
	}
	return release.FromFile(b)
}

// WriteRelease writes the specified release to the backend
func (rs *ReleaseState) WriteRelease(r *rls.Release) error {
	f, err := release.ToFile(r)
	if err != nil {
		return err
	}
	if rs.Config.DryRun {
		return nil
	}
	log.Debugf("Writing remote release %s", release.Filename(r))
	return rs.Backend.Write(release.Filename(r), f)
}

// DeleteRelease deletes the remote release represented by the specified filename
func (rs *ReleaseState) DeleteRelease(f string) error {
	if rs.Config.DryRun {
		return nil
	}
	log.Debugf("Removing remote release %s", f)
	return rs.Backend.Delete(f)
}

// StoredReleases returns the list of release structs currently stored in the backend
func (rs *ReleaseState) StoredReleases() (ret []*rls.Release, err error) {
	filenames, err := rs.StoredReleaseNames()
	if err != nil {
		return ret, err
	}

	var wg sync.WaitGroup
	for _, f := range filenames {
		wg.Add(1)
		go func(f string, ret *[]*rls.Release) {
			defer wg.Done()
			r, e := rs.ReadRelease(f)
			if e != nil {
				log.Warnf("%v", e)
				return
			}
			*ret = append(*ret, r)
		}(f, &ret)
	}
	wg.Wait()
	return ret, err
}

// StoredReleaseNames returns the list of release filenames currently stored in the backend
func (rs *ReleaseState) StoredReleaseNames() (ret []string, err error) {
	log.Debugf("Finding releases stored in the backend.")
	names, err := rs.Backend.List()
	if err != nil {
		return ret, err
	}

	// ignore non release files in path, e.g. state, other cruft outside our control
	r, err := regexp.Compile(fmt.Sprintf("^.+%s$", constants.ReleaseExtension))
	if err != nil {
		return nil, err
	}

	for _, n := range names {
		if r.MatchString(n) {
			ret = append(ret, n)
		}
	}
	return ret, err
}
