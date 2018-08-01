package delete

import (
	"fmt"

	"github.com/logicmonitor/k8s-release-manager/pkg/config"
	"github.com/logicmonitor/k8s-release-manager/pkg/release"
	"github.com/logicmonitor/k8s-release-manager/pkg/state"
	log "github.com/sirupsen/logrus"
)

// Delete deletes remotely stored release info
type Delete struct {
	Config *config.Config
	State  *state.State
}

// New instantiates and returns a Delete and an error if any.
func New(rlsmgrconfig *config.Config, state *state.State) (*Delete, error) {
	return &Delete{
		Config: rlsmgrconfig,
		State:  state,
	}, nil
}

// Run the Delete.
func (d *Delete) Run() error {
	releaseNames, err := d.State.Releases.StoredReleaseNames()
	if err != nil {
		log.Fatalf("Error retrieving stored releases: %v", err)
	}

	err = d.deleteReleases(releaseNames)
	if err != nil {
		log.Warnf("%v", err)
	}
	return d.deleteState()
}

func (d *Delete) deleteReleases(releaseNames []string) error {
	for _, f := range releaseNames {
		fmt.Printf("Removing release: %s\n", f)
		switch true {
		case d.Config.DryRun:
			r, e := d.State.Releases.ReadRelease(f)
			if e != nil {
				log.Errorf("Error retrieving remote release %s: %v", f, e)
			}
			if d.Config.DebugMode {
				fmt.Printf("%s\n", release.ToString(r, d.Config.VerboseMode))
			}
			continue
		default:
			e := d.State.Releases.DeleteRelease(f)
			if e != nil {
				log.Errorf("Error removing remote release %s: %v", f, e)
				continue
			}
		}
	}
	return nil
}

func (d *Delete) deleteState() error {
	return d.State.Remove()
}
