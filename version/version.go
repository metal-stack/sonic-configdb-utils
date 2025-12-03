package version

import (
	"strings"

	"gopkg.in/yaml.v3"
)

type Branch string

const (
	Branch202111 Branch = "ec202111"
	Branch202211 Branch = "ec202211"
)

type Version struct {
	Branch string `yaml:"branch"`
}

func GetVersion(in []byte) (Branch, error) {
	var version Version
	err := yaml.Unmarshal(in, &version)
	if err != nil {
		return "", err
	}

	if strings.Contains(version.Branch, string(Branch202111)) {
		return Branch202111, nil
	}
	if strings.Contains(version.Branch, string(Branch202211)) {
		return Branch202211, nil
	}

	return "", nil
}
