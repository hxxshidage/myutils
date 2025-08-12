package uguc

import (
	"fmt"
	"github.com/labstack/gommon/log"
	urtime "myutils/runtime"
	"runtime/debug"
)

type Runnable func() error

type PanicHandler func(err any) bool

type ErrorHandler func(err error)

func GoRun(run Runnable, ph PanicHandler, eh ErrorHandler) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				if ph == nil || !ph(r) {
					//buf := make([]byte, 4096)
					//layers := runtime.Stack(buf, false)
					bytes := debug.Stack()

					errStack := string(bytes)

					gid := urtime.GoId()

					errInfo := fmt.Sprintf("ErrorStackInGucRun Begin(%d):\n%s\nErrorStackInGucRun End(%d)\n", gid, errStack, gid)

					log.Error(errInfo)
				}
			}
		}()

		err := run()

		if err != nil {
			if eh != nil {
				eh(err)
			} else {
				log.Errorf("execute async runnable error occurred:%v", err)
			}
		}
	}()
}

func PositiveRun(run Runnable, ph PanicHandler) {
	GoRun(run, ph, nil)
}

func CatchRun(run Runnable, eh ErrorHandler) {
	GoRun(run, nil, eh)
}

func JustRun(run Runnable) {
	PositiveRun(run, nil)
}
