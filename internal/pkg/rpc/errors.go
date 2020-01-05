package rpc

import (
	"io"
	"net"
	"strings"
)

func IsClosedConn(err error) bool {
	if err == nil {
		return false
	}
	// see https://github.com/golang/net/blob/master/http2/server.go#L676-L679
	if strings.Contains(err.Error(), "use of closed network connection") {
		return true
	}
	return false
}

func IsTimeout(err error) bool {
	e, ok := err.(net.Error)
	return ok && e.Timeout()
}

func IsEOF(err error) bool {
	return err == io.EOF
}
