package main

import (
	"fmt"
	"github.com/davidforest123/goutil/container/gvolume"
)

func main() {
	vol, err := gvolume.ParseString("10 MB")
	fmt.Println(vol.String(), err)
}
