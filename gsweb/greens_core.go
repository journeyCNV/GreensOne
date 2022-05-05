package gsweb

import (
	"github.com/journeycnv/greensone/gsweb/container"
	"log"
	"net/http"
	"strings"
)

/**
实现http标准库的handler接口
*/
type GreensCore struct {
	router      GRouter       // 所有路由
	middlewares []HandlerFunc // 中间件
	con         container.GContainer
}

type GRouter map[string]*TrieTree // 方法 : 路由 : 实际处理函数

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

func NewGreensCore() *GreensCore {
	return &GreensCore{
		router: routerMap(),
		con:    container.NewContainer(),
	}
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
	ctx.con = g.con

	// 寻找路由匹配的控制器们
	node := g.FindRouteNode(request)
	if node == nil {
		ctx.Json("not found").SetStatus(http.StatusBadRequest)
		return
	}
	ctx.SetHandlers(node.handlers)

	// 设置路由参数
	params := node.parseParamsFormEndNode(request.URL.Path)
	ctx.SetParams(params)

	// 调用路由函数, 访问控制器链上的函数
	if err := ctx.Next(); err != nil {
		ctx.Json("inner error").SetStatus(http.StatusInternalServerError)
		return
	}
}

// 注册中间件Use
func (g *GreensCore) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}

// 匹配路由Node
func (g *GreensCore) FindRouteNode(req *http.Request) *node {
	uri := req.URL.Path
	realUri := strings.ToLower(uri)
	method := req.Method
	realMethod := strings.ToUpper(method)

	if treeHandler, ok := g.router[realMethod]; ok {
		return treeHandler.root.matchNode(realUri)
	}
	return nil
}

func (g *GreensCore) Get(url string, h ...HandlerFunc) {
	handlers := append(g.middlewares, h...)
	if err := g.router[GET].AddRouter(url, handlers); err != nil {
		log.Fatal("add router error:", err)
	}
}

func (g *GreensCore) Post(url string, h ...HandlerFunc) {
	handlers := append(g.middlewares, h...)
	if err := g.router[POST].AddRouter(url, handlers); err != nil {
		log.Fatal("add router error:", err)
	}
}

func (g *GreensCore) Put(url string, h ...HandlerFunc) {
	handlers := append(g.middlewares, h...)
	if err := g.router[PUT].AddRouter(url, handlers); err != nil {
		log.Fatal("add router error:", err)
	}
}

func (g *GreensCore) Delete(url string, h ...HandlerFunc) {
	handlers := append(g.middlewares, h...)
	if err := g.router[DELETE].AddRouter(url, handlers); err != nil {
		log.Fatal("add router error:", err)
	}
}

func (g *GreensCore) Group(prefix string) GGroup {
	return NewGroup(g, prefix)
}

//--------------------容器相关-------------------------------------

func (g *GreensCore) Bind(p container.ServiceProvider) error {
	return g.con.Bind(p)
}

func (g *GreensCore) IsBind(key string) bool {
	return g.con.IsBind(key)
}

//-----------------------------------------------------------------

func routerMap() GRouter {
	router := GRouter{}
	router[GET] = NewTrieTree()
	router[POST] = NewTrieTree()
	router[PUT] = NewTrieTree()
	router[DELETE] = NewTrieTree()
	return router
}
