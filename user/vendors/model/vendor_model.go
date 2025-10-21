package model

import (
	"github.com/21strive/redifu"
	"paystore/lib/helper"
	"time"
)

type Vendor struct {
	// TODO: Fill attributes of your payment vendor here
	// This model maps webhook data from your payment provider to your database schema
	// Customize these fields to match your specific payment vendor's webhook format

	*redifu.Record
	ID                     string    `json:"id"`          // Xendit Invoice ID
	ExternalID             string    `json:"external_id"` // Aturjadwal Invoice ID
	UserID                 string    `json:"user_id"`     // Aturjadwal User ID
	Status                 string    `json:"status"`
	Amount                 int64     `json:"amount"`
	PaidAmount             int64     `json:"paid_amount"`
	AdjustedReceivedAmount int64     `json:"adjusted_received_amount"`
	FeesPaidAmount         int64     `json:"fees_paid_amount"`
	PaidAt                 time.Time `json:"paid_at"`
	ExpiryDate             time.Time `json:"expiry_date"`
	InvoiceURL             string    `json:"invoice_url"`
	GivenName              string    `json:"given_names,omitempty"`
	Surname                string    `json:"surname,omitempty"`
	Email                  string    `json:"email,omitempty"`
	MobileNumber           string    `json:"mobile_number,omitempty"`
	Created                time.Time `json:"created"`
	Updated                time.Time `json:"updated"`
	Currency               string    `json:"currency"`
	BankCode               string    `json:"bank_code"`
	PaymentMethod          string    `json:"payment_method"`
	PaymentChannel         string    `json:"payment_channel"`
	PaymentDestination     string    `json:"payment_destination"`
}

func (v *Vendor) ScanDestinations() []interface{} {
	// TODO: Fill scan destinations of your Vendor item here
	// This function returns pointers to struct fields for database row scanning
	// The order must exactly match the field order in your Vendor struct definition

	return []interface{}{
		&v.UUID,
		&v.RandId,
		&v.CreatedAt,
		&v.UpdatedAt,
		&v.ID,
		&v.ExternalID,
		&v.UserID,
		&v.Status,
		&v.Amount,
		&v.PaidAmount,
		&v.AdjustedReceivedAmount,
		&v.FeesPaidAmount,
		&v.PaidAt,
		&v.ExpiryDate,
		&v.InvoiceURL,
		&v.GivenName,
		&v.Surname,
		&v.Email,
		&v.MobileNumber,
		&v.Created,
		&v.Updated,
		&v.Currency,
		&v.BankCode,
		&v.PaymentMethod,
		&v.PaymentChannel,
		&v.PaymentDestination,
	}
}

func (v *Vendor) GetFields() []string {
	return helper.FetchColumns(v)
}

func NewVendor() *Vendor {
	vendor := &Vendor{}
	redifu.InitRecord(vendor)
	return vendor
}
