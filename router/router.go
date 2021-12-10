package router

import (
	"chirpbird/wss"
	"encoding/json"
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

	r.GET("/ws", func(c *gin.Context) {
		wss.ServeWs(c.Writer, c.Request, c.Query("key"))
	})

	r.GET("/history", func(c *gin.Context) {
		roomHistories := wss.LoadMsg(c.Query("room"))
		array := []wss.Response{}
		for _, history := range roomHistories {
			var res wss.Response
			json.Unmarshal(history.Data, &res)
			array = append(array, res)
		}
		c.JSON(http.StatusOK, gin.H{"data": array})
	})

	// run
	r.Run()
}
