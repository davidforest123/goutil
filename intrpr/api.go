package intrpr

import (
	"goutil/basic/gerrors"
	"goutil/container/gany"
)

type (
	// Callable represents a JavaScript function that can be called from Go.
	Callable func(args ...any) ([]gany.Val, error)

	Lang string

	// Vm is script interpreter.
	Vm interface {
		// LoadScript loads script and do some pre-run preparation,
		// such as transpile, pre-compiling, etc., if the runtime can.
		LoadScript(lang Lang, script string) error

		// SetVal maps golang value into script, it allows script to access go runtime `value` with `name`.
		// Usually `value` is a pointer in the go runtime, that means user can call member functions of `value` from script.
		SetVal(name string, value any) error

		// SetFunc maps golang function into script, it allows script to access go runtime function with `fn`.
		SetFunc(name string, fn any) error

		// Run runs script.
		Run() error

		// GetFunc maps script function to golang runtime, it allows go runtime to access script function.
		GetFunc(name string) (Callable, error)

		// GetVal returns value by variable name, if you want returned value, use `return`.
		GetVal(name string) (gany.Val, error)
	}
)

var (
	LangJavaScript = Lang("javascript")
	LangTypeScript = Lang("typescript")
	LangGo         = Lang("go")
	LangTengo      = Lang("tengo")
)

// NewVM creates interpreter.
// It supports Golang script and ECMAScript languages like Javascript, TypeScripts.
func NewVM(engine string) (Vm, error) {
	switch engine {
	case "goja":
		return newVMGoja()
	case "yaegi":
		return newVMYaegi()
	case "tengo":
		return newVMTengo()
	default:
		return nil, gerrors.New("unsupported interpreter engine %s", engine)
	}
}
