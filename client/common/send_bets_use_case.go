package common

import (
	"context"
	"errors"
	"time"
)

type SendBetsUseCase interface {
	Execute(ctx context.Context) error
}

type sendBetsUseCase struct {
	loader        BetLoader
	clientFactory func() (*NetworkClient, error)
	clientID      string
	maxPerBatch   int
	loopDelay     time.Duration
}

func NewSendBetsUseCase(
	loader BetLoader,
	clientFactory func() (*NetworkClient, error),
	clientID string,
	maxPerBatch int,
	loopDelay time.Duration,
) SendBetsUseCase {
	return &sendBetsUseCase{
		loader:        loader,
		clientFactory: clientFactory,
		clientID:      clientID,
		maxPerBatch:   maxPerBatch,
		loopDelay:     loopDelay,
	}
}

func (s *sendBetsUseCase) Execute(ctx context.Context) error {
	defer func() {
		err := s.loader.Close()
		if err != nil {
			log.Errorf("action: close_csv | result: fail | error: %v", err)
		} else {
			log.Infof("action: close_csv | result: success")
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Infof("action: shutdown | result: in_progress")
			return nil
		default:
		}

		batch, err := s.loader.NextBatch(s.maxPerBatch, s.clientID)
		if errors.Is(err, ErrEOF) {
			log.Infof("action: send_bets | result: done | reason: all bets sent")
			return nil
		}
		if err != nil {
			log.Errorf("action: read_batch | result: fail | error: %v", err)
			continue
		}

		client, err := s.clientFactory()
		if err != nil {
			log.Errorf("action: connect | result: fail | error: %v", err)
			continue
		}

		err = client.SendBatch(batch)
		closeErr := client.Close()
		if closeErr != nil {
			log.Errorf("action: close_socket | result: fail | error: %v", closeErr)
		} else {
			log.Infof("action: close_socket | result: success")
		}
		if err != nil {
			log.Errorf("action: send_batch | result: fail | error: %v", err)
			continue
		}

		logSentBets(batch)

		log.Infof("action: send_batch | result: success | amount: %d", len(batch))
		time.Sleep(s.loopDelay)
	}
}

func logSentBets(bets []*Bet) {
	for _, bet := range bets {
		log.Infof(
			"action: apuesta_enviada | result: success | dni: %s | numero: %d",
			bet.documentNumber,
			bet.number,
		)
	}
}
