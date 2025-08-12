package ucmd

import (
	uos "github.com/hxxshidage/myutils/os"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"os/exec"
	"sync/atomic"
	"syscall"
)

var nameVal atomic.Value

func SetExecName(name string) bool {
	return nameVal.CompareAndSwap(nil, name)
}

// 直接执行命令
func ExecCmd(args ...string) (string, error) {
	cmd := exec.Command(getExecName(), append([]string{"-c"}, args...)...)

	hideWin(cmd)

	return doExecute(cmd, false)
}

// 执行shell脚本
// 通过git bash来执行的话, 脚本的路径可以是相对路劲
func ExecShell(shell string, opGbk bool, args ...string) (string, error) {
	cmd := exec.Command(getExecName(), append([]string{shell}, args...)...)

	hideWin(cmd)

	return doExecute(cmd, opGbk)
}

func doExecute(cmd *exec.Cmd, opGbk bool) (string, error) {
	// windows命令输出是gbk编码会导致乱码
	if opGbk && uos.PlatformWin() {
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return "", err
		}

		err = cmd.Start()
		if err != nil {
			return "", err
		}

		render := transform.NewReader(stdout, simplifiedchinese.GBK.NewDecoder())
		output, err := io.ReadAll(render)

		err = cmd.Wait()
		if err != nil {
			return "", err
		}

		return string(output), err
	} else {
		output, err := cmd.CombinedOutput()

		return string(output), err
	}
}

func getExecName() string {
	name := nameVal.Load()
	if name == nil {
		panic("Can't exec without name")
	}

	return name.(string)
}

// 调用shell, 隐藏shell窗口(弹会就消失)
func hideWin(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
}
