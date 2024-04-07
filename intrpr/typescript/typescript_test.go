package typescript

import (
	"goutil/basic/gtest"
	"testing"
)

func TestTranspileToJavaScript(t *testing.T) {
	js, err := TranspileToJavaScript("let a: number = 10;", "")
	if js != "var a = 10;" {
		gtest.PrintlnExit(t, "typescript transpile failed")
	}
	gtest.Assert(t, err)
}
