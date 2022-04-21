package gsweb

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Context struct {
	request  *http.Request
	response http.ResponseWriter
	ctx      context.Context

	handler ControllerHandler

	hasTimeout bool        //超时标记
	writerMux  *sync.Mutex //写保护
}

func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		request:   r,
		response:  w,
		ctx:       r.Context(),
		writerMux: &sync.Mutex{},
	}
}

func (ctx *Context) BaseContext() context.Context {
	return ctx.request.Context()
}

/**---------------实现context接口-------------------**/
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

/**-----------------------------------------------**/

func (ctx *Context) WriterMux() *sync.Mutex {
	return ctx.WriterMux()
}

func (ctx *Context) GetRequest() *http.Request {
	return ctx.request
}

func (ctx *Context) GetResponse() http.ResponseWriter {
	return ctx.response
}

func (ctx *Context) SetHasTimeout() {
	ctx.hasTimeout = true
}

func (ctx *Context) HasTimeout() bool {
	return ctx.hasTimeout
}

/**-------------------------------------------------------**/

func (ctx *Context) FormAll() map[string][]string {
	if ctx.request != nil {
		return map[string][]string(ctx.request.PostForm)
	}
	return map[string][]string{}
}

func (ctx *Context) FormInt(key string, def int) int {
	params := ctx.FormAll()
	return handInt(key, def, params)
}

func (ctx *Context) FormString(key string, def string) string {
	params := ctx.FormAll()
	return handString(key, def, params)
}

func (ctx *Context) FormArray(key string, def []string) []string {
	params := ctx.FormAll()
	return handArray(key, def, params)
}

func (ctx *Context) QueryAll() map[string][]string {
	if ctx.request != nil {
		return map[string][]string(ctx.request.URL.Query())
	}
	return map[string][]string{}
}

func (ctx *Context) QueryInt(key string, def int) int {
	params := ctx.QueryAll()
	return handInt(key, def, params)
}

func (ctx *Context) QueryString(key string, def string) string {
	params := ctx.QueryAll()
	return handString(key, def, params)
}

func (ctx *Context) QueryArray(key string, def []string) []string {
	params := ctx.QueryAll()
	return handArray(key, def, params)
}

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
	} else {
		return errors.New("ctx.request is empty")
	}
	return nil
}

func (ctx *Context) Json(status int, obj interface{}) error {
	if ctx.HasTimeout() {
		return nil
	}
	ctx.response.Header().Set("Content-Type", "application/json")
	ctx.response.WriteHeader(status)
	byteObj, err := json.Marshal(obj)
	if err != nil {
		ctx.response.WriteHeader(500)
		return err
	}
	ctx.response.Write(byteObj)
	return nil
}

func (ctx *Context) HTML(status int, obj interface{}, template string) error {
	return nil
}

func (ctx *Context) Text(status int, obj string) error {
	return nil
}

func handInt(key string, def int, params map[string][]string) int {
	if val, ok := params[key]; ok {
		l := len(val)
		if l > 0 {
			res, err := strconv.Atoi(val[l-1])
			if err != nil {
				return def
			}
			return res
		}
	}
	return def
}

//---------------------------------------------------------------------------
func handString(key string, def string, params map[string][]string) string {
	if val, ok := params[key]; ok {
		l := len(val)
		if l > 0 {
			return val[l-1]
		}
	}
	return def
}

func handArray(key string, def []string, params map[string][]string) []string {
	if val, ok := params[key]; ok {
		return val
	}
	return def
}
