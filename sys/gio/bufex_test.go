package gio

import (
	"fmt"
	"strings"
	"testing"
)

func TestNewBufEx(t *testing.T) {
	r := strings.NewReader("abc")

	b2 := make([]byte, 1)
	n2, err2 := r.Read(b2)
	fmt.Println(b2, n2, err2)
}
