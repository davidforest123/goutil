package anycloud

import (
	"github.com/davidforest123/goutil/basic/gerrors"
	"github.com/davidforest123/goutil/container/gbytes"
	"github.com/davidforest123/goutil/container/gnum"
)

type (
	Msg struct {
		Cmd     CommCmd
		Subj    string
		DataLen gnum.Uint24
		Data    []byte
	}

	CommCmd uint8
)

var (
	CommCmdGossip     = CommCmd(1) // Gossip
	CommCmdPush       = CommCmd(2) // Push/Pop
	CommCmdPop        = CommCmd(3) // Push/Pop
	CommCmdPub        = CommCmd(4) // Pub/Sub
	CommCmdSub        = CommCmd(5) // Pub/Sub
	CommCmdRequest    = CommCmd(6) // Request/Reply
	CommCmdReply      = CommCmd(7) // Request/Reply
	CommCmdRegSub     = CommCmd(8)
	CommCmdRegPop     = CommCmd(9)
	CommCmdRegReply   = CommCmd(10)
	CommCmdUnRegSub   = CommCmd(11)
	CommCmdUnRegPop   = CommCmd(12)
	CommCmdUnRegReply = CommCmd(13)
)

func (m *Msg) totalLen() int {
	return 1 + (len(m.Subj) + 1) + 3 + int(m.DataLen.Uint32())
}

func msgEqual(msg1, msg2 *Msg) bool {
	if msg1.Cmd != msg2.Cmd {
		return false
	}
	if msg1.Subj != msg2.Subj {
		return false
	}
	if msg1.DataLen != msg2.DataLen {
		return false
	}
	if len(msg1.Data) != len(msg2.Data) {
		return false
	}
	for i := range msg1.Data {
		if msg1.Data[i] != msg2.Data[i] {
			return false
		}
	}
	return true
}

func NewMsg(mode CommCmd, subj string, data []byte) (*Msg, error) {
	if len(data) > gnum.MaxUint24 {
		return nil, gerrors.New("msg data length %d is bigger than max limit %d", len(data), gnum.MaxUint24)
	}

	msg := &Msg{
		Cmd:     mode,
		Subj:    subj,
		DataLen: gnum.NewUint24(uint32(len(data))),
		Data:    data,
	}
	return msg, nil
}

func EncodeMsg(msg *Msg) ([]byte, error) {
	msgBuf := make([]byte, msg.totalLen())
	next := gbytes.EncodeUint8(msgBuf, uint8(msg.Cmd))
	next = gbytes.EncodeString(next, msg.Subj)
	next = gbytes.EncodeUint24(next, msg.DataLen)
	copy(next, msg.Data)
	return msgBuf, nil
}

func DecodeMsg(msgBuf []byte) (*Msg, error) {
	msg := &Msg{}
	next := gbytes.DecodeUint8(msgBuf, (*uint8)(&msg.Cmd))
	next = gbytes.DecodeString(next, &msg.Subj)
	next = gbytes.DecodeUint24(next, &msg.DataLen)
	msg.Data = make([]byte, msg.DataLen.Uint32())
	cpyLen := copy(msg.Data, next)
	if cpyLen != int(msg.DataLen.Uint32()) {
		return nil, gerrors.New("cpyLen %d != msg.DataLen %d", cpyLen, msg.DataLen)
	}
	return msg, nil
}

func ReadMsgCmd(msgBuf []byte) CommCmd {
	return CommCmd(msgBuf[0])
}

func ReadMsgSubj(msgBuf []byte) string {
	for i := 1; i < len(msgBuf); i++ {
		if msgBuf[i] == 0 {
			return string(msgBuf[1:i])
		}
	}
	return ""
}

func UpdateMsgCmd(msgBuf []byte, newCmd CommCmd) {
	msgBuf[0] = byte(newCmd)
}
