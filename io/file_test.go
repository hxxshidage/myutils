package uio

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestWinPath2LinuxPath(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{`C:\Windows\System`, `/c/Windows/System`},
		{`D:\Program Files`, `/d/Program Files`},
		{`\\server\share\file`, `//server/share/file`},
		{`/already/unix/path`, `/already/unix/path`},
		{`relative\path`, `relative/path`},
		{`C:\`, `/c`},
		{``, ``},
		{`C:\Mixed/Slash\Path`, `/c/Mixed/Slash/Path`},
	}

	for _, test := range tests {
		got := ToUnixPath(test.input)
		if got != test.expect {
			t.Errorf("winPath2LinuxPath(%q) = %q, want %q", test.input, got, test.expect)
		}
	}

}

func TestToUnixPath(t *testing.T) {
	//println(ToUnixPath("C:\\Program Files\\JetBrains\\IntelliJ IDEA 2022.2\\plugins\\maven\\lib\\maven3\\conf\\settings.xml"))
	//
	//println(ToUnixPath("C:\\Program Files\\JetBrains\\IntelliJ IDEA 2022.2\\plugins\\maven\\lib\\maven3"))

	ts := filepath.ToSlash("C:\\Users\\hq\\Desktop\\fessionx\\i18n\\i18n_in.xlsx")
	lastIdx := strings.LastIndex(ts, "/")
	dir := ts[:lastIdx]
	println(dir)
	println(filepath.FromSlash(dir))

	ts = "/C/des/i18.xls"
	println(filepath.FromSlash(filepath.Dir(ts)))

}
