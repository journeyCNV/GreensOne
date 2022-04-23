package gsweb

import (
	"log"
	"net/http"
	"strings"
)

/**
实现http标准库的handler接口
*/
type GreensCore struct {
	router GRouter
}

type GRouter map[string]*TrieTree // 方法 : 路由 : 实际处理函数

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

	if treeHandler, ok := g.router[realMethod]; ok {
		return treeHandler.FindHandler(realUri)
	}
	return nil
}

func (g *GreensCore) Get(url string, h ControllerHandler) {
	if err := g.router[GET].AddRouter(url, h); err != nil {
		log.Fatal("add router error:", err)
	}
}

func (g *GreensCore) Post(url string, h ControllerHandler) {
	if err := g.router[POST].AddRouter(url, h); err != nil {
		log.Fatal("add router error:", err)
	}
}

func (g *GreensCore) Put(url string, h ControllerHandler) {
	if err := g.router[PUT].AddRouter(url, h); err != nil {
		log.Fatal("add router error:", err)
	}
}

func (g *GreensCore) Delete(url string, h ControllerHandler) {
	if err := g.router[DELETE].AddRouter(url, h); err != nil {
		log.Fatal("add router error:", err)
	}
}

func (g *GreensCore) Group(prefix string) GGroup {
	return NewGroup(g, prefix)
}

func routerMap() GRouter {
	router := GRouter{}
	router[GET] = NewTrieTree()
	router[POST] = NewTrieTree()
	router[PUT] = NewTrieTree()
	router[DELETE] = NewTrieTree()
	return router
}
