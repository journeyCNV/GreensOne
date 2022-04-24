package gsweb

/**
路由分组
*/
type GGroup interface {
	Get(string, ...ControllerHandler)
	Post(string, ...ControllerHandler)
	Put(string, ...ControllerHandler)
	Delete(string, ...ControllerHandler)
	Group(string) GGroup                  //支持多层group
	Use(middlewares ...ControllerHandler) //嵌套中间件
}

type Group struct {
	core   *GreensCore
	prefix string //路由前缀
	parent *Group // 指向上一级路由，方便控制整个Group共用的中间件，从Group级别加中间件

	middlewares []ControllerHandler
}

func NewGroup(g *GreensCore, prefix string) *Group {
	return &Group{
		core:        g,
		prefix:      prefix,
		middlewares: []ControllerHandler{},
	}
}

// 注册中间件
func (g *Group) Use(middlewares ...ControllerHandler) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (g *Group) Get(uri string, h ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	handlers := append(g.getMiddlewares(), h...)
	g.core.Get(uri, handlers...)
}

func (g *Group) Post(uri string, h ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	handlers := append(g.getMiddlewares(), h...)
	g.core.Post(uri, handlers...)
}
func (g *Group) Put(uri string, h ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	handlers := append(g.getMiddlewares(), h...)
	g.core.Put(uri, handlers...)
}

func (g *Group) Delete(uri string, h ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	handlers := append(g.getMiddlewares(), h...)
	g.core.Delete(uri, handlers...)
}

func (g *Group) Group(uri string) GGroup {
	group := NewGroup(g.core, uri)
	group.parent = g
	return group
}

/**
func (g *Group) Group(uri string) GGroup {
	return NewGroup(g.core, g.prefix+uri)
}
*/

func (g *Group) getMiddlewares() []ControllerHandler {
	if g.parent == nil {
		return g.middlewares
	}
	return append(g.parent.getMiddlewares(), g.middlewares...)
}

func (g *Group) getAbsolutePrefix() string {
	if g.parent == nil {
		return g.prefix
	}
	return g.parent.getAbsolutePrefix() + g.prefix
}
