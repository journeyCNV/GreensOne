package gsweb

/**
路由分组
*/
type GGroup interface {
	Get(string, ControllerHandler)
	Post(string, ControllerHandler)
	Put(string, ControllerHandler)
	Delete(string, ControllerHandler)
	Group(string) GGroup
}

type Group struct {
	core   *GreensCore
	prefix string //路由前缀
}

func NewGroup(g *GreensCore, prefix string) *Group {
	return &Group{
		core:   g,
		prefix: prefix,
	}
}

func (g *Group) Get(uri string, h ControllerHandler) {
	uri = g.prefix + uri
	g.core.Get(uri, h)
}
func (g *Group) Post(uri string, h ControllerHandler) {
	uri = g.prefix + uri
	g.core.Post(uri, h)
}
func (g *Group) Put(uri string, h ControllerHandler) {
	uri = g.prefix + uri
	g.core.Put(uri, h)
}

func (g *Group) Delete(uri string, h ControllerHandler) {
	uri = g.prefix + uri
	g.core.Delete(uri, h)
}

func (g *Group) Group(uri string) GGroup {
	return NewGroup(g.core, g.prefix+uri)
}
