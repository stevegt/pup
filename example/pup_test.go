package pup

import (
	"io"
	"testing"

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
