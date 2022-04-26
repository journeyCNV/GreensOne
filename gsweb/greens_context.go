package gsweb

import (
	"GreensOne/cast"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Context struct {
	request  *http.Request
	response http.ResponseWriter
	ctx      context.Context

	handlers []ControllerHandler //当前请求的控制器链
	index    int                 //当前请求控制器链中下标
	params   map[string]string   //路由通配符匹配的参数

	hasTimeout bool        //超时标记
	writerMux  *sync.Mutex //写保护
}

func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		request:   r,
		response:  w,
		ctx:       r.Context(),
		index:     -1,
		writerMux: &sync.Mutex{},
	}
}

/**---------------------------------------------**/

func (ctx *Context) BaseContext() context.Context {
	return ctx.request.Context()
}

/**------------------控制器链调用-----------------------------**/
func (ctx *Context) Next() error {
	ctx.index++
	if ctx.index < len(ctx.handlers) {
		if err := ctx.handlers[ctx.index](ctx); err != nil {
			return err
		}
	}
	return nil
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

/**--------------get--set-----------------------**/

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

func (ctx *Context) SetHandlers(handlers []ControllerHandler) {
	ctx.handlers = handlers
}

func (ctx *Context) SetParams(params map[string]string) {
	ctx.params = params
}

/**-----------------------request--------------------------------**/
/******************************************************************/

// param /xxx/xxx/:xx --------------------------start
func (ctx *Context) ParamInt(key string, def int) (int, bool) {
	if val := ctx.Param(key); val != nil {
		return cast.ToInt(val), true
	}
	return def, false
}

func (ctx *Context) ParamInt64(key string, def int64) (int64, bool) {
	if val := ctx.Param(key); val != nil {
		return cast.ToInt64(val), true
	}
	return def, false
}

func (ctx *Context) ParamFloat64(key string, def float64) (float64, bool) {
	if val := ctx.Param(key); val != nil {
		return cast.ToFloat64(val), true
	}
	return def, false
}

func (ctx *Context) ParamFloat32(key string, def float32) (float32, bool) {
	if val := ctx.Param(key); val != nil {
		return cast.ToFloat32(val), true
	}
	return def, false
}

func (ctx *Context) ParamBool(key string, def bool) (bool, bool) {
	if val := ctx.Param(key); val != nil {
		return cast.ToBool(val), true
	}
	return def, false
}

func (ctx *Context) ParamString(key string, def string) (string, bool) {
	if val := ctx.Param(key); val != nil {
		return cast.ToString(val), true
	}
	return def, false
}

// 获取路由参数
func (ctx *Context) Param(key string) interface{} {
	if ctx.params != nil {
		if val, ok := ctx.params[key]; ok {
			return val
		}
	}
	return nil
}

//param--------------------------------------------------end

// form--------------------------------------------------start
func (ctx *Context) FormAll() map[string][]string {
	if ctx.request != nil {
		return map[string][]string(ctx.request.PostForm)
	}
	return map[string][]string{}
}

func (ctx *Context) FormInt(key string, def int) (int, bool) {
	params := ctx.FormAll()
	return handInt(key, def, params)
}

func (ctx *Context) FormInt64(key string, def int64) (int64, bool) {
	params := ctx.FormAll()
	return handInt64(key, def, params)
}

func (ctx *Context) FormFloat64(key string, def float64) (float64, bool) {
	params := ctx.FormAll()
	return handFloat64(key, def, params)
}

func (ctx *Context) FormFloat32(key string, def float32) (float32, bool) {
	params := ctx.FormAll()
	return handFloat32(key, def, params)
}

func (ctx *Context) FormBool(key string, def bool) (bool, bool) {
	params := ctx.FormAll()
	return handBool(key, def, params)
}

func (ctx *Context) FormString(key string, def string) (string, bool) {
	params := ctx.FormAll()
	return handString(key, def, params)
}

func (ctx *Context) FormStringSlice(key string, def []string) ([]string, bool) {
	params := ctx.FormAll()
	return handArray(key, def, params)
}

/** TODO
func (ctx *Context) FormFile(key string) (*multipart.FileHeader, error) {
	if ctx.request.MultipartForm == nil {
		if err := ctx.request.ParseMultipartForm()
	}
}
*/

// --------------------------------------------------------end

// query  xxx/xxx?xx=xx&xx=xx&a[]=xxx --------------------------start
func (ctx *Context) QueryAll() map[string][]string {
	if ctx.request != nil {
		return map[string][]string(ctx.request.URL.Query())
	}
	return map[string][]string{}
}

func (ctx *Context) QueryInt(key string, def int) (int, bool) {
	params := ctx.QueryAll()
	return handInt(key, def, params)
}

func (ctx *Context) QueryInt64(key string, def int64) (int64, bool) {
	params := ctx.QueryAll()
	return handInt64(key, def, params)
}

func (ctx *Context) QueryFloat32(key string, def float32) (float32, bool) {
	params := ctx.QueryAll()
	return handFloat32(key, def, params)
}

func (ctx *Context) QueryFloat64(key string, def float64) (float64, bool) {
	params := ctx.QueryAll()
	return handFloat64(key, def, params)
}

func (ctx *Context) QueryBool(key string, def bool) (bool, bool) {
	params := ctx.QueryAll()
	return handBool(key, def, params)
}

func (ctx *Context) QueryString(key string, def string) (string, bool) {
	params := ctx.QueryAll()
	return handString(key, def, params)
}

func (ctx *Context) QueryArray(key string, def []string) ([]string, bool) {
	params := ctx.QueryAll()
	return handArray(key, def, params)
}

func (ctx *Context) Query(key string) interface{} {
	params := ctx.QueryAll()
	if val, ok := params[key]; ok {
		return val[0]
	}
	return nil
}

// --------------------------------------------------------end

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

/**--------------------------response-----------------------------------**/
/*************************************************************************/

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

//-------------------辅助函数---------------------------------------------

func handInt(key string, def int, params map[string][]string) (int, bool) {
	if val, ok := params[key]; ok {
		l := len(val)
		if l > 0 {
			return cast.ToInt(val[0]), true
		}
	}
	return def, false
}

func handInt64(key string, def int64, params map[string][]string) (int64, bool) {
	if val, ok := params[key]; ok {
		l := len(val)
		if l > 0 {
			return cast.ToInt64(val[0]), true
		}
	}
	return def, false
}

func handFloat64(key string, def float64, params map[string][]string) (float64, bool) {
	if val, ok := params[key]; ok {
		l := len(val)
		if l > 0 {
			return cast.ToFloat64(val[0]), true
		}
	}
	return def, false
}

func handFloat32(key string, def float32, params map[string][]string) (float32, bool) {
	if val, ok := params[key]; ok {
		l := len(val)
		if l > 0 {
			return cast.ToFloat32(val[0]), true
		}
	}
	return def, false
}

func handBool(key string, def bool, params map[string][]string) (bool, bool) {
	if val, ok := params[key]; ok {
		l := len(val)
		if l > 0 {
			return cast.ToBool(val[0]), true
		}
	}
	return def, false
}

func handString(key string, def string, params map[string][]string) (string, bool) {
	if val, ok := params[key]; ok {
		l := len(val)
		if l > 0 {
			return val[l-1], true
		}
	}
	return def, false
}

func handArray(key string, def []string, params map[string][]string) ([]string, bool) {
	if val, ok := params[key]; ok {
		return val, true
	}
	return def, false
}
