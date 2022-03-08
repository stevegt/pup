package main

import (
	"io"
	"strings"

	. "github.com/stevegt/goadapt"

	"github.com/stevegt/pup"
)

const REGISTER = "sha256:c17dcddbc7b307ab652109d2c1a01fdd53890dffcbce3215da41d8104e551b0b"

type Dispatcher struct {
	server pup.Server
}

func (d *Dispatcher) Dispatch(host string, port int) (err error) {
	defer Return(&err)
	d.server = pup.Server{}
	s := d.server
	s.Register(REGISTER, d.registrar)
	err = s.Serve(host, port)
	Ck(err)
	return
}

func (d *Dispatcher) registrar(hash []byte, peer io.ReadWriteCloser) (err error) {
	defer Return(&err)

	// XXX verify hash == REGISTER

	// parse first line of content
	line, err := pup.Readline(peer, 1024)
	Ck(err)
	parts := strings.Split(string(line), " ")
	// XXX ensure len(parts) is 2
	cmd := parts[0]
	subhash := parts[1]
	switch cmd {
	case "a":
		// make a lambda that routes all streams with hash to peer
		f := func(subhash []byte, caller io.ReadWriteCloser) (err error) {
			proxy(caller, peer)
			return
		}
		d.server.Register(string(subhash), f)
	default:
		Pf("unknown registrar cmd: %s\n", cmd)
	}

	return
}

// proxy routes traffic between caller and peer
func proxy(caller, peer io.ReadWriteCloser) {
	// XXX might need to do something more useful with io.Copy errors
	go func() {
		_, err := io.Copy(caller, peer)
		if err != nil {
			Spf("peer to caller: %v", err)
		}
	}()
	go func() {
		_, err := io.Copy(peer, caller)
		if err != nil {
			Spf("caller to peer: %v", err)
		}
	}()
}
