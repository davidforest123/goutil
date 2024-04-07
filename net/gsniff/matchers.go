package gsniff

import (
	"bufio"
	"fmt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
	"goutil/basic/gerrors"
	"goutil/net/gsocks5/socks5internal"
	"goutil/sys/gio"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type (
	// Matcher is match method implement.
	Matcher = func(r io.Reader, out *Info) (toDelBytes int, err error)
)

var (
	matchReadTimeout = time.Millisecond * 20

	AllMatchers = []Matcher{MatcherUsage, MatcherSocks5, MatcherHTTP2, MatcherHTTP2}
)

func MatcherUsage(r io.Reader, info *Info) (int, error) {
	// read usage
	info.Usage = UsageNone
	usageBuf := make([]byte, 8)
	n, err := gio.ReadFull(r, usageBuf, &matchReadTimeout)
	if err != nil {
		info.Reasons["usage"] = gerrors.Wrap(err, "ReadFull").Error()
		return 0, nil
	}
	if n != 8 {
		info.Reasons["usage"] = "can't read 8 bytes"
		return 0, nil
	}
	usage, err := MakeUsage(string(usageBuf))
	if err != nil {
		info.Reasons["usage"] = err.Error()
		return 0, nil
	}
	info.Usage = usage

	// read usage extended info
	extLen := 0
	if usage[7] != '\n' {
		keepN := false
		extInfo, err := gio.ReadLine(r, keepN, 1024)
		if err != nil {
			info.Reasons["usage"] = gerrors.Wrap(err, "ReadLine").Error()
			return 0, nil
		}
		info.UsageExt = string(extInfo)
		extLen = len(extInfo)
		if !keepN {
			extLen++
		}
	}

	// caller needs to delete user level preset usage info
	return len(UsageNone) + extLen, nil
}

// tested
func MatcherSocks5(r io.Reader, info *Info) (int, error) {
	// first frame: version
	version := make([]byte, 1)
	if _, err := gio.ReadFull(r, version, &matchReadTimeout); err != nil {
		info.Reasons["socks5"] = gerrors.Wrap(err, "ReadFull").Error()
		return 0, nil
	}
	if version[0] != socks5internal.Version {
		info.Reasons["socks5"] = fmt.Sprintf("invalid socks5 version %d", version[0])
		return 0, nil
	}

	// first frame: auth
	authMethodLen := make([]byte, 1)
	if _, err := gio.ReadFull(r, authMethodLen, &matchReadTimeout); err != nil {
		info.Reasons["socks5"] = gerrors.Wrap(err, "ReadFull").Error()
		return 0, nil
	}
	if authMethodLen[0] <= 0 || authMethodLen[0] > 3 {
		return 0, gerrors.New("current support max auth method count is 3, but got %d", authMethodLen[0])
	}
	authMethods := make([]byte, authMethodLen[0])
	if _, err := gio.ReadFull(r, authMethods, &matchReadTimeout); err != nil {
		info.Reasons["socks5"] = gerrors.Wrap(err, "ReadFull").Error()
		return 0, nil
	}
	for _, v := range authMethods {
		amv := socks5internal.AuthMethod(v)
		if amv != socks5internal.AuthMethodNoAuth && amv != socks5internal.AuthMethodGssapi && amv != socks5internal.AuthMethodPassword {
			info.Reasons["socks5"] = fmt.Sprintf("invalid auth method %d", v)
			return 0, nil
		}
	}

	// The first frame sent by the client is called auth frame. After sending it, the client must wait for the response
	// from the server, and only after receiving the response will the client send the second frame,
	// called command frame, which contains the access destination address. So when we sniff, we cannot read the socks5
	// destination address before the socks5 proxy server handles the connection.

	info.Protocols = append(info.Protocols, "socks5")
	return 0, nil
}

// tested except websocket
func MatcherHTTP1(r io.Reader, info *Info) (int, error) {
	// try best to read data until 0 or error is returned.
	buf, readErr := gio.ReadHead(r, &matchReadTimeout)

	// process no data error
	if buf.Len() == 0 {
		if readErr != nil {
			info.Reasons["http1"] = gerrors.Wrap(readErr, "Read").Error()
		} else {
			info.Reasons["http1"] = "no data read at all"
		}
		return 0, nil
	}

	// parse data
	req, err := http.ReadRequest(bufio.NewReader(buf))
	if err != nil {
		info.Reasons["http1"] = gerrors.Wrap(err, "ReadRequest").Error()
		return 0, nil
	}
	info.Protocols = append(info.Protocols, "http1")
	info.HTTP1Req = req
	if req.Header.Get("Upgrade") == "websocket" {
		info.Protocols = append(info.Protocols, "websocket")
	}

	return 0, nil
}

func hasHTTP2Preface(r io.Reader) bool {
	var b [len(http2.ClientPreface)]byte
	last := 0

	for {
		n, err := gio.Read(r, b[last:], &matchReadTimeout)
		if err != nil {
			return false
		}

		last += n
		eq := string(b[:last]) == http2.ClientPreface[:last]
		if last == len(http2.ClientPreface) {
			return eq
		}
		if !eq {
			return false
		}
	}
}

// 尚未测试
func MatcherHTTP2(r io.Reader, info *Info) (int, error) {
	// try best to read data until 0 or error is returned.
	buf, readErr := gio.ReadHead(r, &matchReadTimeout)

	// process no data error
	if buf.Len() == 0 {
		if readErr != nil {
			info.Reasons["http2"] = gerrors.Wrap(readErr, "Read").Error()
		} else {
			info.Reasons["http2"] = "no data read at all"
		}
		return 0, nil
	}

	if !hasHTTP2Preface(buf) {
		info.Reasons["http2"] = "no HTTP2 preface found"
		return 0, nil
	}

	info.Protocols = append(info.Protocols, "http2")
	framer := http2.NewFramer(ioutil.Discard, buf)
	hdec := hpack.NewDecoder(uint32(4<<10), nil)
	for {
		frame, err := framer.ReadFrame()
		if err != nil {
			if err == io.EOF {
				return 0, nil
			} else {
				return 0, err
			}
		}

		if frame.Header().Type == http2.FrameHeaders {
			hFrame := frame.(*http2.HeadersFrame)
			hFields, err := hdec.DecodeFull(hFrame.HeaderBlockFragment())
			if err != nil {
				return 0, err
			}
			for _, hField := range hFields {
				info.HTTP2Fields = append(info.HTTP2Fields, [2]string{hField.Name, hField.Value})
				if hField.Name == "Upgrade" && hField.Value == "websocket" {
					info.Protocols = append(info.Protocols, "websocket")
				}
				if hField.Name == "content-type" && hField.Value == "application/grpc" {
					info.Protocols = append(info.Protocols, "grpc")
				}
			}
			if hFrame.FrameHeader.Flags&http2.FlagHeadersEndHeaders != 0 {
				return 0, nil
			}
		}

		if frame.Header().Type == http2.FrameSettings {
		}
	}
}

/*
	done := false
	hFrame := frame.(*http2.HeadersFrame)
	hFields, err := hdec.DecodeFull(hFrame.HeaderBlockFragment())
	if err != nil {
		return err
	}
	for _, hField := range hFields {
		info.HTTP2Fields = append(info.HTTP2Fields, [2]string{hField.Name, hField.Value})
		if hField.Name == "Upgrade" && hField.Value == "websocket" {
			info.Protocols = append(info.Protocols, "websocket")
		}
		if hField.Name == "content-type" && hField.Value == "application/grpc" {
			info.Protocols = append(info.Protocols, "grpc")
		}
	}

	switch fType := frame.(type) {
	case *http2.SettingsFrame:
		// Sender acknoweldged the SETTINGS frame. No need to write
		// SETTINGS again.
		if fType.IsAck() {
			break
		}
		if err := framer.WriteSettings(); err != nil {
			return err
		}
	case *http2.HeadersFrame, *http2.ContinuationFrame:
		if _, err := hdec.Write(hFrame.HeaderBlockFragment()); err != nil {
			return err
		}
		done = done || hFrame.FrameHeader.Flags&http2.FlagHeadersEndHeaders != 0
	}
	if done {
		return nil
	}*/
