package common

import (
	"context"
	"os"
	"strconv"
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
}

// Client Entity that encapsulates how
type Client struct {
	config    ClientConfig
	betClient *BetClient
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	betClient, err := NewBetClient(config)
	if err != nil {
		return nil
	}
	client := &Client{
		config:    config,
		betClient: betClient,
	}
	return client
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop(ctx context.Context) {

	betNumber, err := strconv.Atoi(os.Getenv("NUMERO"))
	if err != nil {
		log.Errorf("action: parse_bet_number | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}
	bet := &Bet{
		firstName:      os.Getenv("NOMBRE"),
		lastName:       os.Getenv("APELLIDO"),
		documentNumber: os.Getenv("DOCUMENTO"),
		dob:            os.Getenv("NACIMIENTO"),
		number:         betNumber,
		agencyId:       c.config.ID,
	}

	err = c.betClient.SendBet(bet)
	if err != nil {
		log.Errorf("action: apuesta_enviada | result: fail | dni: %v | numero: %v",
			bet.documentNumber,
			betNumber,
		)
		return
	}
	log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
		bet.documentNumber,
		betNumber,
	)

}

func (c *Client) Stop() {
	c.betClient.Close()
}
