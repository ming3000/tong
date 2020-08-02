package tong

import (
	"context"
	"errors"
	"github.com/ming3000/tong/common"
	"net"
	"net/http"
	"reflect"
	"runtime"
	"sync"
)

// $--- handler type define ---
// HandlerFunc defines a function to serve HTTP requests
type HandlerFunc func(c *Context) error

// ErrorHandlerFunc is a centralized HTTP error handler.
type ErrorHandlerFunc func(*Context, error)

// MiddlewareFunc defines a function to process middleware
type MiddlewareFunc func(HandlerFunc) HandlerFunc

// $--- default handler ---
var NotFoundHandler = func(c *Context) error {
	return errors.New("handler not found")
}

// DefaultHTTPErrorHandler the default HTTP error handler.
// it sends a string response with status code StatusInternalServerError.
var DefaultHTTPErrorHandler = func(c *Context, err error) {
	_ = c.String(http.StatusInternalServerError, err.Error())
}

// $--- utils func ---
// reflect name of HandlerFunc
func handlerName(h HandlerFunc) string {
	t := reflect.ValueOf(h).Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	}
	return t.String()
}

func applyMiddleware(h HandlerFunc, middleware ...MiddlewareFunc) HandlerFunc {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}

// $--- tong struct define ---
type Tong struct {
	router             *Router
	sysMiddleware      []MiddlewareFunc
	customerMiddleware []MiddlewareFunc
	pool               sync.Pool
	Server             *http.Server
	Listener           net.Listener
	Debug              bool
	NotFoundHandler    HandlerFunc
	HTTPErrorHandler   ErrorHandlerFunc
}

// New creates an instance of Wu
func New() *Tong {
	tong := &Tong{Server: new(http.Server)}
	tong.Server.Handler = tong
	tong.router = NewRouter()
	tong.pool.New = func() interface{} {
		return tong.NewContext(nil, nil)
	}
	tong.HTTPErrorHandler = DefaultHTTPErrorHandler
	return tong
}

// ServeHTTP implements the http.Handler interface,
// it is used to serve each HTTP requests.
func (t *Tong) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// acquire context instance
	c := t.pool.Get().(*Context)
	c.Reset(r, w, common.NewDefaultRAMCache())

	h := NotFoundHandler
	if t.sysMiddleware == nil {
		t.router.Find(r.Method, parsePath(r), c)
		h = c.Handler()
		h = applyMiddleware(h, t.customerMiddleware...)
	} else {
		h = func(c *Context) error {
			t.router.Find(r.Method, parsePath(r), c)
			h := c.Handler()
			h = applyMiddleware(h, t.customerMiddleware...)
			return h(c)
		}
		h = applyMiddleware(h, t.sysMiddleware...)
	}

	// handle error
	if err := h(c); err != nil {
		t.HTTPErrorHandler(c, err)
	}

	// Release context
	t.pool.Put(c)
}

// Start starts an HTTP server.
func (t *Tong) Start(address string) error {
	t.Server.Addr = address
	return t.StartServer(t.Server)
}

// Close immediately stops the server.
func (t *Tong) Close() error {
	return t.Server.Close()
}

// Shutdown stops the server gracefully.
// It internally calls `http.Server#Shutdown()`.
func (t *Tong) Shutdown(ctx context.Context) error {
	return t.Server.Shutdown(ctx)
}

// StartServer starts a custom http server.
func (t *Tong) StartServer(s *http.Server) error {
	var err error
	s.Handler = t
	t.Listener, err = net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	return s.Serve(t.Listener)
}

func (t *Tong) NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		request:      r,
		response:     NewResponse(w),
		handler:      NotFoundHandler,
		requestCache: common.NewDefaultRAMCache(),
	}
}

func (t *Tong) AddSysMiddleware(middleware ...MiddlewareFunc) {
	t.sysMiddleware = append(t.sysMiddleware, middleware...)
}

func (t *Tong) AddCustomerMiddleware(middleware ...MiddlewareFunc) {
	t.customerMiddleware = append(t.customerMiddleware, middleware...)
}

func (t *Tong) GET(p string, h HandlerFunc, m ...MiddlewareFunc) *RouteInfo {
	return t.Add(http.MethodGet, p, h, m...)
}

func (t *Tong) POST(p string, h HandlerFunc, m ...MiddlewareFunc) *RouteInfo {
	return t.Add(http.MethodPost, p, h, m...)
}

func (t *Tong) Add(method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) *RouteInfo {
	name := handlerName(handler)
	t.router.Add(method, path, func(c *Context) error {
		h := applyMiddleware(handler, middleware...)
		return h(c)
	})
	r := &RouteInfo{
		Method: method,
		Path:   path,
		Name:   name,
	}
	t.router.routes[method+path] = r
	return r
}
