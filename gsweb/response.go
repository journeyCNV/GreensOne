package gsweb

type GResponse interface {

	// 一些输出
	Json(obj interface{}) GResponse
	Jsonp(obj interface{}) GResponse
	Xml(obj interface{}) GResponse
	Html(template string, obj interface{}) GResponse
	Text(format string, values ...interface{}) GResponse

	Redirect(path string) GResponse

	SetHeader(key string, val string) GResponse

	SetCookie(key string, val string, maxAge int, path, domain string, secure, httpOnly bool) GResponse

	SetStatus(code int) GResponse // 设置状态码

	SetOkStatus() GResponse // 设置200状态
}
