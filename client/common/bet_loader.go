package common

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type BetLoader struct {
	reader *csv.Reader
	file   *os.File
}

func NewBetLoader(filePath string) (*BetLoader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return &BetLoader{
		reader: csv.NewReader(file),
		file:   file,
	}, nil
}

func (bl *BetLoader) NextBatch(maxAmount int) ([]*Bet, error) {
	var batch []*Bet
	for len(batch) < maxAmount {
		record, err := bl.reader.Read()
		if err != nil {
			if len(batch) == 0 {
				return nil, err // No data read, return EOF or error
			}
			return batch, nil // Return what we could read
		}
		bet := &Bet{
			agencyId:       "1", // TODO
			firstName:      record[0],
			lastName:       record[1],
			documentNumber: record[2],
			dob:            record[3],
			number:         mustAtoi(record[4]),
		}
		batch = append(batch, bet)
	}
	return batch, nil
}

func (bl *BetLoader) Close() error {
	return bl.file.Close()
}

func mustAtoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
