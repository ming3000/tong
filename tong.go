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
	"time"
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

func prependMiddleware(h HandlerFunc, middleware ...MiddlewareFunc) HandlerFunc {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}

// $--- tong struct define ---
type Tong struct {
	Server             *http.Server
	Listener           net.Listener
	router             *Router
	sysMiddleware      []MiddlewareFunc
	customerMiddleware []MiddlewareFunc
	cronList           []*common.Cron
	pool               sync.Pool
	Debug              bool
	Logger             *common.Logger
	NotFoundHandler    HandlerFunc
	HTTPErrorHandler   ErrorHandlerFunc
}

// New creates an instance of Wu
func New() *Tong {
	tong := &Tong{Server: new(http.Server)}
	tong.Server.Handler = tong
	tong.router = NewRouter()
	tong.sysMiddleware = make([]MiddlewareFunc, 0)
	tong.customerMiddleware = make([]MiddlewareFunc, 0)
	tong.cronList = make([]*common.Cron, 0)
	tong.pool.New = func() interface{} {
		return tong.NewContext(nil, nil)
	}
	tong.Debug = true
	tong.Logger = common.NewDefaultLogger(tong.Debug)
	tong.NotFoundHandler = NotFoundHandler
	tong.HTTPErrorHandler = DefaultHTTPErrorHandler
	return tong
}

// ServeHTTP implements the http.Handler interface,
// it is used to serve each HTTP requests.
func (t *Tong) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// acquire context instance
	c := t.pool.Get().(*Context)
	c.Reset(r, w, c.logger, common.NewDefaultLRUCache())

	h := NotFoundHandler
	if t.sysMiddleware == nil {
		t.router.Find(r.Method, parsePath(r), c)
		h = c.Handler()
		h = prependMiddleware(h, t.customerMiddleware...)
	} else {
		h = func(c *Context) error {
			t.router.Find(r.Method, parsePath(r), c)
			h := c.Handler()
			h = prependMiddleware(h, t.customerMiddleware...)
			return h(c)
		}
		h = prependMiddleware(h, t.sysMiddleware...)
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
	// stop all cron jobs
	t.stopCronJobs()
	return t.Server.Close()
}

// Shutdown stops the server gracefully.
// It internally calls `http.Server#Shutdown()`.
func (t *Tong) Shutdown(ctx context.Context) error {
	// stop all cron jobs
	t.stopCronJobs()
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
	// start all cron jobs
	t.startCronJobs()
	return s.Serve(t.Listener)
}

func (t *Tong) NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		request:      r,
		response:     NewResponse(w),
		handler:      NotFoundHandler,
		logger:       t.Logger,
		requestCache: common.NewDefaultLRUCache(),
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
		h := prependMiddleware(handler, middleware...)
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

func (t *Tong) AddCronJob(initialPeriod, stepPeriod, maxPeriod time.Duration, job common.Job) {
	c := common.NewCron(initialPeriod, stepPeriod, maxPeriod)
	c.Do(job)
	t.cronList = append(t.cronList, c)
}

func (t *Tong) startCronJobs() {
	for i := range t.cronList {
		t.cronList[i].Start()
	}
}

func (t *Tong) stopCronJobs() {
	for i := range t.cronList {
		t.cronList[i].Stop()
	}
}
