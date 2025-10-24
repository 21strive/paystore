package payment

import (
	"errors"
)

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusPaid    PaymentStatus = "paid"
	PaymentStatusFailed  PaymentStatus = "failed"
)

var UnmatchBalance = errors.New("The payment owner must match the account balance.")
var VendorRequired = errors.New("Vendor is required.")
var OrganizationRequired = errors.New("Organization is required")
var ConfigRequired = errors.New("Config is required")
var PaymentRequired = errors.New("Payment is required")
var PaymentNotFound = errors.New("Payment not found")
var FinalAmountLessThanZero = errors.New("Final amount must be greater than zero")
