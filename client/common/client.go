package common

import (
	"context"
	"errors"
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
	config      ClientConfig
	doneSending bool
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config:      config,
		doneSending: false,
	}
	return client
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop(ctx context.Context) {
	loader, err := NewBetLoader("/data/agency-" + c.config.ID + ".csv")
	if err != nil {
		log.Criticalf("action: open_file | result: fail | error: %v", err)
		return
	}
	defer c.closeLoader(loader)

	for {
		select {
		case <-ctx.Done():
			c.shutdown()
			return
		default:
		}
		if !c.doneSending {
			c.sendNextBatch(loader)
		} else {
			retry := c.tryGetWinners()
			if !retry {
				return
			}
		}

		time.Sleep(c.config.LoopPeriod)
	}
}

func (c *Client) closeLoader(loader *BetLoader) {
	if err := loader.Close(); err != nil {
		log.Errorf("action: close_csv | result: fail | error: %v", err)
	} else {
		log.Infof("action: close_csv | result: success")
	}
}

func (c *Client) shutdown() {
	log.Infof("action: shutdown | result: in_progress | reason: received SIGTERM | SIGINT")
	log.Infof("action: shutdown | result: success")
}

func (c *Client) sendNextBatch(loader *BetLoader) {
	batch, err := loader.NextBatch(c.config.BatchSize, c.config.ID)
	if err != nil {
		if err.Error() == "EOF" {
			c.doneSending = true
			return
		}
		log.Errorf("action: read_batch | result: fail | error: %v", err)
		return
	}

	err = c.withBetClient(func(betClient *BetClient) error {
		return betClient.SendBetBatch(batch)
	})
	if err != nil {
		return
	}
}

func (c *Client) tryGetWinners() bool {
	var retry = true
	err := c.withBetClient(func(betClient *BetClient) error {
		winners, err := betClient.GetWinners(c.config.ID)
		switch {
		case err == nil:
			log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %d", len(winners))
			retry = false
			return nil

		case errors.Is(err, ErrLotteryNotEnded):
			log.Infof("action: consulta_ganadores | result: pending")
			return nil

		case errors.Is(err, ErrInvalidAgency), errors.Is(err, ErrInvalidGetWinners):
			log.Errorf("action: consulta_ganadores | result: fail | reason: %v", err)
			return err

		default:
			log.Errorf("action: consulta_ganadores | result: fail | unexpected error: %v", err)
			return err
		}
	})
	if err != nil {
		log.Errorf("action: consulta_ganadores | result: fail | connection error: %v", err)
		return false
	}
	return retry
}

func (c *Client) withBetClient(callback func(*BetClient) error) error {
	betClient, err := NewBetClient(c.config)
	if err != nil {
		log.Errorf("action: connect | result: fail | error: %v", err)
		return err
	}
	defer betClient.Close()

	return callback(betClient)
}
