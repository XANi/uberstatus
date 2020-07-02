package mqtt

import (
	"fmt"
	"github.com/XANi/uberstatus/util"
	"github.com/glycerine/zygomys/zygo"
	"strconv"
)


func AddZygoFuncs(z *zygo.Zlisp) {
	z.AddFunction("bar", BarFunc)
}

func BarFunc(env *zygo.Zlisp, name string, args []zygo.Sexp) (zygo.Sexp, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("%s expects 1 argument", name)
	}
	var pct int
	switch v := args[0].(type) {
	case *zygo.SexpInt:
		pct = int(v.Val)

	case *zygo.SexpStr:
		p, err := strconv.Atoi(v.S)
		if err != nil {
			return nil, err
		}
		pct = p
	case *zygo.SexpRaw:
	default:
		return nil, fmt.Errorf("%s does not handle %+v type", name, args[0])
	}

	return &zygo.SexpStr{S:util.GetBarChar(pct)}, nil
	env.Run()
}