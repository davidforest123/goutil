package gsniff

import "github.com/davidforest123/goutil/basic/gerrors"

// Usage multiplexing connection, dialer and listener.
// It uses the initial 8 bytes of each connection to indicate its purpose,
// in order to multiplex the Listener on one port to cope with clients of different purposes.

type (
	Usage [8]byte
)

var (
	UsageNone = MustMakeUsage("@nonono\n")
)

func MustMakeUsage(s string) Usage {
	usage, err := MakeUsage(s)
	if err != nil {
		panic(err)
	}
	return usage
}

func MakeUsage(s string) (Usage, error) {
	if len(s) != 8 {
		return [8]byte{}, gerrors.New("Usage must be 8 bytes, but got %d bytes: %s", len(s), s)
	}
	if s[0] != '@' || (s[7] != ':' && s[7] != '\n') {
		return [8]byte{}, gerrors.New("invalid Usage `%s`", s)
	}
	result := [8]byte{}
	copy(result[:], s)
	return result, nil
}

func (u Usage) Bytes() []byte {
	return u[:]
}

func (u Usage) String() string {
	return string(u[:])
}
