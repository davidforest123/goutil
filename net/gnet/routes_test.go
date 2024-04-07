package gnet

import (
	"fmt"
	"github.com/davidforest123/goutil/basic/gtest"
	"github.com/davidforest123/goutil/encoding/gjson"
	"testing"
)

func TestGetRoutes(t *testing.T) {
	routes, err := GetRoutes()
	gtest.Assert(t, err)

	fmt.Println(gjson.MarshalString(routes, true))
}
