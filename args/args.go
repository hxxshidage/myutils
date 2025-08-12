package uargs

import (
	"flag"
	utype "myutils/type"
)

type CmdArgTrans struct {
	Name   string
	DefVal string
}

type CmdParser struct {
	resultMap map[string]string
}

func NewCmdParser() *CmdParser {
	return &CmdParser{}
}

// --env=prod or -env prod
func (cp *CmdParser) ParseCmdArgs(args []CmdArgTrans) {
	argsPtr := make([]*string, len(args))
	for idx, arg := range args {
		argsPtr[idx] = flag.String(arg.Name, arg.DefVal, arg.Name)
	}

	flag.Parse()

	argName2value := make(map[string]string, len(args))
	for idx, ptr := range argsPtr {
		argName2value[args[idx].Name] = *ptr
	}

	cp.resultMap = argName2value
}

func (cp *CmdParser) GetIntArg(name string) int {
	return utype.S2i(cp.resultMap[name])
}

func (cp *CmdParser) GetStrArg(name string) string {
	return cp.resultMap[name]
}

func (cp *CmdParser) GetBoolArg(name string) bool {
	return utype.S2b(cp.resultMap[name])
}
