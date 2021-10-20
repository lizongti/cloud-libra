package http

import (
	"fmt"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/gorilla/mux"
)

type Server struct {
	*serverOpt
	router *mux.Router
	server *http.Server
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
		return nil
	}
	return s.serve(addr)
}

func (s *Server) Close() error {
	return s.server.Shutdown(nil)
}

func (s *Server) init() {
	for _, routes := range s.routes {
		s.router.HandleFunc(routes.path, routes.handleFunc)
	}
}

func (s *Server) serve(addr string) (err error) {
	if s.errorFunc != nil {
		defer func() {
			s.errorFunc(err)
			err = nil
		}()
	}

	if s.safety {
		defer func() {
			if e := recover(); e != nil {
				err = fmt.Errorf("%v", e)
			}
		}()
	}

	if s.asProxy {
		s.server = &http.Server{Addr: addr, Handler: goproxy.NewProxyHttpServer()}
	} else {
		s.server = &http.Server{Addr: addr, Handler: s.router}
	}
	return s.server.ListenAndServe()
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
	errorFunc  func(error)
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

func WithErrorFunc(errorFunc func(error)) serverOption {
	return func(s *Server) {
		s.WithErrorFunc(errorFunc)
	}
}

func (s *Server) WithErrorFunc(errorFunc func(error)) *Server {
	s.errorFunc = errorFunc
	return s
}
