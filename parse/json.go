package uparse

import (
	"encoding/json"
	"sync/atomic"
)

type nativeJsonParser struct {
}

var (
	nativeJp = nativeJsonParser{}
	setting  atomic.Bool
	currJp   DataParser = nativeJp
)

func SetJsonParser(tag string) {
	if !setting.CompareAndSwap(false, true) {
		return
	}

	if tag == "native" {
		currJp = nativeJp
	} else {

	}
}

func (nativeJsonParser) Parse(data []byte, val any) error {
	return json.Unmarshal(data, val)
}

func (nativeJsonParser) Fmt(val any) ([]byte, error) {
	return json.Marshal(val)
}

func (nativeJsonParser) FmtPretty(val any) ([]byte, error) {
	if contents, err := json.MarshalIndent(val, "", "\t"); err != nil {
		return nil, err
	} else {
		return contents, nil
	}
}

func ParseJson(data []byte, val any) error {
	return currJp.Parse(data, val)
}

func FmtJson(val any) ([]byte, error) {
	return currJp.Fmt(val)
}

func FmtJsonPretty(val any) ([]byte, error) {
	return currJp.FmtPretty(val)
}
