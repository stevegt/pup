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

func readLine(buf []byte, stream io.Reader) (n int, err error) {
	c := make([]byte, 1)
	for n = 0; n < len(buf); n++ {
		_, err = stream.Read(c)
		if err != nil {
			return
		}
		if c[0] == '\n' {
			return
		}
		buf[n] = c[0]
	}
	return n, ELONGLINE
}

func handleStream(stream io.ReadWriteCloser) (err error) {
	defer Return(&err)
	// read the leading hash
	hash := make([]byte, 1024)
	_, err = readLine(hash, stream)
	Ck(err)
	// lookup the hash in the builtins

	// lookup the hash in peer registrations
	// pipe the rest of the stream to the peer
	err = echoContent(hash, stream)
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
	_, err = stream.Write(append(hash, byte('\n')))
	Ck(err)
	return
}
