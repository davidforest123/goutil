package intrpr

import (
	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
	"goutil/basic/gerrors"
	"goutil/container/gany"
	"goutil/intrpr/typescript"
)

type (
	vmGoja struct {
		vm     *goja.Runtime
		prog   *goja.Program
		retVal gany.Val
	}
)

func valGoja2Comm(gojaVal goja.Value) any {
	return gojaVal.Export()
}

func newVMGoja() (*vmGoja, error) {
	res := &vmGoja{vm: goja.New()}
	res.vm.SetParserOptions(parser.WithDisableSourceMaps) // prevented filesystem access
	return res, nil
}

func (vm *vmGoja) toGojaValue(anyValue any) goja.Value {
	return vm.vm.ToValue(anyValue)
}

func (vm *vmGoja) LoadScript(lang Lang, script string) error {
	if lang != LangJavaScript && lang != LangTypeScript {
		return gerrors.New("engine yaegj doesn't support language %s", lang)
	}
	if lang == LangTypeScript {
		err := error(nil)
		script, err = typescript.TranspileToJavaScript(script, "")
		if err != nil {
			return err
		}
	}

	err := error(nil)
	vm.prog, err = goja.Compile("", script, true)
	return err
}

func (vm *vmGoja) LoadScriptAndRun(script string) (gany.Val, error) {
	val, err := vm.vm.RunString(script)
	if err != nil {
		return gany.ValNil, err
	}
	return gany.NewVal(val.Export()), nil
}

func (vm *vmGoja) SetVal(name string, value any) error {
	return vm.vm.Set(name, value)
}

func (vm *vmGoja) SetFunc(name string, fn any) error {
	return vm.vm.Set(name, fn)
}

func (vm *vmGoja) Run() error {
	if vm.vm == nil {
		return gerrors.New("goja vm not initialized")
	}
	if vm.prog == nil {
		return gerrors.New("goja prog not initialized")
	}
	vm.retVal = gany.ValNil
	val, err := vm.vm.RunProgram(vm.prog)
	if err != nil {
		return err
	}
	vm.retVal = gany.NewVal(val.Export())
	return nil
}

func (vm *vmGoja) GetFunc(name string) (Callable, error) {
	callable, ok := goja.AssertFunction(vm.vm.Get(name))
	if !ok {
		return nil, gerrors.New("%s is not a valid function", name)
	}

	return func(args ...any) ([]gany.Val, error) {
		var items []goja.Value
		for _, item := range args {
			items = append(items, vm.toGojaValue(item))
		}
		retGoja, err := callable(goja.Undefined(), items...)
		if err != nil {
			return nil, err
		}
		return []gany.Val{gany.NewVal(retGoja.Export())}, nil
	}, nil
}

func (vm *vmGoja) GetVal(name string) (gany.Val, error) {
	if name == "return" {
		return vm.retVal, nil
	}
	val := vm.vm.Get(name)
	return gany.NewVal(val.Export()), nil
}
