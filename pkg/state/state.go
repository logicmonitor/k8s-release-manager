package state

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"

	"github.com/logicmonitor/k8s-release-manager/pkg/backend"
	"github.com/logicmonitor/k8s-release-manager/pkg/config"
	"github.com/logicmonitor/k8s-release-manager/pkg/constants"
	"github.com/logicmonitor/k8s-release-manager/pkg/release"
	log "github.com/sirupsen/logrus"
	rls "k8s.io/helm/pkg/proto/hapi/release"
)

// State represents the release manager's state information
type State struct {
	Backend backend.Backend
	Config  *config.Config
	init    bool
}

// Init the release manager state
func (s *State) Init() error {
	s.init = false
	if s.Config.Manager.ReleaseName != "" {
		path := s.remoteFilePath(constants.ManagerStateFilename)
		log.Infof("Deleting old state %s", path)
		err := s.Backend.Delete(path)
		if err != nil {
			log.Warnf("Error cleaning up old release manager state: %v", err)
		}
	}
	return nil
}

// Update updates the release manager state on the backend
func (s *State) Update(releases []*rls.Release) error {
	if s.Config.Manager.ReleaseName == "" {
		log.Debugf("--release-name not specified. Ignoring state.")
		return nil
	}

	// locate the release managing this application
	for _, r := range releases {
		if s.isManagerRelease(r.GetName()) {
			return s.updateState(&Info{
				ReleaseFilename: release.Filename(r),
				ReleaseName:     s.Config.Manager.ReleaseName,
				ReleaseVersion:  r.GetVersion(),
			})
		}
	}

	// if the manager release no longer exists, delete the remote state
	log.Debugf("Release manager release %s doesn't exist. Deleting.", s.Config.Manager.ReleaseName)
	err := s.delete()
	if err != nil {
		return err
	}
	return nil
}

func (s *State) updateState(i *Info) (err error) {
	update := false

	// don't attempt to read the remote state if this is our first update
	if s.init {
		// check to see if the state is stale
		oldInfo, e := s.read()
		if e != nil {
			log.Warnf("Error reading remote state: %v", e)
			update = true
		}

		if !reflect.DeepEqual(i, oldInfo) {
			update = true
		}
	} else {
		s.init, update = true, true
	}

	if update || !s.init {
		log.Debugf("Updating state %s.", i.ReleaseName)
		err = s.write(i)
		if err != nil {
			return
		}
	}
	return err
}

func (s *State) read() (i *Info, err error) {
	path := s.remoteFilePath(constants.ManagerStateFilename)
	log.Infof("Reading state from %s", path)
	f, err := s.Backend.Read(path)
	if err != nil {
		return nil, err
	}

	i = &Info{}
	err = json.Unmarshal(f, i)
	return i, err
}

func (s *State) write(i *Info) error {
	f, err := i.Serialize()
	if err != nil {
		return err
	}
	return s.Backend.Write(s.remoteFilePath(constants.ManagerStateFilename), f)
}

func (s *State) delete() error {
	return s.Backend.Delete(s.remoteFilePath(constants.ManagerStateFilename))
}

func (s *State) isManagerRelease(name string) bool {
	// check if the release name was explicitly set by flag
	if s.Config.Manager.ReleaseName != "" && name == s.Config.Manager.ReleaseName {
		return true
	}
	return false
}

// ReadRelease returns the remote release represented by the specified filename
func (s *State) ReadRelease(f string) (*rls.Release, error) {
	// TODO test
	path := s.remoteFilePath(f)
	log.Debugf("Reading remote release %s", path)
	r, err := s.Backend.Read(path)
	if err != nil {
		return nil, err
	}
	return release.FromFile(r)
}

// WriteRelease writes the specified release to the backend
func (s *State) WriteRelease(r *rls.Release) error {
	f, err := release.ToFile(r)
	if err != nil {
		return err
	}

	path := s.remoteFilePath(release.Filename(r))
	log.Debugf("Writing remote release %s", path)
	return s.Backend.Write(path, f)
}

// DeleteRelease deletes the remote release represented by the specified filename
func (s *State) DeleteRelease(f string) error {
	path := s.remoteFilePath(f)
	log.Debugf("Deleting remote release %s", path)
	return s.Backend.Delete(path)
}

// StoredReleases returns the list of release structs currently stored in the backend
func (s *State) StoredReleases() (ret []*rls.Release, err error) {
	// TODO test
	filenames, err := s.StoredReleaseNames()
	if err != nil {
		return ret, err
	}

	for _, f := range filenames {
		r, e := s.ReadRelease(f)
		if e != nil {
			log.Warnf("%v", e)
			continue
		}
		ret = append(ret, r)
	}
	return ret, err
}

// StoredReleaseNames returns the list of release filenames currently stored in the backend
func (s *State) StoredReleaseNames() (ret []string, err error) {
	log.Debugf("Finding releases stored in the backend.")
	names, err := s.Backend.List(s.Config.Manager.StoragePath)
	if err != nil {
		return ret, err
	}

	// ignore non release files in path, e.g. state, other cruft outside our control
	r, _ := regexp.Compile(fmt.Sprintf("^.+%s$", constants.ReleaseExtension))
	for _, n := range names {
		if r.MatchString(n) {
			ret = append(ret, n)
		}
	}
	return ret, err
}

// remoteFilePath returns the full appropriate backend file path based on the app's configuration
func (s *State) remoteFilePath(name string) string {
	if s.Config.Manager.StoragePath == s.Backend.PathSeparator() {
		return name
	}
	return s.Config.Manager.StoragePath + s.Backend.PathSeparator() + name
}
