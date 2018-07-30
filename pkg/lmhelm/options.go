package lmhelm

import (
	"k8s.io/helm/pkg/helm"
	rls "k8s.io/helm/pkg/proto/hapi/release"
)

// func installOpts(r *rls.Release, vals []byte, config *config.Config) []helm.InstallOption {
// 	return []helm.InstallOption{
// 		helm.InstallReuseName(true),
// 		helm.InstallTimeout(config.Client.ReleaseTimeoutSec),
// 		helm.InstallWait(true),
// 		helm.ReleaseName(r.Name),
// 		helm.ValueOverrides(vals),
// 	}
// }

func listOpts() []helm.ReleaseListOption {
	return []helm.ReleaseListOption{
		helm.ReleaseListStatuses([]rls.Status_Code{
			rls.Status_DEPLOYED,
			rls.Status_FAILED,
			rls.Status_PENDING_INSTALL,
			rls.Status_PENDING_ROLLBACK,
			rls.Status_PENDING_UPGRADE,
		}),
	}
}