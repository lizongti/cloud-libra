package http

import (
	"fmt"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/gorilla/mux"
)

type Server struct {
	opts       []serverOpt
	routes     []*route
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

func NewServer(opts ...serverOpt) *Server {
	return &Server{opts: opts}
}

func Serve(addr string, options ...serverOpt) error {
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
	s.router = mux.NewRouter()

	for _, opt := range s.opts {
		opt(s)
	}

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

	if s.proxy {
		s.server = &http.Server{Addr: addr, Handler: goproxy.NewProxyHttpServer()}
	} else {
		s.server = &http.Server{Addr: addr, Handler: s.router}
	}
	return s.server.ListenAndServe()
}

type serverOpt func(*Server)
type serverOption struct{}

var ServerOption serverOption

func (serverOption) WithRoute(path string, f func(http.ResponseWriter, *http.Request)) serverOpt {
	return func(s *Server) {
		s.routes = append(s.routes, &route{path, f})
	}
}

func (s *Server) WithRoute(path string, f func(http.ResponseWriter, *http.Request)) *Server {
	s.opts = append(s.opts, ServerOption.WithRoute(path, f))
	return s
}

func (serverOption) WithProxy() serverOpt {
	return func(s *Server) {
		s.proxy = true
	}
}

func (s *Server) WithProxy() *Server {
	s.opts = append(s.opts, ServerOption.WithProxy())
	return s
}

func (serverOption) WithBackground() serverOpt {
	return func(s *Server) {
		s.background = true
	}
}

func (s *Server) WithBackground() *Server {
	s.opts = append(s.opts, ServerOption.WithBackground())
	return s
}

func (serverOption) WithServerSafety() serverOpt {
	return func(s *Server) {
		s.safety = true
	}
}

func (s *Server) WithServerSafety() *Server {
	s.opts = append(s.opts, ServerOption.WithServerSafety())
	return s
}

func (serverOption) WithErrorFunc(errorFunc func(error)) serverOpt {
	return func(s *Server) {
		s.errorFunc = errorFunc
	}
}

func (s *Server) WithErrorFunc(errorFunc func(error)) *Server {
	s.opts = append(s.opts, ServerOption.WithErrorFunc(errorFunc))
	return s
}
