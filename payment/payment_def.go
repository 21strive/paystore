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
	ItemPerPage   int64
	RecordAge     time.Duration
	PaginationAge time.Duration
}

var UnmatchBalance = errors.New("The payment owner must match the account balance.")
var VendorRequired = errors.New("Vendor is required.")
var OrganizationRequired = errors.New("Organization is required")
var ConfigRequired = errors.New("Config is required")
