package model

import (
	"github.com/21strive/redifu"
	"paystore/lib/def"
	"time"
)

type Withdraw struct {
	*redifu.Record
	Amount               int64              `json:"amount"`
	BalanceBeforePayment int64              `json:"balanceBeforePayment"`
	BalanceAfterPayment  int64              `json:"balanceAfterPayment"`
	BalanceUUID          string             `json:"BalanceUUID"`
	OrganizationUUID     string             `json:"organizationUUID"`
	VendorRecordID       string             `json:"vendorRecordID"`
	Status               def.WithdrawStatus `json:"status"`
	Hash                 string             `json:"hash"`
}

type WithdrawHashPayload struct {
	UUID                 string             `json:"uuid"`
	RandId               string             `json:"randid"`
	CreatedAt            time.Time          `json:"createdAt"`
	Amount               int64              `json:"amount"`
	BalanceBeforePayment int64              `json:"balanceBeforePayment"`
	BalanceAfterPayment  int64              `json:"balanceAfterPayment"`
	BalanceUUID          string             `json:"balanceUUID"`
	OrganizationUUID     string             `json:"organizationUUID"`
	VendorRecordID       string             `json:"vendorRecordID"`
	Status               def.WithdrawStatus `json:"status"`
	PreviousWithdrawHash string             `json:"previousWithdrawHash"`
}

func (w *Withdraw) SetBalance(balance *Balance) {
	w.BalanceUUID = balance.UUID
}

func (w *Withdraw) SetAmount(amount int64, currentBalanceAmount int64) {
	w.Amount = amount
	w.BalanceBeforePayment = currentBalanceAmount
	w.BalanceAfterPayment = w.BalanceBeforePayment - w.Amount
}

func (w *Withdraw) SetOrganization(organization *Organization) {
	w.OrganizationUUID = organization.UUID
}

func (w *Withdraw) SetVendorRecord(uuid string) {
	w.VendorRecordID = uuid
}

func (w *Withdraw) SetSuccess() {
	w.Status = def.StatusSuccess
}

func (w *Withdraw) SetFailed() {
	w.Status = def.StatusFailed
}

func (w *Withdraw) ScanDestinations() []interface{} {
	return []interface{}{
		&w.UUID,
		&w.RandId,
		&w.CreatedAt,
		&w.UpdatedAt,
		&w.Amount,
		&w.BalanceBeforePayment,
		&w.BalanceAfterPayment,
		&w.BalanceUUID,
		&w.OrganizationUUID,
		&w.VendorRecordID,
		&w.Status,
		&w.Hash,
	}
}

func NewWithdraw() *Withdraw {
	withdraw := &Withdraw{}
	redifu.InitRecord(withdraw)
	withdraw.Status = def.StatusPending
	return withdraw
}
