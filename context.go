package tong

import (
	"encoding/json"
	"errors"
	"github.com/ming3000/tong/common"
	"net/http"
	"strconv"
)

// Context is context for every goroutine
type Context struct {
	request      *http.Request
	response     *Response
	path         string
	handler      HandlerFunc
	requestCache common.RAMCache
	logger       common.Logger
}

// $--- utils ---
func (c *Context) Reset(r *http.Request, w http.ResponseWriter, cache common.RAMCache) {
	c.request = r
	c.response.Reset(w)
	c.path = ""
	c.handler = NotFoundHandler
	c.requestCache = cache
}

func (c *Context) Redirect(code int, url string) error {
	if code < http.StatusMultipleChoices || code > http.StatusPermanentRedirect {
		return errors.New("redirect code error")
	} // if>
	c.response.Header().Set(common.HeaderLocation, url)
	c.response.WriteHeader(code)
	return nil
}

// $--- Getter ---
func (c *Context) Request() *http.Request {
	return c.request
}

func (c *Context) Response() *Response {
	return c.response
}

func (c *Context) Path() string {
	return c.path
}

func (c *Context) Handler() HandlerFunc {
	return c.handler
}

func (c *Context) RequestCache() common.RAMCache {
	return c.requestCache
}

func (c *Context) Logger() common.Logger {
	return c.logger
}

// $--- Writer ---
func (c *Context) WriteContentType(value string) {
	head := c.response.Header()
	if head.Get(common.HeaderContentType) == "" {
		head.Set(common.HeaderContentType, value)
	}
}

func (c *Context) Blob(code int, contentType string, data []byte) error {
	c.response.WriteHeader(code)
	c.WriteContentType(contentType)
	_, err := c.response.Write(data)
	return err
}

func (c *Context) Json(code int, value interface{}, indent string) error {
	enc := json.NewEncoder(c.response)
	if indent != "" {
		enc.SetIndent("", indent)
	}
	c.WriteContentType(common.MIMEApplicationJSONCharsetUTF8)
	c.response.Status = code
	return enc.Encode(value)
}

func (c *Context) String(code int, value string) error {
	return c.Blob(code, common.MIMETextPlainCharsetUTF8, []byte(value))
}

// $--- Query Reader ---
func (c *Context) QueryInt(key string, defaultValue int) int {
	value := c.request.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	} // if>

	ret, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	} // if>
	return ret
}

func (c *Context) QueryFloat(key string, defaultValue float64) float64 {
	value := c.request.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	} // if>

	ret, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	} // if>
	return ret
}

func (c *Context) QueryString(key string, defaultValue string) string {
	value := c.request.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	} // if>
	return value
}

// $--- Post Reader ---
func (c *Context) PostInt(key string, defaultValue int) int {
	value := c.request.PostFormValue(key)
	if value == "" {
		return defaultValue
	} // if>

	ret, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	} // if>
	return ret
}

func (c *Context) PostFloat(key string, defaultValue float64) float64 {
	value := c.request.PostFormValue(key)
	if value == "" {
		return defaultValue
	} // if>

	ret, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	} // if>
	return ret
}

func (c *Context) PostString(key string, defaultValue string) string {
	value := c.request.PostFormValue(key)
	if value == "" {
		return defaultValue
	} // if>

	return value
}
