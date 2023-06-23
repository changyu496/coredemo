package middleware

import (
	"coredemo/framework"
	"fmt"
)

func Test1() framework.ControllerHandler {
	return func(c *framework.Context) error {
		fmt.Println("middle pre test1")
		c.Next()
		fmt.Println("middle post test1")
		return nil
	}
}

func Test2() framework.ControllerHandler {
	return func(c *framework.Context) error {
		fmt.Println("middle pre test2")
		c.Next()
		fmt.Println("middle post test2")
		return nil
	}
}
