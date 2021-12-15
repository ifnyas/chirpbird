package chat

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	r      = gin.Default()
	domain = ""
)

func initRoute() {
	setRoute()
	for _, port := range ports {
		addEngine(port)
	}
}

func setRoute() {
	r.GET("/chat", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	r.GET("/chat/ws", func(c *gin.Context) {
		ServeWs(c.Writer, c.Request, c.Query("key"))
	})

	r.GET("/chat/logs", func(c *gin.Context) {
		roomLogs := loadMsg(c.Query("room"))
		c.JSON(http.StatusOK, gin.H{"data": roomLogs})
	})
}

func addEngine(port string) {
	go r.Run(port)
}

func Proxy() gin.HandlerFunc {
	return func(c *gin.Context) {
		scheme := request("scheme", c.Request.URL.Scheme)
		host := request("host", c.Request.Host)
		setDomain(scheme, host)

		port := request("port", "")
		path := request("path", c.Param("path"))

		director := func(req *http.Request) {
			req.Header = c.Request.Header
			req.URL.Scheme = scheme
			req.URL.Host = host + port
			req.URL.Path = path
		}

		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)

	}
}

func request(part string, url string) string {
	switch part {
	case "scheme":
		if url == "" {
			url = "http"
		}
		fmt.Println("SSSCHEMEE", url)
	case "host":
		if strings.Contains(url, ":") {
			url = strings.Split(url, ":")[0]
		}
		fmt.Println("HOOOSSSTT", url)
	case "port":
		url = healthyPort()
		fmt.Println("POORRTT", url)
	case "path":
		url = "/chat" + url
		fmt.Println("PAATTHH", url)
	}
	return url
}

func setDomain(scheme string, host string) {
	domain = scheme + "://" + host
}
