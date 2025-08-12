package uparse

import (
	"fmt"
	"testing"
)

func TestParseYaml(t *testing.T) {
	var m map[string]any
	err := ParseYamlFromPath("test.yaml", &m)
	if err != nil {
		panic(err)
	}
}

func TestFmtYaml2Path(t *testing.T) {
	m := map[string]any{
		"Name": "jack",
		"Age":  11,
		"hobby": map[string]any{
			"sun": 1,
			"tom": "tom hobby11",
		},
	}

	err := FmtYaml2Path(&m, "test1.yaml")
	if err != nil {
		panic(err)
	}
}

func TestParseYamlFromPath(t *testing.T) {
	var m map[string]any
	err := ParseYamlFromPath("test1.yaml", &m)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", m)
}
