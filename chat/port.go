package chat

import (
	"log"
	"net"
	"net/http"
	"strconv"
)

var (
	ports         = []string{}
	lastPortIndex = 0
)

func setPorts(n int) {
	if n < 1 {
		n = 1
	}

	portsTemp := make([]string, n)
	for i := range portsTemp {
		if i < len(ports) {
			portsTemp[i] = ports[i]
		}
	}
	ports = portsTemp

	for i := 0; i < n; i++ {
		if ports[i] == "" {
			err := addPort(i)
			if err != nil {
				continue
			}
		}
	}
}

func addPort(i int) error {
	freePort, err := GetFreePort()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	newPort := ":" + strconv.Itoa(freePort)
	ports[i] = newPort
	return nil
}

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", ":0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func healthCheck(i int) bool {
	url := domain + ports[i] + "/chat"
	res, err := http.Get(url)

	if err != nil {
		log.Println(err.Error())
	}

	if res.StatusCode != http.StatusOK {
		log.Println(res.Status)
		return false
	}

	return true
}

func healthyPort() string {
	nextPort := ""

	for i := 0; i <= len(ports); i++ {
		nextPortIndex := (lastPortIndex + 1) % len(ports)
		if i == len(ports) {
			runNewPort()
		} else {
			isPortHealthy := healthCheck(nextPortIndex)
			if isPortHealthy {
				lastPortIndex = nextPortIndex
				nextPort = ports[nextPortIndex]
				break
			}
		}
	}

	return nextPort
}

func runNewPort() {
	index := len(ports)
	ports = append(ports, "")
	err := addPort(index)

	if err != nil {
		log.Println(err.Error())
		startPanicPlan()
	}

	addEngine(ports[index])
}

func startPanicPlan() {
	log.Println("lol")
}
