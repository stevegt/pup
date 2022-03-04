package pup

import (
	"io"
	"net"
	"testing"
	"time"

	. "github.com/stevegt/goadapt"
)

type MockReadWriteCloser struct {
	readbuf  []byte
	writebuf []byte
	readpos  int
}

func (c *MockReadWriteCloser) Read(out []byte) (n int, err error) {
	if c.readpos >= len(c.readbuf) {
		return 0, io.EOF
	}
	// set end to the end of the readbuf or out, depending on which comes
	// first
	end := c.readpos + len(out)
	if end > len(c.readbuf) {
		end = len(c.readbuf)
	}
	n = copy(out, c.readbuf[c.readpos:end])
	c.readpos += n
	return n, nil
}

func (c *MockReadWriteCloser) Write(data []byte) (n int, err error) {
	c.writebuf = append(c.writebuf, data...)
	return len(data), nil
}

func (c *MockReadWriteCloser) Close() error {
	return nil
}

func echoContent(hash []byte, stream io.ReadWriteCloser) (err error) {
	defer Return(&err)
	_, err = io.Copy(stream, stream)
	Ck(err)
	return
}

func echoHash(hash []byte, stream io.ReadWriteCloser) (err error) {
	defer Return(&err)
	_, err = stream.Write(hash)
	Ck(err)
	return
}

var s1hash = "somehash"
var s1content = "first line\nsecond line\n"
var s1 = Spf("%s\n%s", s1hash, s1content)

var s2hash = "anotherhash"
var s2content = "1 first line\n2 second line\n"
var s2 = Spf("%s\n%s", s2hash, s2content)

func TestStream(t *testing.T) {

	s := &Server{}

	s.Register("somehash", echoContent)
	s.Register("anotherhash", echoHash)

	rwc := &MockReadWriteCloser{readbuf: []byte(s1)}
	err := s.handleStream(rwc)
	Tassert(t, err == nil, "handleStream %v", err)
	Tassert(t, string(rwc.writebuf) == s1content, "writebuf '%v'", string(rwc.writebuf))

	rwc = &MockReadWriteCloser{readbuf: []byte(s2)}
	err = s.handleStream(rwc)
	Tassert(t, err == nil, "handleStream %v", err)
	Tassert(t, string(rwc.writebuf) == s2hash, "writebuf '%v'", string(rwc.writebuf))

}

func TestServer(t *testing.T) {
	port := 10842
	go func() {
		s := &Server{}
		s.Register("somehash", echoContent)
		err := s.Serve("127.0.0.1", port)
		Tassert(t, err == nil, "Serve: %v", err)
	}()
	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", Spf(":%d", port))
	Tassert(t, err == nil, "Dial: %v", err)
	defer conn.Close()

	go func() {
		_, err := conn.Write([]byte(s1))
		Tassert(t, err == nil, "conn.Write: %v", err)
	}()
	time.Sleep(1 * time.Second)

	got := make([]byte, 1024)
	n, err := conn.Read(got)
	Tassert(t, string(got[:n]) == s1content, "wanted '%v' got '%v'", []byte(s1content), got)

}
