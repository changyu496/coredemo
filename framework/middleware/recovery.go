package middleware

import "coredemo/framework"

func Recovery() framework.ControllerHandler {
	return func(c *framework.Context) error {
		defer func() {
			if err := recover(); err != nil {
				err := c.Json(500, err)
				if err != nil {
					return
				}
			}
		}()
		err := c.Next()
		if err != nil {
			err := c.Json(500, err)
			if err != nil {
				return nil
			}
		}
		return nil
	}
}
