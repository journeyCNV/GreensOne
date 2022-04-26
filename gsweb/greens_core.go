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
	router      GRouter             // 所有路由
	middlewares []ControllerHandler // 中间件
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
所有的请求都会进入到这个函数，这个函数负责路由分发
*/
func (g *GreensCore) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request, response)

	// 寻找路由
	handlers := g.FindRouteHandler(request)
	if handlers == nil {
		ctx.Json(404, "not found")
		return
	}
	ctx.SetHandlers(handlers)

	// 设置路由参数
	/**
	params := node.parseParamsFormEndNode(request.URL.Path)
	*/

	// 调用路由函数, 访问控制器链上的函数
	if err := ctx.Next(); err != nil {
		ctx.Json(500, "inner error")
		return
	}
}

// 注册中间件Use
func (g *GreensCore) Use(middlewares ...ControllerHandler) {
	g.middlewares = append(g.middlewares, middlewares...)
}

// 匹配路由，如果没有匹配到返回nil
func (g *GreensCore) FindRouteHandler(req *http.Request) []ControllerHandler {
	uri := req.URL.Path
	realUri := strings.ToLower(uri)
	method := req.Method
	realMethod := strings.ToUpper(method)

	if treeHandler, ok := g.router[realMethod]; ok {
		return treeHandler.FindHandler(realUri)
	}
	return nil
}

func (g *GreensCore) Get(url string, h ...ControllerHandler) {
	handlers := append(g.middlewares, h...)
	if err := g.router[GET].AddRouter(url, handlers); err != nil {
		log.Fatal("add router error:", err)
	}
}

func (g *GreensCore) Post(url string, h ...ControllerHandler) {
	handlers := append(g.middlewares, h...)
	if err := g.router[POST].AddRouter(url, handlers); err != nil {
		log.Fatal("add router error:", err)
	}
}

func (g *GreensCore) Put(url string, h ...ControllerHandler) {
	handlers := append(g.middlewares, h...)
	if err := g.router[PUT].AddRouter(url, handlers); err != nil {
		log.Fatal("add router error:", err)
	}
}

func (g *GreensCore) Delete(url string, h ...ControllerHandler) {
	handlers := append(g.middlewares, h...)
	if err := g.router[DELETE].AddRouter(url, handlers); err != nil {
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
