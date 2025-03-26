package common

import (
	"fmt"
	"net"
	"strings"
)

// NetworkClient handles TCP communication with the bet server
// It implements the domain.BetSender and domain.WinnerFetcher interfaces
type NetworkClient struct {
	conn net.Conn
}

func NewNetworkClient(serverAddr string) (*NetworkClient, error) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}
	return &NetworkClient{conn: conn}, nil
}

func (c *NetworkClient) SendBatch(bets []*Bet) error {
	var sb strings.Builder
	for _, b := range bets {
		sb.WriteString(FormatBetMessage(b))
		sb.WriteString("\n")
	}

	if err := SendString(c.conn, sb.String()); err != nil {
		return fmt.Errorf("send failed: %w", err)
	}

	if err := ReceiveAck(c.conn); err != nil {
		return fmt.Errorf("ack failed: %w", err)
	}

	return nil
}

func (c *NetworkClient) GetWinners(agencyID string) ([]*Bet, error) {
	msg := fmt.Sprintf("GET_WINNERS|%s", agencyID)
	if err := SendString(c.conn, msg); err != nil {
		return nil, fmt.Errorf("failed to send GET_WINNERS: %w", err)
	}

	raw, err := ReceiveResponse(c.conn)
	if err != nil {
		return nil, err
	}

	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []*Bet{}, nil
	}

	lines := strings.Split(raw, "\n")
	winners := make([]*Bet, 0, len(lines))
	for _, line := range lines {
		bet, err := ParseBetMessage(line)
		if err != nil {
			return nil, fmt.Errorf("parse error: %w", err)
		}
		winners = append(winners, bet)
	}

	return winners, nil
}

func (c *NetworkClient) Close() {
	err := c.conn.Close()
	if err != nil {
		log.Errorf("action: close_socket | result: fail | error: %v", err)
	} else {
		log.Infof("action: close_socket | result: success")
	}
}
