package release

import (
	"bytes"
	"fmt"
	"io"

	"github.com/golang/protobuf/proto"
	"github.com/logicmonitor/k8s-release-manager/pkg/constants"
	log "github.com/sirupsen/logrus"
	rls "k8s.io/helm/pkg/proto/hapi/release"
)

// ToString returns the string representation of the release
func ToString(r *rls.Release, verbose bool) (ret string, err error) {
	ret = fmt.Sprintf(`Name: %s
Filename: %s
Status: %s
Version: %d
Namespace: %s
Values: %s`,
		r.GetName(),
		Filename(r),
		r.GetInfo().GetStatus().GetCode().String(),
		r.GetVersion(),
		r.GetNamespace(),
		r.GetConfig().GetRaw(),
	)

	if verbose {
		ret += fmt.Sprintf("\nChart: %s", r.GetChart().String())
	}
	return ret, err
}

// Filename returns the calculated filename for the specified release
func Filename(r *rls.Release) string {
	info := r.GetInfo()
	if info == nil {
		log.Warnf("Unable to get info for release %s", r.GetName())
		return fmt.Sprintf("%s-%d.%s", r.GetName(), r.GetVersion(), constants.ReleaseExtension)
	}

	t := info.GetLastDeployed()
	return fmt.Sprintf("%s-%d-%d.%s", r.GetName(), r.GetVersion(), t.Seconds, constants.ReleaseExtension)
}

// FromFile returns a release struct from raw bytes
func FromFile(f []byte) (r *rls.Release, err error) {
	r = &rls.Release{}
	err = proto.Unmarshal(f, r)
	if err != nil {
		return nil, err
	}
	return r, err
}

// ToFile serializes the release into an io.Reader for writing to a filesystem
func ToFile(m proto.Message) (io.Reader, error) {
	f, err := proto.Marshal(m)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(f), nil
}

// func (m *Manager) waitForReleaseToDeploy(rls *lmhelm.Release) error {
// 	timeout := time.After(2 * time.Minute)
// 	ticker := time.NewTicker(30 * time.Second)

// 	for c := ticker.C; ; <-c {
// 		select {
// 		case <-timeout:
// 			return errors.New("Timed out waiting for release to deploy")
// 		default:
// 			log.Debugf("Checking status of release %s", rls.Name())
// 			if rls.Deployed() {
// 				return nil
// 			}
// 		}
// 	}
// }
