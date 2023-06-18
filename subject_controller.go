package main

import "coredemo/framework"

func SubjectDelController(c *framework.Context) error {
	c.Json(200, "SubjectDelController")
	return nil
}

func SubjectGetController(c *framework.Context) error {
	c.Json(200, "SubjectGetController")
	return nil
}

func SubjectPostController(c *framework.Context) error {
	c.Json(200, "SubjectPostController")
	return nil
}

func SubjectPutController(c *framework.Context) error {
	c.Json(200, "SubjectPutController")
	return nil
}

func SubjectListController(c *framework.Context) error {
	c.Json(200, "SubjectListController")
	return nil
}

func SubjectNameController(c *framework.Context) error {
	c.Json(200, "SubjectNameController")
	return nil
}
