package main

import (
	"coredemo/framework"
	"time"
)

func registerRouter(core *framework.Core) {
	core.Get("/user/login", framework.TimeoutHandler(UserLoginController, time.Second))

	subjectApi := core.Group("/subject")
	{
		subjectApi.Delete("/:id", SubjectDelController)
		subjectApi.Put("/:id", SubjectPutController)
		subjectApi.Get("/:id", SubjectGetController)
		subjectApi.Get("/list/all", SubjectListController)

		subjectInnerApi := subjectApi.Group("/info")
		{
			subjectInnerApi.Get("/name", SubjectNameController)
		}
	}
}
