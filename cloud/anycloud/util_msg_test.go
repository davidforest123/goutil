package anycloud

import (
	"github.com/davidforest123/goutil/basic/gerrors"
	"github.com/davidforest123/goutil/basic/gtest"
	"testing"
)

func TestEncodeMsgDecodeMsg(t *testing.T) {
	msg, err := NewMsg(CommCmdPub, "234", []byte{5, 6, 7, 8})
	gtest.Assert(t, err)
	msgBuf, err := EncodeMsg(msg)
	gtest.Assert(t, err)
	decMsg, err := DecodeMsg(msgBuf)
	gtest.Assert(t, err)
	if !msgEqual(msg, decMsg) {
		gtest.Assert(t, gerrors.New("original msg != decoded msg"))
	}
}
