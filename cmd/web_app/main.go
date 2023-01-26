package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// var db = make(map[string]string)

func main() {
	server := gin.Default()

	server.LoadHTMLGlob("web/templates/*.html")
	server.Static("/figures", "./figures")
	server.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"message": "OK!",
		})
	})
	server.Run(":8080")
}
