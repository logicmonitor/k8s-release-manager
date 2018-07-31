package transfer

import (
	"fmt"
	"sync"

	"github.com/logicmonitor/k8s-release-manager/pkg/backend"
	"github.com/logicmonitor/k8s-release-manager/pkg/client"
	"github.com/logicmonitor/k8s-release-manager/pkg/config"
	"github.com/logicmonitor/k8s-release-manager/pkg/lmhelm"
	"github.com/logicmonitor/k8s-release-manager/pkg/release"
	"github.com/logicmonitor/k8s-release-manager/pkg/state"
	log "github.com/sirupsen/logrus"
	rls "k8s.io/helm/pkg/proto/hapi/release"
)

// Transferer deploys remotely stored releases
type Transferer struct {
	Config     *config.Config
	HelmClient *lmhelm.Client
	State      *state.State
}

// New instantiates and returns a Deleter and an error if any.
func New(rlsmgrconfig *config.Config, backend backend.Backend) (*Transferer, error) {
	helmClient := &lmhelm.Client{}

	// dry run's don't need to interact with tiller, so skip config setup
	if !rlsmgrconfig.DryRun {
		kubernetesClient, kubernetesConfig, err := client.KubernetesClient(rlsmgrconfig.ClusterConfig)
		if err != nil {
			return nil, err
		}

		err = helmClient.Init(rlsmgrconfig.Helm, kubernetesClient, kubernetesConfig)
		if err != nil {
			return nil, err
		}
	}

	return &Transferer{
		Config:     rlsmgrconfig,
		HelmClient: helmClient,
		State: &state.State{
			Backend: backend,
			Config:  rlsmgrconfig,
		},
	}, nil
}

// Run the transferer.
func (t *Transferer) Run() error {
	if t.Config.DryRun {
		fmt.Println("Dry run. No changes will be made.")
	}

	releases, err := t.State.StoredReleases()
	if err != nil {
		log.Fatalf("Error retrieving stored releases: %v", err)
	}

	return t.deployReleases(releases)
}

func (t *Transferer) deployReleases(releases []*rls.Release) error {
	var stateInfo *state.Info
	stateExists, err := t.State.Exists()
	if err != nil {
		log.Warnf("%v", err)
	}

	if stateExists {
		stateInfo, err = t.State.Read()
		if err != nil {
			log.Errorf("Error retrieving remote state: ")
		}
	}

	var wg sync.WaitGroup
	for _, r := range releases {
		fmt.Printf("Deploying release: %s\n", r.GetName())

		// check if this release is managing the release manager
		if t.Config.Transfer.NewStoragePath != "" && stateExists && stateInfo != nil && r.GetName() == stateInfo.ReleaseName {
			r, err = t.updateManagerStoragePath(r, t.Config.Transfer.NewStoragePath)
			if err != nil {
				log.Errorf("Unable to update the output path for the new release manager chart. Skipping.")
				continue
			}
		}

		if t.Config.DryRun {
			fmt.Printf("%s\n", release.ToString(r, t.Config.VerboseMode))
			continue
		}

		wg.Add(1)
		go func(r *rls.Release, client *lmhelm.Client) {
			defer wg.Done()
			e := client.Install(r)
			if e != nil {
				log.Errorf("Error deploying release %s: %v", r.GetName(), e)
			}
		}(r, t.HelmClient)
	}
	wg.Wait()
	return nil
}

func (t *Transferer) updateManagerStoragePath(r *rls.Release, path string) (*rls.Release, error) {
	// TODO don't use hardcoded path. how do we handle non-official charts? i guess we don't
	return release.UpdateValue(r, "backend.storagePath", path)
}
