package pup

import (
	"bufio"
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

func handleStream(stream io.ReadWriteCloser) (err error) {
	defer Return(&err)
	// read the leading hash
	r := bufio.NewReader(stream)
	hash, err := r.ReadString('\n')
	Ck(err)
	// lookup the hash in the builtins

	// lookup the hash in peer registrations
	// pipe the rest of the stream to the peer
	err = builtinEcho(hash, stream)
	return
}

func builtinEcho(hash string, stream io.ReadWriteCloser) (err error) {
	defer Return(&err)
	/*
		scanner := bufio.NewScanner(stream)

		for scanner.Scan() {
			_, err = stream.Write(Spf("%s\n", scanner.Text()))
			Ck(err)
		}

		err = scanner.Err()
		Ck(err)
	*/

	_, err = io.Copy(stream, stream)
	Ck(err)
	return
}
