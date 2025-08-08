package version

import "gopkg.in/yaml.v3"

type Branch string

const (
	Branch202111 Branch = "ec202111"
	Branch202211 Branch = "ec202211_ecsonic"
)

type Version struct {
	Branch string `yaml:"branch"`
}

func UnmarshalVersion(in []byte) (*Version, error) {
	var version Version
	err := yaml.Unmarshal(in, &version)
	if err != nil {
		return nil, err
	}

	return &version, nil
}
