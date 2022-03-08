package pup

import (
	"errors"
	"io"
	"net"
	"syscall"

	. "github.com/stevegt/goadapt"
)

type Registration struct {
	Hash   string
	Lambda Lambda
}

type Server struct {
	registry *registry
}

func (s *Server) Register(hash string, lambda Lambda) {
	if s.registry == nil {
		s.registry = &registry{}
	}
	s.registry.put(hash, lambda)
}

func (s *Server) Dereference(hash string) (lambda Lambda) {
	if s.registry == nil {
		s.registry = new(registry)
	}
	return s.registry.get(hash)
}

func (s *Server) Registrations() (res []Registration) {
	if s.registry == nil {
		s.registry = new(registry)
	}
	return s.registry.ls()
}

func (s *Server) Serve(host string, port int) (err error) {
	defer Return(&err)
	l, err := net.Listen("tcp", Spf("%s:%d", host, port))
	Ck(err)
	defer l.Close()
	Pl("Listening on", Spf("%s:%d", host, port))

	for {
		conn, err := l.Accept()
		if err != nil {
			Pl("error accepting:", err.Error())
		}
		go s.handleTcp(conn)
	}
}

func (s *Server) handleTcp(conn net.Conn) {
	// XXX deal with whitelist/blacklist here
	err := s.handleStream(conn)
	if err != nil {
		Pl("error handling stream:", err.Error())
		return
	}
}

type Error struct {
	Errno syscall.Errno
	Msg   string
}

func (e Error) Error() string {
	return Spf("%s: %s", e.Errno.Error(), e.Msg)
}

func (s *Server) handleStream(stream io.ReadWriteCloser) (err error) {
	defer Return(&err)

	// read the leading hash
	hash, err := Readline(stream, 1024)
	Ck(err)

	// get lambda by looking up the hash in the registry
	lambda := s.Dereference(string(hash))

	if lambda == nil {
		return Error{syscall.ENOSYS, string(hash)}
	}

	// pipe the rest of the stream to the lambda
	err = lambda(hash, stream)
	return
}

var ELONGLINE = errors.New("no newline found -- would overflow Readline output buffer")

func Readline(stream io.Reader, max int) (line []byte, err error) {
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

type registry map[string]Lambda

func (r *registry) put(hash string, lambda Lambda) {
	(*r)[hash] = lambda
	return
}

func (r *registry) get(hash string) (lambda Lambda) {
	lambda, _ = (*r)[hash]
	return
}

func (r *registry) ls() (res []Registration) {
	for k, v := range *r {
		res = append(res, Registration{k, v})
	}
	return
}
