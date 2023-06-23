package framework

import (
	"log"
	"net/http"
	"strings"
)

type Core struct {
	router      map[string]*Tree
	middlewares []ControllerHandler
}

func NewCore() *Core {
	router := map[string]*Tree{}
	router["GET"] = NewTree()
	router["POST"] = NewTree()
	router["PUT"] = NewTree()
	router["DELETE"] = NewTree()

	return &Core{router: router}
}

func (c *Core) Get(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	err := c.router["GET"].AddRoute(url, allHandlers)
	if err != nil {
		log.Fatal("add router error:", err)
	}
}

func (c *Core) Post(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	err := c.router["POST"].AddRoute(url, allHandlers)
	if err != nil {
		log.Fatal("add router error", err)
	}
}

func (c *Core) Put(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	err := c.router["PUT"].AddRoute(url, allHandlers)
	if err != nil {
		log.Fatal("add router error", err)
	}
}

func (c *Core) Delete(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	err := c.router["DELETE"].AddRoute(url, allHandlers)
	if err != nil {
		log.Fatal("add router error", err)
	}
}

func (c *Core) Use(middlewares ...ControllerHandler) {
	c.middlewares = append(c.middlewares, middlewares...)
}

func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}

func (c *Core) FindRoutesByRequest(request *http.Request) []ControllerHandler {
	uri := request.URL.Path
	method := request.Method
	upperMethod := strings.ToUpper(method)

	if methodHandlers, ok := c.router[upperMethod]; ok {
		return methodHandlers.FindHandler(uri)
	}
	return nil
}

func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	log.Println("core.ServeHTTP")
	ctx := NewContext(request, response)

	handlers := c.FindRoutesByRequest(request)
	if handlers == nil {
		err := ctx.Json(404, "not found")
		if err != nil {
			return
		}
		return
	}
	ctx.SetHandlers(handlers)
	// 调用路由函数，如果返回err，代表存在内部错误，返回500状态码
	if err := ctx.Next(); err != nil {
		err := ctx.Json(500, "inner error")
		if err != nil {
			return
		}
	}
}
