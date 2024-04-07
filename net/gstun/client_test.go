package gstun

import (
	"fmt"
	"testing"
)

func TestClient_Dial(t *testing.T) {
	fmt.Println(Discover("stun.l.google.com:19302"))
}
