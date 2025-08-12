package utype

import "testing"

func TestS2i(t *testing.T) {
	// 9223372036854775807
	// 243333333333322222222222222222222
	println(S2i64("243333333333322222222222222222222"))
}

func TestB2s(t *testing.T) {
	println(B2s(true))
}

func TestS2b(t *testing.T) {
	println(S2b("true"))
}
