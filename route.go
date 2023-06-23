package main

import (
	"coredemo/framework"
	"coredemo/framework/middleware"
)

func registerRouter(core *framework.Core) {
	core.Get("/user/login", middleware.Test1(), UserLoginController)

	subjectApi := core.Group("/subject")
	{
		subjectApi.Delete("/:id", SubjectDelController)
		subjectApi.Put("/:id", SubjectPutController)
		subjectApi.Get("/:id", SubjectGetController)
		subjectApi.Get("/list/all", SubjectListController)

		subjectInnerApi := subjectApi.Group("/info")
		{
			subjectInnerApi.Get("/name", middleware.Test1(), SubjectNameController)
		}
	}
}
