package uos

import "runtime"

var curOs = runtime.GOOS

func GetOs() string {
	return curOs
}

func PlatformWin() bool {
	return "windows" == curOs
}

func PlatformUnix() bool {
	return PlatformLinux() || PlatformMacOs()
}

func PlatformLinux() bool {
	return "linux" == curOs
}

func PlatformMacOs() bool {
	return "darwin" == curOs
}
