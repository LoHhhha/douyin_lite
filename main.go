package main

import (
	"douyin_lite/repository"
	"douyin_lite/service"
	"douyin_lite/settings"
	"douyin_lite/tools"
	"github.com/gin-gonic/gin"
)

func main() {
	go service.RunMessageServer()

	r := gin.Default()

	initRouter(r)

	err := repository.InitSql()
	if err != nil {
		tools.ErrorPrint(err)
		return
	}
	defer repository.Database.Close()

	r.Run(settings.ServerIP + settings.ServerPort) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
