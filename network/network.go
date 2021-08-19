package network

import "net"

type Handler func(net.Conn)

type Network interface {
	GetAddr() string
	SetAddr(string)
	GetHandler() Handler
	SetHandler(Handler)
}
