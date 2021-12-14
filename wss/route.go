package wss

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/gin-gonic/gin"
)

func setRoute() {
	r := gin.Default()

	r.GET("/ws", func(c *gin.Context) {
		ServeWs(c.Writer, c.Request, c.Query("key"))
	})

	r.GET("/logs", func(c *gin.Context) {
		roomLogs := loadMsg(c.Query("room"))
		array := []Response{}
		for _, log := range roomLogs {
			var res Response
			json.Unmarshal(log.Data, &res)
			array = append(array, res)
		}
		c.JSON(http.StatusOK, gin.H{"data": array})
	})

	for _, port := range ports {
		go r.Run(port)
	}
}

func WsProxy() gin.HandlerFunc {
	return func(c *gin.Context) {
		scheme := c.Request.URL.Scheme
		if scheme == "" {
			scheme = "http"
		}

		host := c.Request.Host
		if strings.Contains(host, ":") {
			host = strings.Split(host, ":")[0]
		}

		director := func(req *http.Request) {
			req.URL.Scheme = scheme
			req.URL.Host = host + ports[0]
			req.Header["my-header"] = []string{c.Request.Header.Get("my-header")}
			delete(req.Header, "My-Header")
		}

		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
