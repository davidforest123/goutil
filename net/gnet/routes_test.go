package gnet

import (
	"fmt"
	"goutil/basic/gtest"
	"goutil/encoding/gjson"
	"testing"
)

func TestGetRoutes(t *testing.T) {
	routes, err := GetRoutes()
	gtest.Assert(t, err)

	fmt.Println(gjson.MarshalString(routes, true))
}
