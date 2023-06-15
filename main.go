package main

import (
	"coredemo/framework"
	"fmt"
	"net/http"
)

func main() {
	core := framework.NewCore()
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
