package payment

import (
	"errors"
	"time"
)

type PaymentStatus string

const (
	StatusPending PaymentStatus = "pending"
	StatusPaid    PaymentStatus = "paid"
	StatusFailed  PaymentStatus = "failed"
)

type Config struct {
	EntityName    string
	ItemPerPage   int64
	RecordAge     time.Duration
	PaginationAge time.Duration
}

var UnmatchBalance = errors.New("The payment owner must match the account balance.")
