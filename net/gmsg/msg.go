package gmsg

import (
	"github.com/davidforest123/goutil/basic/gerrors"
	"github.com/davidforest123/goutil/container/gbytes"
	"github.com/davidforest123/goutil/container/gnum"
	"io"
	"net"
	"time"
)

type (
	msg struct {
		head    byte
		dataLen gnum.Uint24 // message total length include `len` and `data`
		data    []byte      // payload
	}

	connMsg struct {
		streamConn net.Conn
		connName   string
	}

	Conn interface {
		Name() string
		Read() ([]byte, error)
		Write(b []byte) error
		Close() error
		LocalAddr() net.Addr
		RemoteAddr() net.Addr
		SetDeadline(t time.Time) error
		SetReadDeadline(t time.Time) error
		SetWriteDeadline(t time.Time) error
	}
)

const (
	msgHead = byte(0xF2)
)

func msgEqual(msg1, msg2 *msg) bool {
	if msg1.head != msg2.head {
		return false
	}
	if msg1.dataLen.Uint32() != msg2.dataLen.Uint32() {
		return false
	}
	if len(msg1.data) != len(msg2.data) {
		return false
	}
	for i := range msg1.data {
		if msg1.data[i] != msg2.data[i] {
			return false
		}
	}
	return true
}

func newMsg(data []byte) (*msg, error) {
	if len(data) > gnum.MaxUint24 {
		return nil, gerrors.New("msg data length %d is bigger than max limit %d", len(data), gnum.MaxUint24)
	}

	msg := &msg{
		head:    msgHead,
		dataLen: gnum.NewUint24(uint32(len(data))),
		data:    data,
	}
	return msg, nil
}

func (m *msg) totalLen() int {
	return 1 + 3 + int(m.dataLen.Uint32())
}

func encodeMsg(msg *msg) ([]byte, error) {
	buf := make([]byte, msg.totalLen())
	next := gbytes.EncodeUint8(buf, msg.head)
	next = gbytes.EncodeUint24(next, msg.dataLen)
	copy(next, msg.data)
	return buf, nil
}

func decodeMsg(msgBuf []byte) (*msg, error) {
	msg := &msg{}
	next := gbytes.DecodeUint8(msgBuf, &msg.head)
	next = gbytes.DecodeUint24(msgBuf, &msg.dataLen)
	if msg.totalLen() != len(msgBuf) {
		return nil, gerrors.New("msgBuf len %d != msg.totalLen %d", len(msgBuf), msg.totalLen())
	}
	msg.data = next[:]
	return msg, nil
}

func ClientConn(s net.Conn, connName string) (Conn, error) {
	c := &connMsg{streamConn: s}
	if err := c.Write([]byte(connName)); err != nil {
		return nil, err
	}
	return c, nil
}

func ServerConn(s net.Conn) (Conn, error) {
	c := &connMsg{streamConn: s}
	firstMsg, err := c.Read()
	if err != nil {
		return nil, err
	}
	c.connName = string(firstMsg)
	return c, nil
}

func NewConn(s net.Conn) Conn {
	return &connMsg{streamConn: s}
}

func (c *connMsg) Name() string {
	return c.connName
}

func (c *connMsg) Read() ([]byte, error) {
	// read msg head
	headBuf := make([]byte, 1) // uint8 requires 1 byte
	n, err := io.ReadFull(c.streamConn, headBuf)
	if err != nil {
		return nil, err
	}
	if n != 1 {
		return nil, gerrors.New("read msg head byte %d != 1", n)
	}
	head := uint8(0)
	gbytes.DecodeUint8(headBuf, &head)

	// read msg len
	dataLenBuf := make([]byte, 3) // uint24 requires 3 bytes
	n, err = io.ReadFull(c.streamConn, dataLenBuf)
	if err != nil {
		return nil, err
	}
	if n != 3 {
		return nil, gerrors.New("read msg length bytes %d != 2", n)
	}
	dataLen := gnum.Uint24{}
	gbytes.DecodeUint24(dataLenBuf, &dataLen)

	// read msg data
	dataBuf := make([]byte, dataLen.Uint32())
	n, err = io.ReadFull(c.streamConn, dataBuf)
	if err != nil {
		return nil, err
	}
	if n != int(dataLen.Uint32()) {
		return nil, gerrors.New("read msg bytes %d != %d", n, dataLen)
	}

	return dataBuf, nil
}

func (c *connMsg) Write(data []byte) error {
	maxUint16, err := gnum.Max(uint16(0))
	if err != nil {
		return err
	}
	if len(data) > int(maxUint16.(uint16)) {
		return gerrors.New("msg data length %d is bigger than max limit %d", len(data), maxUint16.(uint16))
	}

	msg, err := newMsg(data)
	if err != nil {
		return err
	}
	msgBuf, err := encodeMsg(msg)
	if err != nil {
		return err
	}

	wTotal := 0
	for wTotal < len(msgBuf) {
		wOnce, err := c.streamConn.Write(msgBuf[wTotal:])
		if err != nil {
			return err
		}
		wTotal += wOnce
	}
	return nil
}

func (c *connMsg) LocalAddr() net.Addr {
	return c.streamConn.LocalAddr()
}

func (c *connMsg) RemoteAddr() net.Addr {
	return c.streamConn.RemoteAddr()
}

func (c *connMsg) SetDeadline(t time.Time) error {
	return c.streamConn.SetDeadline(t)
}

func (c *connMsg) SetReadDeadline(t time.Time) error {
	return c.streamConn.SetReadDeadline(t)
}

func (c *connMsg) SetWriteDeadline(t time.Time) error {
	return c.streamConn.SetWriteDeadline(t)
}

func (c *connMsg) Close() error {
	return c.streamConn.Close()
}
