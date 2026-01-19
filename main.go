package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Hello World")
	})

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
