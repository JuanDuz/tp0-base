package common

import "errors"

var (
	ErrLotteryNotEnded   = errors.New("lottery has not ended yet")
	ErrInvalidAgency     = errors.New("invalid agency")
	ErrInvalidGetWinners = errors.New("invalid get winners request")
	ErrInvalidBatch      = errors.New("invalid batch format")
	ErrEmptyBatch        = errors.New("empty batch received")
	ErrEOF               = errors.New("EOF")
)
