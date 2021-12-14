package wss

import (
	"log"
	"net"
	"strconv"
)

var (
	ports = []string{":230"}
)

func setPorts(n int) {
	num := n
	if n < 2 {
		num = 2
	}

	portsTemp := make([]string, num)
	for i := range portsTemp {
		if i < len(ports) {
			portsTemp[i] = ports[i]
		}
	}
	ports = portsTemp

	for i := 0; i < num; i++ {
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
