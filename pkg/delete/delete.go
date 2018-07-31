package delete

import (
	"fmt"

	"github.com/logicmonitor/k8s-release-manager/pkg/backend"
	"github.com/logicmonitor/k8s-release-manager/pkg/config"
	"github.com/logicmonitor/k8s-release-manager/pkg/release"
	"github.com/logicmonitor/k8s-release-manager/pkg/state"
	log "github.com/sirupsen/logrus"
)

// Deleter deletes remotely stored release info
type Deleter struct {
	Config *config.Config
	State  *state.State
}

// New instantiates and returns a Deleter and an error if any.
func New(rlsmgrconfig *config.Config, backend backend.Backend) (*Deleter, error) {
	return &Deleter{
		Config: rlsmgrconfig,
		State: &state.State{
			Backend: backend,
			Config:  rlsmgrconfig,
		},
	}, nil
}

// Run the deleter.
func (d *Deleter) Run() error {
	if d.Config.DryRun {
		fmt.Println("Dry run. No changes will be made.")
	}

	releaseNames, err := d.State.StoredReleaseNames()
	if err != nil {
		log.Fatalf("Error retrieving stored releases: %v", err)
	}

	err = d.deleteReleases(releaseNames)
	if err != nil {
		log.Warnf("%v", err)
	}
	return d.deleteState()
}

func (d *Deleter) deleteReleases(releaseNames []string) error {
	for _, f := range releaseNames {
		fmt.Printf("Removing release: %s\n", f)
		switch true {
		case d.Config.DryRun:
			r, e := d.State.ReadRelease(f)
			if e != nil {
				log.Errorf("Error retrieving remote release %s: %v", f, e)
			} else {
				fmt.Printf("%s\n", release.ToString(r, d.Config.VerboseMode))
			}
			continue
		default:
			e := d.State.DeleteRelease(f)
			if e != nil {
				log.Errorf("Error removing remote release %s: %v", f, e)
				continue
			}
		}
	}
	return nil
}

func (d *Deleter) deleteState() error {
	if d.Config.DryRun {
		return nil
	}
	return d.State.Remove()
}
