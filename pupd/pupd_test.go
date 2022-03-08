package main

import (
	"io"
	"net"
	"testing"
	"time"

	. "github.com/stevegt/goadapt"
	"github.com/stevegt/pup"
)

const port = 10843

const CALLBACK = "sha256:cdcae2a18f7bc3980dbe3a5e173cd6edb9258bd496a6eb067fe8103fd43d2a05"

var s1hash = REGISTER
var s1content = Spf("a %s\n", CALLBACK)
var s1 = Spf("%s\n%s", s1hash, s1content)

var s2hash = CALLBACK
var s2content = "testing callback\n"
var s2 = Spf("%s\n%s", s2hash, s2content)

func peer(t *testing.T) (err error) {
	// connect to pupd
	conn, err := net.Dial("tcp", Spf(":%d", port))
	Tassert(t, err == nil, "Dial: %v", err)
	defer conn.Close()

	// register us as a lambda
	go func() {
		_, err := conn.Write([]byte(s1))
		Tassert(t, err == nil, "conn.Write: %v", err)
	}()
	time.Sleep(1 * time.Second)
	// expect empty output
	got := make([]byte, 1024)
	n, err := conn.Read(got)
	Tassert(t, n == 0, "wanted '' got '%v'", got)

	// verify hash
	hash, err := pup.Readline(conn, 1024)
	Tassert(t, err == nil, "%v", err)
	Tassert(t, string(hash) == CALLBACK, "wanted '%s' got '%v'", CALLBACK, hash)
	// echo back the content
	_, err = io.Copy(conn, conn)
	Tassert(t, err == nil, "%v", err)
	return
}

func TestDispatcher(t *testing.T) {

	// start dispatcher
	go func() {
		err := Dispatcher("127.0.0.1", port)
		Tassert(t, err == nil, "Dispatcher: %v", err)
	}()
	time.Sleep(1 * time.Second)

	// start peer -- this peer will register a lambda that in turn
	// just echoes back the content of any message sent to it
	go peer(t)
	time.Sleep(1 * time.Second)

	// connect to dispatcher
	conn, err := net.Dial("tcp", Spf(":%d", port))
	Tassert(t, err == nil, "Dial: %v", err)
	defer conn.Close()

	// send a message to the peer
	go func() {
		_, err := conn.Write([]byte(s2))
		Tassert(t, err == nil, "conn.Write: %v", err)
	}()
	time.Sleep(1 * time.Second)

	// verify the response content matches what we sent
	got := make([]byte, 1024)
	n, err := conn.Read(got)
	Tassert(t, string(got[:n]) == s2content, "wanted '%v' got '%v'", []byte(s2content), got)

}
