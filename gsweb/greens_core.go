package gsweb

import (
	"net/http"
	"strings"
)

/**
实现http标准库的handler接口
*/
type GreensCore struct {
	router GRouter
}

type GRouter map[string]map[string]ControllerHandler // 方法 : 路由 : 实际处理函数

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

func NewGreensCore() *GreensCore {
	return &GreensCore{router: routerMap()}
}

// 实现http标准库的handler的方法
/**
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
*/
func (g *GreensCore) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request, response)
	handler := g.FindRouteHandler(request)
	if handler == nil {
		ctx.Json(404, "not found")
		return
	}
	if err := handler(ctx); err != nil {
		ctx.Json(500, "inner error")
		return
	}
}

func (g *GreensCore) FindRouteHandler(req *http.Request) ControllerHandler {
	uri := req.URL.Path
	realUri := strings.ToLower(uri)
	method := req.Method
	realMethod := strings.ToUpper(method)

	if handlerMap, ok := g.router[realMethod]; ok {
		if handler, ok := handlerMap[realUri]; ok {
			return handler
		}
	}
	return nil
}

func (g *GreensCore) Get(url string, h ControllerHandler) {
	realUrl := strings.ToLower(url)
	g.router[GET][realUrl] = h
}

func (g *GreensCore) Post(url string, h ControllerHandler) {
	realUrl := strings.ToLower(url)
	g.router[POST][realUrl] = h
}

func (g *GreensCore) Put(url string, h ControllerHandler) {
	realUrl := strings.ToLower(url)
	g.router[PUT][realUrl] = h
}

func (g *GreensCore) Delete(url string, h ControllerHandler) {
	realUrl := strings.ToLower(url)
	g.router[DELETE][realUrl] = h
}

func (g *GreensCore) Group(prefix string) GGroup {
	return NewGroup(g, prefix)
}

func routerMap() GRouter {
	getRouters := map[string]ControllerHandler{}
	postRouters := map[string]ControllerHandler{}
	putRouters := map[string]ControllerHandler{}
	deleteRouters := map[string]ControllerHandler{}
	router := GRouter{}
	router[GET] = getRouters
	router[POST] = postRouters
	router[PUT] = putRouters
	router[DELETE] = deleteRouters
	return router
}
