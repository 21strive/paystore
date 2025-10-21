package model

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/21strive/redifu"
	"paystore/lib/def"
	"paystore/user/vendors/model"
	"time"
)

type Payment struct {
	*redifu.Record
	Amount               int64             `json:"amount"`
	BalanceBeforePayment int64             `json:"balanceBeforePayment"`
	BalanceAfterPayment  int64             `json:"balanceAfterPayment"`
	BalanceUUID          string            `json:"BalanceUUID"`
	OrganizationUUID     string            `json:"organizationUUID"`
	VendorRecordID       string            `json:"vendorRecordID"`
	Status               def.PaymentStatus `json:"status"`
	Hash                 string            `json:"hash"`
	VendorRandId         string            `json:"vendorRandId,omitempty"`
	Vendor               model.Vendor      `json:"vendor,omitempty"`
}

type PaymentHashPayload struct {
	UUID                 string            `json:"uuid"`
	RandId               string            `json:"randid"`
	CreatedAt            time.Time         `json:"createdAt"`
	Amount               int64             `json:"amount"`
	BalanceBeforePayment int64             `json:"balanceBeforePayment"`
	BalanceAfterPayment  int64             `json:"balanceAfterPayment"`
	BalanceUUID          string            `json:"balanceUUID"`
	OrganizationUUID     string            `json:"organizationUUID"`
	VendorRecordID       string            `json:"vendorRecordID"`
	Status               def.PaymentStatus `json:"status"`
	PreviousPaymentHash  string            `json:"previousPaymentHash"`
}

func (p *Payment) SetBalance(balance *Balance) {
	p.BalanceUUID = balance.UUID
}

func (p *Payment) SetAmount(amount int64, currentBalanceAmount int64) {
	p.Amount = amount
	p.BalanceBeforePayment = currentBalanceAmount
	p.BalanceAfterPayment = p.BalanceBeforePayment + p.Amount
}

func (p *Payment) SetOrganization(organization Organization) {
	p.OrganizationUUID = organization.UUID
}

func (p *Payment) SetVendorRecord(uuid string) {
	p.VendorRecordID = uuid
}

func (p *Payment) GenerateHash(previousPayment *Payment) error {
	hashPayload := PaymentHashPayload{
		UUID:                 p.UUID,
		RandId:               p.RandId,
		CreatedAt:            p.CreatedAt,
		Amount:               p.Amount,
		BalanceBeforePayment: p.BalanceBeforePayment,
		BalanceAfterPayment:  p.BalanceAfterPayment,
		BalanceUUID:          p.BalanceUUID,
		OrganizationUUID:     p.OrganizationUUID,
		VendorRecordID:       p.VendorRecordID,
		Status:               p.Status,
	}

	if previousPayment != nil {
		hashPayload.PreviousPaymentHash = previousPayment.Hash
	}

	paymentHash, errHash := createHash(hashPayload)
	if errHash != nil {
		return errHash
	}

	p.Hash = paymentHash
	return nil
}

func (p *Payment) Verify(previousPayment *Payment) (bool, error) {
	hashPayload := PaymentHashPayload{
		UUID:                 p.UUID,
		RandId:               p.RandId,
		CreatedAt:            p.CreatedAt,
		Amount:               p.Amount,
		BalanceBeforePayment: p.BalanceBeforePayment,
		BalanceAfterPayment:  p.BalanceAfterPayment,
		BalanceUUID:          p.BalanceUUID,
	}

	if previousPayment != nil {
		hashPayload.PreviousPaymentHash = previousPayment.Hash
	}

	currentHash, errHash := createHash(hashPayload)
	if errHash != nil {
		return false, errHash
	}

	return currentHash == p.Hash, nil
}

func (p *Payment) ScanDestinations() []interface{} {
	return []interface{}{
		&p.UUID,
		&p.RandId,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.Amount,
		&p.BalanceBeforePayment,
		&p.BalanceAfterPayment,
		&p.BalanceUUID,
		&p.OrganizationUUID,
		&p.VendorRecordID,
		&p.Status,
		&p.Hash,
	}
}

func (p *Payment) SetPaid() {
	p.Status = def.PaymentStatusPaid
}

func (p *Payment) SetFailed() {
	p.Status = def.PaymentStatusFailed
}

func NewPayment() *Payment {
	payment := &Payment{}
	redifu.InitRecord(payment)
	payment.Status = def.PaymentStatusPending
	return payment
}

func createHash(payload PaymentHashPayload) (string, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(jsonData)
	return hex.EncodeToString(hash[:]), nil
}
