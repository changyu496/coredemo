package main

import (
	"coredemo/framework"
	"coredemo/framework/middleware"
	"fmt"
	"net/http"
)

func main() {
	core := framework.NewCore()
	core.Use(middleware.Recovery())
	registerRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    ":8080",
	}
	fmt.Println("Started......")
	err := server.ListenAndServe()
	if err != nil {
		return
	}
}
