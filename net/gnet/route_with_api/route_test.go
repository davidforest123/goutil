package route_with_api

import (
	"fmt"
	"testing"
)

func TestGetRoutesByAPI(t *testing.T) {
	data, err := getRoutesByAPI()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range data {
		fmt.Println(v.ToTableString())
	}
}
