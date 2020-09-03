package state

import (
	"fmt"
	"reflect"

	"github.com/logicmonitor/k8s-release-manager/pkg/backend"
	"github.com/logicmonitor/k8s-release-manager/pkg/config"
	"github.com/logicmonitor/k8s-release-manager/pkg/constants"
	"github.com/logicmonitor/k8s-release-manager/pkg/release"
	log "github.com/sirupsen/logrus"
	rls "helm.sh/helm/v3/pkg/release"
)

// State represents the release manager's state information
type State struct {
	Backend  backend.Backend
	Config   *config.Config
	Info     *Info
	Releases *ReleaseState
	init     bool
}

// Init the release manager state
func (s *State) Init() error {
	s.init = false
	s.Releases = &ReleaseState{
		Backend: s.Backend,
		Config:  s.Config,
	}
	return nil
}

// Name returns the state file name
func (s *State) Name() string {
	return constants.ManagerStateFilename
}

// Update updates the release manager state on the backend
func (s *State) Update(releases []*rls.Release) error {
	if s.Config.Export.ReleaseName == "" {
		return nil
	}

	log.Debugf("Updating state")
	// locate the release managing this application
	for _, r := range releases {
		if s.isManagerRelease(r.Name) {
			return s.updateState(&Info{
				ReleaseFilename: release.Filename(r),
				ReleaseName:     s.Config.Export.ReleaseName,
				ReleaseVersion:  int32(r.Version),
			})
		}
	}

	// if the manager release no longer exists, delete the remote state
	log.Debugf("Release manager release %s doesn't exist. Removing state.", s.Config.Export.ReleaseName)
	return s.delete()
}

// Read the release manager state from the backend
func (s *State) Read() error {
	exists, err := s.exists()
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	info, err := s.read()
	if err != nil {
		return err
	}
	s.Info = info
	return nil
}

// Remove the release manager state from the backend
func (s *State) Remove() error {
	return s.delete()
}

// Exists returns true if the remote state file exists
func (s *State) exists() (bool, error) {
	log.Infof("Check if remote state file %s exists", constants.ManagerStateFilename)
	f, err := s.Backend.List()
	if err != nil {
		return false, err
	}

	count := 0
	for _, f := range f {
		if f == constants.ManagerStateFilename {
			count++
		}
	}

	switch count {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		return false, fmt.Errorf("Found %d state files", len(f))
	}
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
	log.Debugf("Reading state from %s", constants.ManagerStateFilename)
	f, err := s.Backend.Read(constants.ManagerStateFilename)
	if err != nil {
		return nil, err
	}

	i = &Info{}
	err = i.Deserialize(f)
	return i, err
}

func (s *State) write(i *Info) error {
	f, err := i.Serialize()
	if err != nil {
		return err
	}
	if s.Config.DryRun {
		return nil
	}
	return s.Backend.Write(constants.ManagerStateFilename, f)
}

func (s *State) delete() error {
	if s.Config.DryRun {
		return nil
	}
	log.Debugf("Removing remote state %s", constants.ManagerStateFilename)
	return s.Backend.Delete(constants.ManagerStateFilename)
}

func (s *State) isManagerRelease(name string) bool {
	// check if the release name was explicitly set by flag
	if s.Config.Export.ReleaseName != "" && name == s.Config.Export.ReleaseName {
		return true
	}
	return false
}
