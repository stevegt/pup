package pup

import (
	"errors"
	"io"
	"net"

	. "github.com/stevegt/goadapt"
)

func Serve(host string, port int) {
	l, err := net.Listen("tcp", Spf("%s:%d", host, port))
	Ck(err)
	defer l.Close()
	Pl("Listening on", Spf("%s:%d", host, port))

	for {
		conn, err := l.Accept()
		if err != nil {
			Pl("error accepting:", err.Error())
		}
		go handleTcp(conn)
	}
}

func handleTcp(conn net.Conn) {
	// XXX deal with whitelist/blacklist here
	err := handleStream(conn)
	if err != nil {
		Pl("error handling stream:", err.Error())
		return
	}
}

var ELONGLINE = errors.New("no newline found -- would overflow readLine output buffer")

func readLine(stream io.Reader, max int) (line []byte, err error) {
	c := make([]byte, 1)
	buf := make([]byte, max)
	var n int
	for n = 0; n < len(buf); n++ {
		_, err = stream.Read(c)
		if err != nil {
			return buf[:n], err
		}
		if c[0] == '\n' {
			return buf[:n], err
		}
		buf[n] = c[0]
	}
	return buf[:n], ELONGLINE
}

type Lambda func([]byte, io.ReadWriteCloser) error

type Registry map[string]Lambda

func (r *Registry) Get(hash string) (lambda Lambda) {
	lambda, _ = (*r)[hash]
	return
}

var registry = Registry{
	"somehash":    echoContent,
	"anotherhash": echoHash,
}

func handleStream(stream io.ReadWriteCloser) (err error) {
	defer Return(&err)

	// read the leading hash
	hash, err := readLine(stream, 1024)
	Ck(err)

	// get lambda by looking up the hash in the registry
	lambda := registry.Get(string(hash))

	// pipe the rest of the stream to the lambda
	err = lambda(hash, stream)
	return
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
