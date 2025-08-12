package uparse

import (
	"testing"
)

func TestParseYaml(t *testing.T) {
}

func TestFmtYaml2Path(t *testing.T) {
	_ = map[string]any{
		"Name": "jack",
		"Age":  11,
		"hobby": map[string]any{
			"sun": 1,
			"tom": "tom hobby11",
		},
	}

}

func TestParseYamlFromPath(t *testing.T) {
}
