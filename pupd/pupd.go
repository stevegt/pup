package main

import (
	"fmt"
	"io"
	"net"

	. "github.com/stevegt/goadapt"
)

func main() {
	fmt.Println("vim-go")
	s.Register("sha256:c17dcddbc7b307ab652109d2c1a01fdd53890dffcbce3215da41d8104e551b0b", s.registerRemote)

}

// proxy represents a remote peer and routes traffic to it
func proxy(caller io.ReadWriteCloser, lambdaConn net.Conn) {
	// XXX
}

func (s *Server) registerRemote(hash []byte, lambdaConn net.Conn) (err error) {
	// XXX make an f() that routes all messages with hash to conn
	f := func(caller io.ReadWriteCloser) (err error) {
		err = proxy(caller, lambdaConn)
		return
	}
	s.Register(hash, f)
	return
}
