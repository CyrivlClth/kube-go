package main

import (
	"github.com/CyrivlClth/kube-go/app/config"
	"github.com/CyrivlClth/kube-go/app/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	e := gin.Default()
	handler.Register(e, must(config.GetDB()))
	e.Run(":8080")
}

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
