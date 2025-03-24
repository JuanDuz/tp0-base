package common

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

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

func ReceiveAck(conn net.Conn) error {
	msg, err := ReceiveString(conn)
	if err != nil {
		return fmt.Errorf("failed to receive ACK: %w", err)
	}

	if msg != "ACK" {
		return fmt.Errorf("unexpected response: expected 'ACK', got '%s'", msg)
	}

	return nil
}
