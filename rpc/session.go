package rpc

import (
	"encoding/binary"
	"net"
)

type Session struct {
	conn          net.Conn
	sessionOption *sessionOption
}

// Read implements SessionHandler
func (s *Session) Read() ([]byte, error) {
	header := make([]byte, s.sessionOption.len)
	_, err := s.conn.Read(header)
	if err != nil {
		panic(err)
	}
	dataLens := binary.BigEndian.Uint32(header)
	buffer := make([]byte, dataLens)
	_, err = s.conn.Read(buffer)
	if err != nil {
		panic(err)
	}
	return buffer, nil
}

// Write implements SessionHandler
func (s *Session) Write(data []byte) error {
	buffer := make([]byte, int(s.sessionOption.len)+len(data))
	binary.BigEndian.PutUint32(buffer[:s.sessionOption.len], uint32(len(data)))
	copy(buffer[s.sessionOption.len:], data)
	_, err := s.conn.Write(buffer)
	if err != nil {
		panic(err)
	}
	return nil
}

func NewSession(conn net.Conn, sessionOpts ...SessionOptioner) *Session {
	opt := &sessionOption{}
	for _, v := range sessionOpts {
		v.apply(opt)
	}
	return &Session{
		conn:          conn,
		sessionOption: opt,
	}
}

type sessionOption struct {
	len uint32
}

type SessionOptioner interface {
	apply(*sessionOption)
}

type funcSessionOption func(s *sessionOption)

func (f funcSessionOption) apply(s *sessionOption) {
	f(s)
}

func WithLen(len uint32) funcSessionOption {
	return func(s *sessionOption) {
		s.len = len
	}
}

type SessionHandler interface {
	Write([]byte) error
	Read() ([]byte, error)
}

var _ SessionHandler = (*Session)(nil)
