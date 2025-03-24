package common

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
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
	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
		select {
		case <-ctx.Done():
			log.Infof("action: loop_interrupted | result: success | client_id: %v", c.config.ID)
			return
		default:
		}

		// c.ensureConnection()

		err := sendAndReceive(c, msgID)

		if err != nil {
			log.Errorf("action: send_receive | result: fail | client_id: %v | error: %v", c.config.ID, err)
			c.conn.Close()
			c.conn = nil
			continue
		}

		time.Sleep(c.config.LoopPeriod)
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

/*func (c *Client) ensureConnection() {
	if c.betClient == nil {
		if betClient, err := NewBetClient(c.config); err != nil {
			log.Errorf("action: reconnect | result: fail | client_id: %v | error: %v", c.config.ID, err)
			return
		}
		c.betClient = betClient
		log.Infof("action: reconnect | result: success | client_id: %v", c.config.ID)
	}
}*/

func sendAndReceive(c *Client, msgID int) error {
	_, err := fmt.Fprintf(
		c.conn,
		"[CLIENT %v] Message N°%v\n",
		c.config.ID,
		msgID,
	)
	if err != nil {
		return fmt.Errorf("send failed: %w", err)
	}

	msg, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		return fmt.Errorf("receive failed: %w", err)
	}

	log.Infof("action: receive_message | result: success | client_id: %v | msg: %v", c.config.ID, msg)
	return nil
}

func (c *Client) Close() {
	c.betClient.Close()
}

func createConnection(serverAddress string, clientID string) (net.Conn, error) {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			clientID,
			err,
		)
		return nil, err
	}
	return conn, nil
}

//
//  Data -- Protocol
//

// SendString sends a message through the given connection using the protocol:
// message_length\nmessage_body
func SendString(conn net.Conn, message string) error {
	messageLength := len(message)
	formatted := fmt.Sprintf("%d\n%s", messageLength, message)

	totalSent := 0
	for totalSent < len(formatted) {
		n, err := conn.Write([]byte(formatted[totalSent:]))
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
		totalSent += n
	}
	return nil
}

// ReceiveString receives a message using the protocol:
// message_length\nmessage_body
func ReceiveString(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	lengthStr, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read length: %w", err)
	}

	lengthStr = strings.TrimSpace(lengthStr)
	messageLength, err := strconv.Atoi(lengthStr)
	if err != nil {
		return "", fmt.Errorf("invalid length: %w", err)
	}

	message := make([]byte, messageLength)
	_, err = io.ReadFull(reader, message)
	if err != nil {
		return "", fmt.Errorf("failed to read message: %w", err)
	}

	return string(message), nil
}

// FormatBetMessage formats a bet as a string separated by '|'.
func FormatBetMessage(nombre, apellido, dni, nacimiento, numero string) string {
	return fmt.Sprintf("%s|%s|%s|%s|%s", nombre, apellido, dni, nacimiento, numero)
}

// ParseBetMessage parses a formatted bet message into individual fields.
func ParseBetMessage(msg string) (nombre, apellido, dni, nacimiento, numero string, err error) {
	fields := strings.Split(msg, "|")
	if len(fields) != 5 {
		return "", "", "", "", "", fmt.Errorf("invalid message format")
	}
	return fields[0], fields[1], fields[2], fields[3], fields[4], nil
}

//
//   Data -- BetClient
//

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

// SendBet formats and sends the bet using the protocol.
func (bc *BetClient) SendBet(nombre, apellido, dni, nacimiento, numero string) error {
	msg := FormatBetMessage(nombre, apellido, dni, nacimiento, numero)
	return SendString(bc.conn, msg)
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
