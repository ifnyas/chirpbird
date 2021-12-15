package router

import (
	"chirpbird/chat"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Run() {
	// init
	var r = gin.New()

	// static
	r.LoadHTMLGlob("./static/*.html")
	r.Static("/static", "./static")

	// index
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// chat
	r.Any("/chat/*path", chat.Proxy())

	// run
	r.Run()
}
