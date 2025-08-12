package uparse

type DataParser interface {
	Parse(data []byte, val any) error

	Fmt(val any) ([]byte, error)

	FmtPretty(val any) ([]byte, error)
}
