package common

import (
	"fmt"
	"strconv"
	"strings"
)

// FormatBetMessage formats a bet as a string separated by '|'.
func FormatBetMessage(bet *Bet) string {
	return fmt.Sprintf("%s|%s|%s|%s|%d|%s",
		bet.firstName,
		bet.lastName,
		bet.documentNumber,
		bet.dob,
		bet.number,
		bet.agencyId,
	)
}

func ParseBetMessage(msg string) (*Bet, error) {
	fields := strings.Split(msg, "|")
	if len(fields) != 6 {
		return nil, fmt.Errorf("invalid bet message format")
	}

	number, err := strconv.Atoi(fields[4])
	if err != nil {
		return nil, fmt.Errorf("invalid number field: %w", err)
	}

	return &Bet{
		firstName:      fields[0],
		lastName:       fields[1],
		documentNumber: fields[2],
		dob:            fields[3],
		number:         number,
		agencyId:       fields[5],
	}, nil
}
