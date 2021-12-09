package router

import (
	"chirpbird/wss"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Run() {
	// init
	var r = gin.Default()

	// static
	r.LoadHTMLGlob("static/*.html")
	r.Static("/static", "./static")

	// index
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// ws port
	r.GET("/ws", wss.WsProxy())

	// run
	r.Run()
}