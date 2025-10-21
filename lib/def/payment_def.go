package def

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

type ReceivePaymentRequest struct {
	AccountUUID    string `json:"account_uuid" binding:"required"`
	Amount         int64  `json:"amount" binding:"required"`
	VendorRecordID string `json:"vendor_record_id" binding:"required"`
}
