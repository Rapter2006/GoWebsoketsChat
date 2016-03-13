package main

import (
	"time"
)

type config struct {
	writeWait                 time.Duration
	readWait                  time.Duration
	pingPeriod                time.Duration
	pongWait                  time.Duration
	maxMessageSize            int64
	broadcastMessageQueueSize int64
}

func Initialize(writeWait, readWait, pingPeriod, pongWait time.Duration, maxMessageSize, broadcastMessageSize int64) {
	cfg = &config{writeWait: writeWait,
		readWait:                  readWait,
		pingPeriod:                pingPeriod,
		pongWait:                  pongWait,
		maxMessageSize:            maxMessageSize,
		broadcastMessageQueueSize: broadcastMessageSize,
	}
}