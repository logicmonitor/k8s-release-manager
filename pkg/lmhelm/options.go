package lmhelm

import (
	"github.com/logicmonitor/k8s-release-manager/pkg/config"
	"k8s.io/helm/pkg/helm"
	rls "k8s.io/helm/pkg/proto/hapi/release"
)

func installOpts(r *rls.Release, vals []byte, config *config.HelmConfig) []helm.InstallOption {
	return []helm.InstallOption{
		helm.InstallReuseName(true),
		helm.InstallTimeout(config.ReleaseTimeoutSec),
		helm.InstallWait(true),
		helm.ReleaseName(r.GetName()),
		helm.ValueOverrides(vals),
	}
}

func listOpts() []helm.ReleaseListOption {
	return []helm.ReleaseListOption{
		helm.ReleaseListStatuses([]rls.Status_Code{
			rls.Status_DEPLOYED,
			rls.Status_FAILED,
		}),
	}
}
