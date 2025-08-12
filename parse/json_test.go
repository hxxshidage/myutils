package uparse

import (
	"fmt"
	"testing"
)

func TestParseJsonFromPath(t *testing.T) {
	var m map[string]any
	err := ParseJsonFromPath("j_test.json", &m)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", m)

}

func TestFmtJson2Path(t *testing.T) {
	var m map[string]any
	err := ParseJsonFromPath("j_test.json", &m)
	if err != nil {
		panic(err)
	}

	err = FmtJson2Path(&m, "j_test1.json", true)
	if err != nil {
		panic(err)
	}
}
