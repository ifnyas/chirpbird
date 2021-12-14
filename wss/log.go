package wss

import (
	"encoding/json"
	"time"
)

var (
	logsTemp = []Message{}
)

func saveMsg(m Message) {
	logsTemp = append(logsTemp, m)
}

func loadMsg(room string) []Message {
	selected := []Message{}
	for _, log := range logsTemp {
		if log.Room == room {
			selected = append(selected, log)
		}
	}
	return selected
}

func modMsg(msg []byte, user string) []byte {
	res := Response{user, string(msg), time.Now().Format("15:04")}
	json, err := json.Marshal(res)
	if err != nil {
		json = []byte("")
	}
	return json
}
