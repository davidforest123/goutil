package main

import (
	"fmt"
	"github.com/davidforest123/goutil/basic/glog"
	"github.com/davidforest123/goutil/container/gstring"
	"github.com/davidforest123/goutil/sys/gfs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	goPath := os.Getenv("GOPATH")
	glog.AssertTrue(goPath != "")

	repoPath := "github.com/davidforest123/goutil"
	commpkgPath := filepath.Join(goPath, "src", repoPath)
	dirs, _, err := gfs.WalkDir(commpkgPath)
	glog.AssertOk(err, "WalkDir")

	res := map[string][]string{}
	for _, v := range dirs {
		ss := strings.Split(v, repoPath)
		if len(ss) != 2 {
			continue
		}
		subItems := strings.Split(ss[1], "/")
		subItems = gstring.RemoveByValue(subItems, "")
		if len(subItems) == 2 {
			original := res[subItems[0]]
			original = append(original, subItems[1])
			res[subItems[0]] = original
		}
	}

	pkgCount := 0
	for pkgSort, pkgs := range res {
		if pkgSort == ".git" {
			continue
		}
		fmt.Println("\n**" + pkgSort + "**\n")
		for _, pkg := range pkgs {
			pkgCount++
			fmt.Print(pkg + "  ")
		}
		fmt.Println("")
	}
	fmt.Println("\ntotal:", pkgCount)
}
