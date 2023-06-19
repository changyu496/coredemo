package main

import (
	"coredemo/framework"
)

func UserLoginController(c *framework.Context) error {
	//time.Sleep(time.Second * 10)/**/
	c.Json(200, "ok,UserLoginController")
	return nil
}
