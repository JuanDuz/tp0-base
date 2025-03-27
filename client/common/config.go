package common

import "time"

type Config struct {
	ClientID     int
	ServerAddr   string
	LoopDelay    time.Duration
	MaxPerBatch  int
	DataFilePath string
}
