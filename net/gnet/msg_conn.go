package gnet

import (
	"goutil/basic/gerrors"
	"goutil/container/gbytes"
	"goutil/container/gnum"
	"io"
	"net"
	"unsafe"
)

type (
	msgLenType uint16

	MsgConn struct {
		net.Conn
		cachedMsgLen *msgLenType
	}

	msgModel struct {
		MsgLen msgLenType
		buf    []byte
	}
)

var (
	msgLenSize = int(unsafe.Sizeof(msgModel{}.MsgLen))
	msgLenMax  = msgLenType(gnum.MustMax(uint16(0)).(uint16))
)

func NewMsgConn(streamConn net.Conn) *MsgConn {
	return &MsgConn{Conn: streamConn}
}

func (mc *MsgConn) Raw() net.Conn {
	return mc.Conn
}

func (mc *MsgConn) ReadMsg() ([]byte, error) {
	msgLen := msgLenType(0)
	if mc.cachedMsgLen != nil {
		msgLen = *mc.cachedMsgLen
		mc.cachedMsgLen = nil
	} else {
		msgLenBuf := make([]byte, msgLenSize)
		_, err := io.ReadFull(mc.Conn, msgLenBuf)
		if err != nil {
			return nil, err
		}
		msgLenIf, err := gbytes.BytesToNum(msgLenBuf, msgLenType(0))
		if err != nil {
			return nil, err
		}
		msgLen = msgLenIf.(msgLenType)
	}
	if msgLen == 0 {
		return nil, gerrors.New("msgLen can't be zero")
	}
	result := make([]byte, msgLen)
	_, err := io.ReadFull(mc.Conn, result)
	return result, err
}

func (mc *MsgConn) Read(b []byte) (int, error) {
	msgLen := msgLenType(0)
	if mc.cachedMsgLen != nil {
		msgLen = *mc.cachedMsgLen
		mc.cachedMsgLen = nil
	} else {
		msgLenBuf := make([]byte, msgLenSize)
		_, err := io.ReadFull(mc.Conn, msgLenBuf)
		if err != nil {
			return 0, err
		}
		msgLenIf, err := gbytes.BytesToNum(msgLenBuf, msgLenType(0))
		if err != nil {
			return 0, err
		}
		msgLen = msgLenIf.(msgLenType)
	}
	if msgLen == 0 {
		return 0, gerrors.New("msgLen can't be zero")
	}
	if len(b) < int(msgLen) {
		mc.cachedMsgLen = &msgLen
		return 0, gerrors.New("len(b) %d < msgLen %d", len(b), msgLen)
	}
	return io.ReadFull(mc.Conn, b[:msgLen])
}

func (mc *MsgConn) Write(b []byte) (int, error) {
	if uint64(len(b)) > uint64(msgLenMax) {
		return 0, gerrors.New("len(b) %d > msgLenSize")
	}
	msgLenBuf, err := gbytes.NumToBytes(msgLenType(len(b)))
	if err != nil {
		return 0, err
	}
	if len(msgLenBuf) != int(msgLenSize) {
		return 0, gerrors.New("len(msgLenBuf) %d != msgLenSize %d", len(msgLenBuf), msgLenSize)
	}
	_, err = mc.Conn.Write(msgLenBuf)
	if err != nil {
		return 0, err
	}
	return mc.Conn.Write(b)
}
