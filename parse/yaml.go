package uparse

import (
	"bytes"
	"gopkg.in/yaml.v3"
)

type v3yamlParser struct {
}

var (
	v3Yp = v3yamlParser{}
)

func (v3yamlParser) Parse(data []byte, val any) error {
	return yaml.Unmarshal(data, val)
}

func (v3yamlParser) Fmt(val any) ([]byte, error) {
	return yaml.Marshal(val)
}

func (v3yamlParser) FmtPretty(val any) ([]byte, error) {
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)

	if err := encoder.Encode(val); err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}

func ParseYaml(data []byte, val any) error {
	return v3Yp.Parse(data, val)
}

func FmtYaml(val any) ([]byte, error) {
	return v3Yp.Fmt(val)
}

func FmtYamlPretty(val any) ([]byte, error) {
	return v3Yp.FmtPretty(val)
}
