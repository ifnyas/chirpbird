package chat

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"net/http/httputil"
// 	"net/url"
// 	"time"

// 	"github.com/go-co-op/gocron"
// )

// type server struct {
// 	Name         string
// 	URL          string
// 	ReverseProxy *httputil.ReverseProxy
// 	Health       bool
// }

// func newServer(name, urlStr string) *server {
// 	u, _ := url.Parse(urlStr)
// 	rp := httputil.NewSingleHostReverseProxy(u)
// 	return &server{
// 		Name:         name,
// 		URL:          urlStr,
// 		ReverseProxy: rp,
// 		Health:       true,
// 	}
// }

// func (s *server) checkHealth() bool {
// 	resp, err := http.Head(s.URL)
// 	if err != nil {
// 		s.Health = false
// 		return s.Health
// 	}
// 	if resp.StatusCode != http.StatusOK {
// 		s.Health = false
// 		return s.Health
// 	}
// 	s.Health = true
// 	return s.Health
// }

// /*
// Load Balancer.go
// */
// var (
// 	serverList = []*server{
// 		newServer("server-1", "http://127.0.0.1:5001"),
// 		newServer("server-2", "http://127.0.0.1:5002"),
// 		newServer("server-3", "http://127.0.0.1:5003"),
// 		newServer("server-4", "http://127.0.0.1:5004"),
// 		newServer("server-5", "http://127.0.0.1:5005"),
// 	}
// 	lastServedIndex = 0
// )

// func main() {
// 	http.HandleFunc("/", forwardRequest)
// 	go startHealthCheck()
// 	log.Fatal(http.ListenAndServe(":8000", nil))
// }

// func forwardRequest(res http.ResponseWriter, req *http.Request) {
// 	server, err := getHealthyServer()
// 	if err != nil {
// 		http.Error(res, "Couldn't process request: "+err.Error(), http.StatusServiceUnavailable)
// 		return
// 	}
// 	server.ReverseProxy.ServeHTTP(res, req)
// }

// func getHealthyServer() (*server, error) {
// 	for i := 0; i < len(serverList); i++ {
// 		server := getServer()
// 		if server.Health {
// 			return server, nil
// 		}
// 	}
// 	return nil, fmt.Errorf("No healthy hosts")
// }

// func getServer() *server {
// 	nextIndex := (lastServedIndex + 1) % len(serverList)
// 	server := serverList[nextIndex]
// 	lastServedIndex = nextIndex
// 	return server
// }

// /*
// HealthCheck.go
// */

// func startHealthCheck() {
// 	s := gocron.NewScheduler(time.Local)
// 	for _, host := range serverList {
// 		_, err := s.Every(2).Seconds().Do(func(s *server) {
// 			healthy := s.checkHealth()
// 			if healthy {
// 				log.Printf("'%s' is healthy!", s.Name)
// 			} else {
// 				log.Printf("'%s' is not healthy", s.Name)
// 			}
// 		}, host)
// 		if err != nil {
// 			log.Fatalln(err)
// 		}
// 	}
// 	<-s.StartAsync()
// }
