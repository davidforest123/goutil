// Listener is a Listener that listens on multiple addresses and ports,
// and it is used like one Listener.
// and it reads from net.Conn and detect what protocol it is.
// references:
// https://github.com/soheilhy/cmux/blob/master/matchers.go

package gsniff

import (
	"fmt"
	"goutil/basic/gerrors"
	"goutil/net/gnet"
	"goutil/sys/gio"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

type (
	mlAccepted struct {
		conn *Conn
		err  *Err
	}

	// listener is a sniffing listener when matchers are NOT empty.
	listener struct {
		network  string       // listening network
		addr     string       // listen address
		ln       net.Listener // raw listener like "tcp", "unix", "udp" ...
		chDie    chan struct{}
		matchers []Matcher
	}

	// Listener is a multiple networks / ports / protocols listener.
	Listener struct {
		lns       []*listener
		lnsMtx    sync.RWMutex
		chAccepts chan mlAccepted
		chDie     chan struct{}
		wgClose   *sync.WaitGroup
	}

	Err struct {
		network string // listening network
		addr    string // listen address
		err     error
	}

	// Info is sniffed info.
	Info struct {
		Network     string // listening network
		Addr        string // listen address
		Usage       Usage
		UsageExt    string
		Protocols   []string // "http1", "socks5", TODO "SMux", "http2", "websocket", "grpc"...
		HTTP1Req    *http.Request
		HTTP2Fields [][2]string
		Reasons     map[string]string
	}

	// Conn is server-side connection which could sniff info.
	Conn struct {
		net.Conn
		buf  *gio.SniffBuf
		info *Info
	}
)

const acceptBacklog = 4096

func (se *Err) Error() string {
	return se.err.Error()
}

// NewListener creates new Listener.
func NewListener() *Listener {
	return &Listener{
		chAccepts: make(chan mlAccepted, acceptBacklog),
		chDie:     make(chan struct{}),
	}
}

func (ml *Listener) listenRoutine(l *listener) {
	for {
		select {
		case <-ml.chDie:
			return
		case <-l.chDie:
			return
		default:
			rawConn, err := l.ln.Accept()
			rawAcceptErr := err

			// rebuild connection
			newConn := (*Conn)(nil)
			if rawConn != nil {
				if len(l.matchers) > 0 {
					newConn, err = trySniff(rawConn, l.matchers...)
				} else {
					newConn = &Conn{
						Conn: rawConn,
						buf:  gio.NewSniffBuf(rawConn),
						info: &Info{
							Network: l.network,
							Addr:    l.addr,
							Usage:   UsageNone,
						},
					}
				}
			}

			// rebuild error
			newErr := (*Err)(nil)
			if err != nil {
				newErr = &Err{
					network: l.network,
					addr:    l.addr,
					err:     err,
				}
			}

			// output connection and error
			ml.chAccepts <- mlAccepted{
				conn: newConn,
				err:  newErr,
			}

			if rawAcceptErr != nil {
				return
			}
		}
	}
}

// AddListener add new listener.
func (ml *Listener) AddListener(rawLn net.Listener, matchers ...Matcher) error {
	ln := &listener{
		network:  rawLn.Addr().Network(),
		addr:     rawLn.Addr().String(),
		ln:       rawLn,
		chDie:    make(chan struct{}),
		matchers: matchers,
	}

	ml.lnsMtx.Lock()
	ml.lns = append(ml.lns, ln)
	ml.lnsMtx.Unlock()

	go ml.listenRoutine(ln)
	return nil
}

// AddListen add new listen address.
func (ml *Listener) AddListen(network, addr string, matchers ...Matcher) error {
	rawLn, err := gnet.ListenAny(network, addr)
	if err != nil {
		return err
	}
	ln := &listener{
		network:  network,
		addr:     addr,
		ln:       rawLn,
		chDie:    make(chan struct{}),
		matchers: matchers,
	}

	ml.lnsMtx.Lock()
	ml.lns = append(ml.lns, ln)
	ml.lnsMtx.Unlock()

	go ml.listenRoutine(ln)
	return nil
}

// Accept waits for and returns the next connection to the listener.
func (ml *Listener) Accept() (net.Conn, error) {
	select {
	case <-ml.chDie:
		return nil, nil
	case accepted := <-ml.chAccepts:
		return accepted.conn, accepted.err
	}
}

// CloseOne closes one listener.
func (ml *Listener) CloseOne(network, addr string) error {
	err := error(nil)
	ml.lnsMtx.RLock()
	for _, v := range ml.lns {
		if v.network == network && v.addr == addr {
			close(v.chDie)
			err = v.ln.Close()
		}
	}
	ml.lnsMtx.RUnlock()
	return err
}

// Close closes all the listeners.
// Any blocked 'Accept' operations will be unblocked and return errors.
func (ml *Listener) Close() error {
	close(ml.chDie)

	var errs []error
	ml.lnsMtx.RLock()
	for _, v := range ml.lns {
		close(v.chDie)
		if err := v.ln.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	ml.lnsMtx.RUnlock()

	return gerrors.JoinArray(errs)
}

// Addr returns all the listener's network address.
func (ml *Listener) Addr() []string {
	var res []string

	ml.lnsMtx.RLock()
	for _, v := range ml.lns {
		res = append(res, fmt.Sprintf("%s://%s", v.network, v.addr))
	}
	ml.lnsMtx.RUnlock()

	return res
}

// trySniff trys to sniff server-side connection.
func trySniff(svrSideConn net.Conn, matchers ...Matcher) (*Conn, error) {
	if svrSideConn == nil {
		panic("conn can't be nil when Wrap")
	}
	if len(matchers) == 0 {
		return nil, gerrors.New("no matchers")
	}
	result := &Conn{
		Conn: svrSideConn,
		buf:  gio.NewSniffBuf(svrSideConn),
		info: &Info{
			Reasons: map[string]string{},
		},
	}

	rewindReader := result.buf.RewindReader()
	for _, matcher := range matchers {
		// set read timeout before execute matcher function
		err := svrSideConn.SetReadDeadline(time.Now().Add(matchReadTimeout))
		if err != nil {
			return nil, err
		}

		// try match
		rewindReader.Rewind()
		toDelBytes, err := matcher(rewindReader, result.info)
		if err != nil {
			return nil, err
		}

		// delete user level preset usage info
		if toDelBytes > 0 {
			io.ReadFull(result.buf.NormalReader(), make([]byte, toDelBytes))
		}

		if len(result.info.Protocols) > 0 {
			break
		}
	}

	return result, nil
}

// WARNING: data read by NormalReader can't be read from RewindReader again.
func (c *Conn) Read(p []byte) (n int, err error) {
	return c.buf.NormalReader().Read(p)
}

// SniffedInfo returns sniffed info.
func (c *Conn) SniffedInfo() *Info {
	return c.info
}
