package context

type Context interface {
	GetAddr() string
	GetLeftAddr() string
}
