package common

import (
	"fmt"
	"net"
	"strings"
)

type BetClient struct {
	conn net.Conn
}

// NewBetClient establishes a connection and returns a BetClient instance.
func NewBetClient(config ClientConfig) (*BetClient, error) {
	conn, err := createConnection(config.ServerAddress, config.ID)
	if err != nil {
		return nil, err
	}
	return &BetClient{conn: conn}, nil
}

func (c *BetClient) SendBetBatch(bets []*Bet) error {
	var sb strings.Builder

	for _, bet := range bets {
		encoded := FormatBetMessage(bet)
		sb.WriteString(encoded)
		sb.WriteString("\n")
	}

	err := SendString(c.conn, sb.String())
	if err != nil {
		log.Errorf("action: send_batch | result: fail | error: %v", err)
		return err
	}

	err = ReceiveAck(c.conn)
	if err != nil {
		log.Errorf("action: send_batch | result: fail | error: %v", err)
		return err
	}
	log.Info("action: send_batch | result: success")
	return nil
}

func (c *BetClient) GetWinners(agencyId string) ([]*Bet, error) {
	message := fmt.Sprintf("GET_WINNERS|%s", agencyId)
	if err := SendString(c.conn, message); err != nil {
		return nil, fmt.Errorf("failed to send GET_WINNERS: %w", err)
	}

	bets, err := ReceiveBets(c.conn)
	if err != nil {
		return nil, fmt.Errorf("failed to receive winners: %w", err)
	}

	log.Infof("action: notify_finish | result: success | agency_id: %s", agencyId)
	return bets, nil
}

func ReceiveBets(conn net.Conn) ([]*Bet, error) {
	response, err := ReceiveString(conn)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(response), "\n")
	var bets []*Bet
	for _, line := range lines {
		bet, err := ParseBetMessage(line)
		if err != nil {
			return nil, fmt.Errorf("failed to parse winner bet: %w", err)
		}
		bets = append(bets, bet)
	}
	return bets, nil
}

// Close closes the connection.
func (c *BetClient) Close() {
	log.Infof("action: close_socket | result: in_progress")
	err := c.conn.Close()
	if err != nil {
		log.Infof("action: close_socket | result: fail")
	}
	log.Infof("action: close_socket | result: success")
}

func createConnection(serverAddress string, clientId string) (net.Conn, error) {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			clientId,
			err,
		)
		return nil, err
	}
	return conn, nil
}
