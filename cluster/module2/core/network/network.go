package network

import (
	"net"
)

type Conn interface {
	Send()
	Close()
}

type Handler interface {
	Handle(net.Conn) Conn
}

type Server interface {
	ListenAndServe(string, Handler) error
	String()
}

type Client interface {
	Dial(string, Handler) error
	String()
}
