package common

import (
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

func (bc *BetClient) SendBetBatch(bets []*Bet) error {
	var sb strings.Builder

	for _, bet := range bets {
		encoded := FormatBetMessage(bet)
		sb.WriteString(encoded)
		sb.WriteString("\n")
	}

	err := SendString(bc.conn, sb.String())
	if err != nil {
		log.Errorf("action: send_batch | result: fail | error: %v", err)
		return err
	}

	err = ReceiveAck(bc.conn)
	if err != nil {
		log.Errorf("action: send_batch | result: fail | error: %v", err)
		return err
	}
	log.Info("action: send_batch | result: success")
	return nil
}

// Close closes the connection.
func (bc *BetClient) Close() {
	log.Infof("action: close_socket | result: in_progress")
	err := bc.conn.Close()
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
