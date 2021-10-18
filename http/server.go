package http

import (
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/gorilla/mux"
)

type Server struct {
	*serverOpt
	router *mux.Router
}

func NewServer(options ...serverOption) *Server {
	s := &Server{
		serverOpt: newServerOpt(options),
		router:    mux.NewRouter(),
	}
	s.doOpt(s)
	return s
}

func Serve(addr string, options ...serverOption) error {
	return NewServer(options...).Serve(addr)
}

func (s *Server) Serve(addr string) error {
	s.init()
	if s.background {
		go s.serve(addr)
	}
	return s.serve(addr)
}

func (s *Server) serve(addr string) (err error) {
	if s.safety {
		defer func() {
			if v := recover(); v != nil {
				err = v.(error)
			}
		}()
	}
	if s.asProxy {
		return http.ListenAndServe(addr, goproxy.NewProxyHttpServer())
	}
	return http.ListenAndServe(addr, s.router)
}

func (s *Server) init() {
	for _, routes := range s.routes {
		s.router.HandleFunc(routes.path, routes.handleFunc)
	}
}

type serverOption func(*Server)
type serverOptions []serverOption
type route struct {
	path       string
	handleFunc func(http.ResponseWriter, *http.Request)
}

type serverOpt struct {
	serverOptions
	routes     []*route
	asProxy    bool
	background bool
	safety     bool
}

func newServerOpt(options []serverOption) *serverOpt {
	return &serverOpt{
		serverOptions: options,
	}
}

func (opt *serverOpt) doOpt(s *Server) {
	for _, option := range opt.serverOptions {
		option(s)
	}
}

func WithRoute(path string, f func(http.ResponseWriter, *http.Request)) serverOption {
	return func(s *Server) {
		s.WithRoute(path, f)
	}
}

func (s *Server) WithRoute(path string, f func(http.ResponseWriter, *http.Request)) *Server {
	s.routes = append(s.routes, &route{path, f})
	return s
}

func WithAsProxy() serverOption {
	return func(s *Server) {
		s.WithAsProxy()
	}
}

func (s *Server) WithAsProxy() *Server {
	s.asProxy = true
	return s
}

func WithBackground() serverOption {
	return func(s *Server) {
		s.WithBackground()
	}
}

func (s *Server) WithBackground() *Server {
	s.background = true
	return s
}

func WithServerSafety() serverOption {
	return func(s *Server) {
		s.WithServerSafety()
	}
}

func (s *Server) WithServerSafety() *Server {
	s.safety = true
	return s
}
