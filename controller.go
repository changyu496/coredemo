package main

import (
	"context"
	"coredemo/framework"
	"fmt"
	"log"
	"time"
)

func FooControllerHandler(ctx *framework.Context) error {
	durationCtx, cancel := context.WithTimeout(ctx.BaseContext(), time.Duration(1*time.Second))
	defer cancel()
	// 这个channel负责通知结束
	finish := make(chan struct{}, 1)
	// 这个channel负责通知 panic 异常
	panicChan := make(chan interface{}, 1)
	go func() {
		// 增加异常
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		// 具体业务
		time.Sleep(10 * time.Second)
		ctx.Json(200, map[string]interface{}{
			"code": 0,
		})
		finish <- struct{}{}
	}()
	select {
	case p := <-panicChan:
		ctx.WriteMux().Lock()
		defer ctx.WriteMux().Unlock()
		log.Println(p)
		ctx.Json(500, "panic")
	case <-finish:
		fmt.Println("finish")
	case <-durationCtx.Done():
		ctx.WriteMux().Lock()
		defer ctx.WriteMux().Unlock()
		ctx.Json(500, "time out")
		ctx.SetHasTimeout()
	}
	return nil
}
