package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Context struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	ctx            context.Context
	handler        ControllerHandler

	// 是否超时标记位
	hasTimeout bool
	// 写保护机制
	writeMux *sync.Mutex
}

func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		request:        r,
		responseWriter: w,
		ctx:            r.Context(),
		writeMux:       &sync.Mutex{},
	}
}

// Context #region base function

func (ctx *Context) WriteMux() *sync.Mutex {
	return ctx.writeMux
}
func (ctx *Context) GetRequest() *http.Request {
	return ctx.request
}

func (ctx *Context) GetResponse() http.ResponseWriter {
	return ctx.responseWriter
}

func (ctx *Context) SetHasTimeout() {
	ctx.hasTimeout = true
}

func (ctx *Context) HasTimeout() bool {
	return ctx.hasTimeout
}

// #endregion

// #region context function

func (ctx *Context) BaseContext() context.Context {
	return ctx.request.Context()
}

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.BaseContext().Deadline()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.BaseContext().Done()
}

func (ctx *Context) Err() error {
	return ctx.BaseContext().Err()
}

func (ctx *Context) Value(key interface{}) interface{} {
	return ctx.BaseContext().Value(key)
}

func (ctx *Context) QueryAll() map[string][]string {
	if ctx.request != nil {
		return ctx.request.URL.Query()
	}
	return map[string][]string{}
}

func (ctx *Context) QueryInt(key string, def int) int {
	params := ctx.QueryAll()
	if values, ok := params[key]; ok {
		len := len(values)
		if len > 0 {
			intval, err := strconv.Atoi(values[len-1])
			if err != nil {
				return def
			}
			return intval
		}
	}
	return def
}

func (ctx *Context) QueryString(key string, def string) string {
	params := ctx.QueryAll()
	if values, ok := params[key]; ok {
		len := len(values)
		if len > 0 {
			return values[len-1]
		}
	}
	return def
}

func (ctx *Context) QueryArray(key string, def []string) []string {
	params := ctx.QueryAll()
	if values, ok := params[key]; ok {
		return values
	}
	return def
}

func (ctx *Context) FormAll() map[string][]string {
	if ctx.request != nil {
		return ctx.request.PostForm
	}
	return map[string][]string{}
}

func (ctx *Context) FormInt(key string, def int) int {
	params := ctx.FormAll()
	if value, ok := params[key]; ok {
		len := len(value)
		if len > 0 {
			intval, err := strconv.Atoi(value[len-1])
			if err != nil {
				return def
			}
			return intval
		}
	}
	return def
}

func (ctx *Context) FormString(key string, def string) string {
	params := ctx.FormAll()
	if values, ok := params[key]; ok {
		len := len(values)
		if len > 0 {
			return values[len-1]
		}
	}
	return def
}

func (ctx *Context) FormArray(key string, def []string) []string {
	params := ctx.FormAll()
	if values, ok := params[key]; ok {
		return values
	}
	return def
}

// #endregion

// BindJson #region application/json post
func (ctx *Context) BindJson(obj interface{}) error {
	if ctx.request != nil {
		body, err := ioutil.ReadAll(ctx.request.Body)
		if err != nil {
			return err
		}
		ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		err = json.Unmarshal(body, obj)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ctx *Context) Json(status int, obj interface{}) error {
	if ctx.HasTimeout() {
		return nil
	}
	ctx.responseWriter.Header().Set("Content-Type", "application/json")
	ctx.responseWriter.WriteHeader(status)
	byt, err := json.Marshal(obj)
	if err != nil {
		ctx.responseWriter.WriteHeader(500)
		return err
	}
	ctx.responseWriter.Write(byt)
	return nil
}

func (ctx *Context) Text(status int, obj interface{}) error {
	return nil
}

func (ctx *Context) HTML(status int, obj interface{}) error {
	return nil
}

// #endregion
