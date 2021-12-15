package chat

import (
	"encoding/json"
	"time"
)

var (
	logsTemp = []Message{}
	//queuesTemp = []Message{}
)

func saveMsg(m Message) {
	logsTemp = append(logsTemp, m)
}

func loadMsg(room string) []Response {
	selected := []Message{}
	for _, log := range logsTemp {
		if log.Room == room {
			selected = append(selected, log)
		}
	}

	array := []Response{}
	for _, log := range selected {
		var res Response
		json.Unmarshal(log.Data, &res)
		array = append(array, res)
	}

	return array
}

func modMsg(msg []byte, user string) []byte {
	res := Response{user, string(msg), time.Now().Format("15:04")}
	json, err := json.Marshal(res)
	if err != nil {
		json = []byte("")
	}
	return json
}
