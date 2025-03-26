package common

import "time"

type Config struct {
	ClientID     string
	ServerAddr   string
	LoopDelay    time.Duration
	MaxPerBatch  int
	DataFilePath string
}
