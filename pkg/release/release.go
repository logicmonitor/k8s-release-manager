package release

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/logicmonitor/k8s-release-manager/pkg/constants"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	rls "helm.sh/helm/v3/pkg/release"
)

// ToString returns the string representation of the release
func ToString(r *rls.Release, verbose bool) (ret string) {
	ret = fmt.Sprintf(`Name: %s
Filename: %s
Status: %s
Version: %d
Namespace: %s
Values:
%s`,
		r.Name,
		Filename(r),
		r.Info.Status.String(),
		r.Version,
		r.Namespace,
		getConfigYaml(r.Config),
	)

	if verbose {
		ret += fmt.Sprintf("\nChart: %s", getConfigYaml(r.Chart.Values))
	}
	return ret
}

func getConfigYaml(c map[string]interface{}) string {

	conf, err := yaml.Marshal(c)
	if err != nil {
		return ""
	}
	return string(conf)
}

// Filename returns the calculated filename for the specified release
func Filename(r *rls.Release) string {
	info := r.Info
	if info == nil {
		log.Warnf("Unable to get info for release %s", r.Name)
		return fmt.Sprintf("%s-%d.%s", r.Name, r.Version, constants.ReleaseExtension)
	}

	t := info.LastDeployed
	return fmt.Sprintf("%s-%d-%d.%s", r.Name, r.Version, t.Second(), constants.ReleaseExtension)
}

// FromFile returns a release struct from raw bytes
func FromFile(f []byte) (r *rls.Release, err error) {
	r = &rls.Release{}

	err = json.Unmarshal(f, r)
	if err != nil {
		return nil, err
	}
	return r, err
}

// ToFile serializes the release into an io.Reader for writing to a filesystem
func ToFile(r *rls.Release) (io.Reader, error) {

	f, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(f), nil
}

// UpdateValue configured in the release
func UpdateValue(r *rls.Release, name, value string) (*rls.Release, error) {

	config, err := updateConfigValues(r.Config, name, value)
	if err != nil {
		log.Errorf("Error updating values struct: %v", err)
	}

	r.Config = config
	return r, nil
}

func updateConfigValues(values map[string]interface{}, name, value string) (map[string]interface{}, error) {

	if values != nil {
		if _, ok := values[name]; !ok {
			log.Warnf("Value %s doesn't exist in the stored release, so refusing to update it", name)
			return values, nil
		}

	}

	values[name] = value

	return values, nil
}
