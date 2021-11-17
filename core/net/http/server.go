package http

import (
	"fmt"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/gorilla/mux"
)

type Server struct {
	opts   serverOptions
	router *mux.Router
	server *http.Server
}

type Route struct {
	Path string
	Func func(http.ResponseWriter, *http.Request)
}

func NewServer(opt ...ApplyServerOption) *Server {
	opts := defaultServerOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	s := &Server{
		opts:   opts,
		router: mux.NewRouter(),
	}

	return s
}

func Serve(addr string, opt ...ApplyServerOption) error {
	return NewServer(opt...).Serve(addr)
}

func (s *Server) Serve(addr string) error {
	for _, route := range s.opts.routes {
		s.router.HandleFunc(route.Path, route.Func)
	}

	if s.opts.background {
		go s.serve(addr)
		return nil
	}
	return s.serve(addr)
}

func (s *Server) Close() error {
	return s.server.Shutdown(nil)
}

func (s *Server) serve(addr string) (err error) {
	if s.opts.safety {
		defer func() {
			if e := recover(); e != nil {
				err = fmt.Errorf("%v", e)
			}
		}()
	}

	if s.opts.proxy {
		s.server = &http.Server{Addr: addr, Handler: goproxy.NewProxyHttpServer()}
	} else {
		s.server = &http.Server{Addr: addr, Handler: s.router}
	}
	return s.server.ListenAndServe()
}

type serverOptions struct {
	proxy      bool
	background bool
	safety     bool
	routes     []Route
}

var defaultServerOptions = serverOptions{
	proxy:      false,
	background: false,
	safety:     false,
	routes:     nil,
}

type ApplyServerOption interface {
	apply(*serverOptions)
}

type funcServerOption func(*serverOptions)

func (fso funcServerOption) apply(so *serverOptions) {
	fso(so)
}

type serverOption int

var ServerOption serverOption

func (serverOption) Route(routes ...Route) funcServerOption {
	return func(so *serverOptions) {
		so.routes = append(so.routes, routes...)
	}
}

func (s *Server) WithRoute(routes ...Route) *Server {
	ServerOption.Route(routes...).apply(&s.opts)
	return s
}

func (serverOption) Proxy() funcServerOption {
	return func(so *serverOptions) {
		so.proxy = true
	}
}

func (s *Server) WithProxy() *Server {
	ServerOption.Proxy().apply(&s.opts)
	return s
}

func (serverOption) Background() funcServerOption {
	return func(so *serverOptions) {
		so.background = true
	}
}

func (s *Server) WithBackground() *Server {
	ServerOption.Background().apply(&s.opts)
	return s
}

func (serverOption) Safety() funcServerOption {
	return func(so *serverOptions) {
		so.safety = true
	}
}

func (s *Server) WithSafety() *Server {
	ServerOption.Safety().apply(&s.opts)
	return s
}
