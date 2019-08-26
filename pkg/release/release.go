package release

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/logicmonitor/k8s-release-manager/pkg/constants"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"k8s.io/helm/pkg/proto/hapi/chart"
	rls "k8s.io/helm/pkg/proto/hapi/release"
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
	return ret
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

// UpdateValue configured in the release
func UpdateValue(r *rls.Release, name string, value string) (*rls.Release, error) {
	// it's unclear how helm is utilizing config.Raw() vs config.Values
	// to be safe, i'm just going to attempt to update both and see what happens
	var err error
	config := r.GetConfig()
	values := config.GetValues()
	raw := config.GetRaw()

	config.Values, err = updateValueStruct(values, name, value)
	if err != nil {
		log.Errorf("Error updating values struct: %v", err)
	}

	config.Raw, err = updateValueRaw(raw, name, value)
	if err != nil {
		log.Errorf("Error updating raw values: %v", err)
	}

	r.Config = config
	return r, nil
}

func updateValueStruct(values map[string]*chart.Value, name string, value string) (map[string]*chart.Value, error) {
	if values != nil {
		if _, ok := values[name]; !ok {
			log.Warnf("Value %s doesn't exist in the stored release, so refusing to update it", name)
			return values, nil
		}

		// again, it's not exactly clear if/how helm/tiller is using config.Values,
		// so i'm just kind of assuming that the key should be the full value
		// path string and assigning as such. This may not end up being correct,
		// but given that config.Values is map[string]*Value, i can't possibly see
		// how it could be nested.
		values[name] = &chart.Value{
			Value: value,
		}
	}
	return values, nil
}

func updateValueRaw(raw string, name string, value string) (string, error) {
	var y map[interface{}]interface{}
	err := yaml.Unmarshal([]byte(raw), &y)
	if err != nil {
		return "", err
	}

	path := strings.Split(name, ".")
	y, err = updateNestedMapString(y, path, value)
	if err != nil {
		return "", err
	}

	r, err := yaml.Marshal(y)
	if err != nil {
		return "", err
	}
	return string(r[:]), nil
}

func updateNestedMapString(m map[interface{}]interface{}, path []string, value string) (map[interface{}]interface{}, error) {
	var err error
	if len(path) < 1 {
		return m, nil
	}

	key := path[0]
	_, exists := m[key]
	switch true {
	case exists && len(path) == 1:
		m[key] = value
		return m, nil
	case exists:
		m[key], err = updateNestedMapString(m[key].(map[interface{}]interface{}), path[1:], value)
		return m, err
	default:
		log.Warnf("Value %s doesn't exist in the stored release, so refusing to update it", key)
		return m, nil
	}
}
