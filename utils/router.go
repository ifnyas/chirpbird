package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var Router = gin.Default()

func RouterRun() {
	routeStatic()
	routeIndex()
	routeWs()
	Router.Run()
}

func routeStatic() {
	Router.LoadHTMLGlob("static/*.html")
	Router.Static("/static", "./static")
}

func routeIndex() {
	Router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
}

func routeWs() {
	Router.GET("/ws", func(c *gin.Context) {
		user := c.Query("user")
		room := c.Query("room")
		ServeWs(c.Writer, c.Request, room, user)
	})
}
