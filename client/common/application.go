package common

import (
	"context"
	"fmt"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

type Application struct {
	cfg                Config
	betLoader          BetLoader
	sendBetsUseCase    SendBetsUseCase
	pollWinnersUseCase PollWinnersUseCase
}

func NewApplication(cfg Config) (*Application, error) {
	loader, err := NewBetCsvLoader(cfg.DataFilePath, cfg.ClientID)
	if err != nil {
		return nil, fmt.Errorf("failed to create bet loader: %w", err)
	}

	clientFactory := func() (*NetworkClient, error) {
		return NewNetworkClient(cfg.ServerAddr)
	}

	sendBets := NewSendBetsUseCase(loader, clientFactory, cfg.ClientID, cfg.MaxPerBatch, cfg.LoopDelay)
	pollWinners := NewPollWinnersUseCase(clientFactory, cfg.ClientID, cfg.LoopDelay)

	return &Application{
		cfg:                cfg,
		betLoader:          loader,
		sendBetsUseCase:    sendBets,
		pollWinnersUseCase: pollWinners,
	}, nil
}

func (a *Application) Run(ctx context.Context) error {
	if err := a.sendBetsUseCase.Execute(ctx); err != nil {
		return err
	}
	if err := a.pollWinnersUseCase.Execute(ctx); err != nil {
		return err
	}
	return nil
}

func (a *Application) Close() {
	if err := a.betLoader.Close(); err != nil {
		log.Errorf("action: close_csv | result: fail | error: %v", err)
	} else {
		log.Infof("action: close_csv | result: success")
	}
}
