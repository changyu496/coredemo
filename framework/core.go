package framework

import (
	"log"
	"net/http"
	"strings"
)

type Core struct {
	router map[string]*Tree
}

func NewCore() *Core {
	router := map[string]*Tree{}
	router["GET"] = NewTree()
	router["POST"] = NewTree()
	router["PUT"] = NewTree()
	router["DELETE"] = NewTree()

	return &Core{router: router}
}

func (c *Core) Get(url string, handler ControllerHandler) {
	err := c.router["GET"].AddRoute(url, handler)
	if err != nil {
		log.Fatal("add router error:", err)
	}
}

func (c *Core) Post(url string, handler ControllerHandler) {
	err := c.router["POST"].AddRoute(url, handler)
	if err != nil {
		log.Fatal("add router error", err)
	}
}

func (c *Core) Put(url string, handler ControllerHandler) {
	err := c.router["PUT"].AddRoute(url, handler)
	if err != nil {
		log.Fatal("add router error", err)
	}
}

func (c *Core) Delete(url string, handler ControllerHandler) {
	err := c.router["DELETE"].AddRoute(url, handler)
	if err != nil {
		log.Fatal("add router error", err)
	}
}

func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}

func (c *Core) FindRoutesByRequest(request *http.Request) ControllerHandler {
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

	router := c.FindRoutesByRequest(request)
	if router == nil {
		err := ctx.Json(404, "not found")
		if err != nil {
			return
		}
		return
	}
	err := router(ctx)
	if err != nil {
		err := ctx.Json(500, "inner error")
		if err != nil {
			return
		}
	}
}
