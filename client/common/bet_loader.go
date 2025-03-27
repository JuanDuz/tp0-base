package common

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

const MaxBatchSizeBytes = 8192

type BetLoader interface {
	NextBatch(max int, agencyId int) ([]*Bet, error)
	Close() error
}

type BetCsvLoader struct {
	file     *os.File
	reader   *csv.Reader
	clientID int
}

func NewBetCsvLoader(filePath string, clientID int) (*BetCsvLoader, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return &BetCsvLoader{
		file:     f,
		reader:   csv.NewReader(f),
		clientID: clientID,
	}, nil
}

func (bl *BetCsvLoader) NextBatch(maxAmount int, agencyId int) ([]*Bet, error) {
	var batch []*Bet
	var currentSize int
	for len(batch) < maxAmount {
		record, err := bl.reader.Read()
		if err != nil {
			if err == io.EOF {
				if len(batch) > 0 {
					return batch, nil // último batch, incompleto pero válido
				}
				return nil, ErrEOF // no hay más apuestas
			}
			return nil, err // error real
		}
		bet := &Bet{
			agencyId:       agencyId,
			firstName:      record[0],
			lastName:       record[1],
			documentNumber: record[2],
			dob:            record[3],
			number:         mustAtoi(record[4]),
		}
		encoded := FormatBetMessage(bet) + "\n"
		if currentSize+len(encoded) > MaxBatchSizeBytes {
			break
		}
		batch = append(batch, bet)
		currentSize += len(encoded)
	}
	return batch, nil
}

func (bl *BetCsvLoader) Close() error {
	return bl.file.Close()
}

func mustAtoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
