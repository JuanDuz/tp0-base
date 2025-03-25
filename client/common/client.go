package common

import (
	"context"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	BatchSize     int
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop(ctx context.Context) {
	filePath := "/data/agency-" + c.config.ID + ".csv"
	loader, err := NewBetLoader(filePath)
	if err != nil {
		log.Criticalf("action: open_file | result: fail | file: %s | error: %v", filePath, err)
		return
	}
	defer func() {
		if err := loader.Close(); err != nil {
			log.Errorf("action: close_csv | result: fail | error: %v", err)
		} else {
			log.Infof("action: close_csv | result: success")
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Infof("action: shutdown | result: in_progress | reason: received SIGTERM | SIGINT")
			log.Infof("action: shutdown | result: success")
			return
		default:
		}
		batch, err := loader.NextBatch(c.config.BatchSize, c.config.ID)
		if err != nil {
			if err.Error() == "EOF" {
				log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
				return
			}
			log.Errorf("action: read_batch | result: fail | error: %v", err)
			return
		}

		betClient, err := NewBetClient(c.config)
		if err != nil {
			log.Errorf("action: connect | result: fail | error: %v", err)
			continue
		}

		err = betClient.SendBetBatch(batch)
		betClient.Close()
		if err != nil {
			continue
		}
		time.Sleep(c.config.LoopPeriod)
	}
}
