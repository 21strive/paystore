package request

import "time"

type XenditReceivePayment struct {
	ID                     string    `json:"id"`                       // Xendit Invoice ID
	ExternalID             string    `json:"external_id"`              // Your Invoice ID
	UserID                 string    `json:"user_id"`                  // Xendit User ID
	IsHigh                 bool      `json:"is_high"`                  // High value transaction flag
	PaymentMethod          string    `json:"payment_method"`           // BANK_TRANSFER, EWALLET, etc.
	Status                 string    `json:"status"`                   // PAID, PENDING, EXPIRED
	MerchantName           string    `json:"merchant_name"`            // Your merchant name
	Amount                 int64     `json:"amount"`                   // Invoice amount
	PaidAmount             int64     `json:"paid_amount"`              // Actually paid amount
	BankCode               string    `json:"bank_code"`                // BCA, BNI, MANDIRI, etc.
	PaidAt                 time.Time `json:"paid_at"`                  // Payment timestamp
	PayerEmail             string    `json:"payer_email"`              // Customer email (can be empty)
	Description            string    `json:"description"`              // Invoice description
	AdjustedReceivedAmount int64     `json:"adjusted_received_amount"` // Adjusted amount
	FeesPaidAmount         int64     `json:"fees_paid_amount"`         // Transaction fees
	Updated                time.Time `json:"updated"`                  // Last update timestamp
	Created                time.Time `json:"created"`                  // Creation timestamp
	Currency               string    `json:"currency"`                 // IDR, USD, etc.
	PaymentChannel         string    `json:"payment_channel"`          // Specific channel (BCA, DANA, etc.)
	PaymentDestination     string    `json:"payment_destination"`      // Virtual account number or destination
}
