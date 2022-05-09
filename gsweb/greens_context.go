package gsweb

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/journeycnv/greensone/cast"
	"github.com/journeycnv/greensone/gsweb/container"
	"html/template"
	"io"
	"io/ioutil"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

type Context struct {
	request  *http.Request
	response http.ResponseWriter
	ctx      context.Context

	handlers []HandlerFunc     //当前请求的控制器链
	index    int               //当前请求控制器链中下标
	params   map[string]string //路由通配符匹配的参数

	hasTimeout bool        //超时标记
	writerMux  *sync.Mutex //写保护

	Errors errorMsgs
	con    container.GContainer
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

const abortIndex int8 = math.MaxInt8 / 2

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

func (ctx *Context) SetHandlers(handlers []HandlerFunc) {
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
		ctx.request.ParseForm()
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

func (ctx *Context) FormFile(key string) (*multipart.FileHeader, error) {
	if ctx.request.MultipartForm == nil { // 如果没有设置这个大小
		// 如果上传的文件大小大于maxMemory,将存在临时文件里  见ctx.request.ParseMultipartForm源码
		if err := ctx.request.ParseMultipartForm(defaultMultipartMemory); err != nil {
			return nil, err
		}
	}
	f, fh, err := ctx.request.FormFile(key)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, err
}

func (ctx *Context) Form(key string) interface{} {
	params := ctx.FormAll()
	if val, ok := params[key]; ok {
		if len(val) > 0 {
			return val[0]
		}
	}
	return nil
}

func (ctx *Context) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

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

func (ctx *Context) QueryStringSlice(key string, def []string) ([]string, bool) {
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

// 解析body为object
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

func (ctx *Context) BindXml(obj interface{}) error {
	if ctx.request != nil {
		body, err := ioutil.ReadAll(ctx.request.Body)
		if err != nil {
			return err
		}
		ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		err = xml.Unmarshal(body, obj)
		if err != nil {
			return err
		}
	} else {
		return errors.New("ctx.request empty")
	}
	return nil
}

// 其他格式
func (ctx *Context) GetRawData() ([]byte, error) {
	if ctx.request != nil {
		body, err := ioutil.ReadAll(ctx.request.Body)
		if err != nil {
			return nil, err
		}
		ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		return body, err
	}
	return nil, errors.New("ctx.request empty")
}

func (ctx *Context) Uri() string {
	return ctx.request.RequestURI
}

func (ctx *Context) Method() string {
	return ctx.request.Method
}

func (ctx *Context) Host() string {
	return ctx.request.URL.Host
}

func (ctx *Context) ClientIp() string {
	r := ctx.request
	ipAddress := r.Header.Get("X-Real-Ip") // 真实客户端Ip
	if ipAddress == "" {
		ipAddress = r.Header.Get("X-Forwarded-For") //代理信息
	}
	if ipAddress == "" {
		ipAddress = r.RemoteAddr // 上一个存在的客户端的地址或者上一个代理服务器地址
	}
	return ipAddress
}

func (ctx *Context) Headers() map[string][]string {
	return map[string][]string(ctx.request.Header)
}

func (ctx *Context) Header(key string) (string, bool) {
	val := ctx.request.Header.Values(key)
	if val == nil || len(val) <= 0 {
		return "", false
	}
	return val[0], true
}

func (ctx *Context) Cookies() map[string]string {
	cookies := ctx.request.Cookies()
	ret := map[string]string{}
	for _, cookie := range cookies {
		ret[cookie.Name] = cookie.Value
	}
	return ret
}

func (ctx *Context) Cookie(key string) (string, bool) {
	cookies := ctx.Cookies()
	if val, ok := cookies[key]; ok {
		return val, true
	}
	return "", false
}

/**--------------------------response-----------------------------------**/
/*************************************************************************/
func (ctx *Context) Jsonp(obj interface{}) GResponse {
	// 获取请求参数
	callbackFunc, _ := ctx.QueryString("callback", "callback_function")
	ctx.SetHeader("Context-Type", "application/javascript")
	// 输出到前端要进行字符过滤，否则可能造成xss攻击
	callback := template.JSEscapeString(callbackFunc)

	// 输出函数名
	_, err := ctx.response.Write([]byte(callback))
	if err != nil {
		return ctx
	}

	// 输出左括号
	_, err = ctx.response.Write([]byte("("))
	if err != nil {
		return ctx
	}

	ret, err := json.Marshal(obj)
	if err != nil {
		return ctx
	}

	_, err = ctx.response.Write(ret)
	if err != nil {
		return ctx
	}

	_, err = ctx.response.Write([]byte(")"))
	if err != nil {
		return ctx
	}

	return ctx

}

func (ctx *Context) Json(obj interface{}) GResponse {
	byteObj, err := json.Marshal(obj)
	if err != nil {
		return ctx.SetStatus(http.StatusInternalServerError)
	}
	ctx.SetHeader("Content-Type", "application/json")
	ctx.response.Write(byteObj)
	return ctx
}

func (ctx *Context) Xml(obj interface{}) GResponse {
	byt, err := xml.Marshal(obj)
	if err != nil {
		return ctx.SetStatus(http.StatusInternalServerError)
	}
	ctx.SetHeader("Content-Type", "application/xml")
	ctx.response.Write(byt)
	return ctx
}

func (ctx *Context) Html(file string, obj interface{}) GResponse {
	// 读取模板文件，创建template实例
	t, err := template.New("output").ParseFiles(file)
	if err != nil {
		return ctx
	}
	// 将obj和模板结合
	if err := t.Execute(ctx.response, obj); err != nil {
		return ctx
	}
	ctx.SetHeader("Content-Type", "application/html")
	return ctx
}

func (ctx *Context) Text(format string, values ...interface{}) GResponse {
	out := fmt.Sprintf(format, values...)
	ctx.SetHeader("Content-Type", "application/text")
	ctx.response.Write([]byte(out))
	return ctx
}

// 301 重定向
func (ctx *Context) Redirect(path string) GResponse {
	http.Redirect(ctx.response, ctx.request, path, http.StatusMovedPermanently)
	return ctx
}

func (ctx *Context) SetHeader(key string, val string) GResponse {
	ctx.response.Header().Add(key, val)
	return ctx
}

func (ctx *Context) SetCookie(key string, val string, maxAge int, path string, domain string, secure bool, httpOnly bool) GResponse {
	if path == "" {
		path = "/"
	}
	http.SetCookie(ctx.response, &http.Cookie{
		Name:     key,
		Value:    url.QueryEscape(val),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: 1,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
	return ctx
}

func (ctx *Context) SetStatus(code int) GResponse {
	ctx.response.WriteHeader(code)
	return ctx
}

func (ctx *Context) SetOkStatus() GResponse {
	ctx.response.WriteHeader(http.StatusOK)
	return ctx
}

//----------------------------------------------------------------
// Abort 可防止调用挂起的处理程序
func (ctx *Context) Abort() {
	ctx.index = int(abortIndex)
}

func (ctx *Context) IsAborted() bool {
	return ctx.index >= int(abortIndex)
}

func (ctx *Context) AbortWithStatus(code int) {
	ctx.SetStatus(code)
	ctx.Abort()
}

//----------------------------------------------------------------

func (ctx *Context) Error(err error) *Error {
	if err == nil {
		panic("err is nil")
	}

	parsedError, ok := err.(*Error)
	if !ok {
		parsedError = &Error{
			Err:  err,
			Type: ErrorTypePrivate,
		}
	}

	ctx.Errors = append(ctx.Errors, parsedError)
	return parsedError
}

//------------------容器-------------------------------------------------

func (ctx *Context) Make(key string) (interface{}, error) {
	return ctx.con.Make(key)
}

func (ctx *Context) MustMake(key string) interface{} {
	return ctx.con.MustMake(key)
}

func (ctx *Context) MakeNew(key string, params []interface{}) (interface{}, error) {
	return ctx.con.MakeNew(key, params)
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
