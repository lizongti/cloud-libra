package http

import (
	"fmt"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/gorilla/mux"
)

type Server struct {
	proxy      bool
	background bool
	safety     bool
	errorFunc  func(error)
	router     *mux.Router
	server     *http.Server
}

type route struct {
	path       string
	handleFunc func(http.ResponseWriter, *http.Request)
}

func NewServer(opts ...funcServerOption) *Server {
	s := &Server{
		router: mux.NewRouter(),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func Serve(addr string, options ...funcServerOption) error {
	return NewServer(options...).Serve(addr)
}

func (s *Server) Serve(addr string) error {
	if s.background {
		go s.serve(addr)
		return nil
	}
	return s.serve(addr)
}

func (s *Server) Close() error {
	return s.server.Shutdown(nil)
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

	if s.proxy {
		s.server = &http.Server{Addr: addr, Handler: goproxy.NewProxyHttpServer()}
	} else {
		s.server = &http.Server{Addr: addr, Handler: s.router}
	}
	return s.server.ListenAndServe()
}

type funcServerOption func(*Server)
type serverOption struct{}

var ServerOption serverOption

func (serverOption) WithRoute(path string, f func(http.ResponseWriter, *http.Request)) funcServerOption {
	return func(s *Server) {
		s.WithRoute(path, f)
	}
}

func (s *Server) WithRoute(path string, f func(http.ResponseWriter, *http.Request)) *Server {
	s.router.HandleFunc(path, f)
	return s
}

func (serverOption) WithProxy() funcServerOption {
	return func(s *Server) {
		s.WithProxy()
	}
}

func (s *Server) WithProxy() *Server {
	s.proxy = true
	return s
}

func (serverOption) WithBackground() funcServerOption {
	return func(s *Server) {
		s.WithBackground()
	}
}

func (s *Server) WithBackground() *Server {
	s.background = true
	return s
}

func (serverOption) WithServerSafety() funcServerOption {
	return func(s *Server) {
		s.WithServerSafety()
	}
}

func (s *Server) WithServerSafety() *Server {
	s.safety = true
	return s
}

func (serverOption) WithErrorFunc(errorFunc func(error)) funcServerOption {
	return func(s *Server) {
		s.WithErrorFunc(errorFunc)
	}
}

func (s *Server) WithErrorFunc(errorFunc func(error)) *Server {
	s.errorFunc = errorFunc
	return s
}
