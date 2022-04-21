package gsweb

import "net/http"

/**
实现http标准库的handler接口
*/
type GreensCore struct {
}

func NewGreensCore() *GreensCore {
	return &GreensCore{}
}

// 实现http标准库的handler的方法
func (g *GreensCore) ServeHTTP(response http.ResponseWriter, request *http.Request) {

}
