package common

import (
	"context"
	"errors"
	"time"
)

type PollWinnersUseCase interface {
	Execute(ctx context.Context) error
}

type pollWinnersUseCase struct {
	clientFactory func() (*NetworkClient, error)
	clientId      string
	pollInterval  time.Duration
}

func NewPollWinnersUseCase(
	clientFactory func() (*NetworkClient, error),
	clientID string,
	pollInterval time.Duration,
) PollWinnersUseCase {
	return &pollWinnersUseCase{
		clientId:      clientID,
		clientFactory: clientFactory,
		pollInterval:  pollInterval,
	}
}

func (u *pollWinnersUseCase) Execute(ctx context.Context) error {
	for ctx.Err() == nil {

		client, err := u.clientFactory()
		if err != nil {
			log.Errorf("action: connect | result: fail | error: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		winners, err := client.GetWinners(u.clientId)
		client.Close()

		switch {
		case err == nil:
			log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %d", len(winners))
			return nil

		case errors.Is(err, ErrLotteryNotEnded):
			log.Infof("action: consulta_ganadores | result: in_progress")

		case errors.Is(err, ErrInvalidAgency), errors.Is(err, ErrInvalidGetWinners):
			log.Errorf("action: consulta_ganadores | result: fail | reason: %v", err)
			return err

		default:
			log.Errorf("action: consulta_ganadores | result: fail | unexpected error: %v", err)
		}

		time.Sleep(u.pollInterval)
	}
	return nil
}
