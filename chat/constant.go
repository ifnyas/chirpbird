package chat

import "time"

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512

	offlineMsg = "<i>(is offline)</i>"
	onlineMsg  = "<i>(is online)</i>"
)
