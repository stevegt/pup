package main

import (
	"io"

	. "github.com/stevegt/goadapt"

	"github.com/stevegt/pup"
)

const REGISTER = "sha256:c17dcddbc7b307ab652109d2c1a01fdd53890dffcbce3215da41d8104e551b0b"

func Dispatcher(host string, port int) (err error) {
	defer Return(&err)
	s := pup.Server{}
	_ = s
	s.Register(REGISTER, registrar)
	err = s.Serve(host, port)
	Ck(err)
	return
}

/*
// proxy represents a remote peer and routes traffic to it
func proxy(caller io.ReadWriteCloser, lambdaConn net.Conn) {
	// XXX
}
*/

func registrar(hash []byte, peer io.ReadWriteCloser) (err error) {
	// XXX verify hash

	// XXX make an f() that routes all messages with hash to peer
	f := func(peer io.ReadWriteCloser) (err error) {
		// err = proxy(caller, lambdaConn)
		return
	}
	_ = f
	// s.Register(hash, f)
	return
}
