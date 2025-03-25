package common

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

const MaxBatchSizeBytes = 8192

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
	var currentSize int
	log.Info("Previous to next batch loop")
	for len(batch) < maxAmount {
		record, err := bl.reader.Read()
		if err != nil {
			if err == io.EOF {
				if len(batch) > 0 {
					return batch, nil // último batch, incompleto pero válido
				}
				return nil, err // no hay más apuestas
			}
			return nil, err // error real
		}
		bet := &Bet{
			agencyId:       "1", // TODO
			firstName:      record[0],
			lastName:       record[1],
			documentNumber: record[2],
			dob:            record[3],
			number:         mustAtoi(record[4]),
		}
		encoded := FormatBetMessage(bet) + "\n"
		log.Info("Loading bet: " + encoded)
		if currentSize+len(encoded) > MaxBatchSizeBytes {
			break
		}
		batch = append(batch, bet)
		currentSize += len(encoded)
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
