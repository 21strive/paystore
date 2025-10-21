package withdraw

import (
	"github.com/21strive/redifu"
	"paystore/balance"
	"paystore/organization"
	"time"
)

type Withdraw struct {
	*redifu.Record
	Amount               int64          `json:"amount"`
	BalanceBeforePayment int64          `json:"balanceBeforePayment"`
	BalanceAfterPayment  int64          `json:"balanceAfterPayment"`
	BalanceUUID          string         `json:"BalanceUUID"`
	OrganizationUUID     string         `json:"organizationUUID"`
	VendorRecordID       string         `json:"vendorRecordID"`
	Status               WithdrawStatus `json:"status"`
	Hash                 string         `json:"hash"`
}

type HashPayload struct {
	UUID                 string         `json:"uuid"`
	RandId               string         `json:"randid"`
	CreatedAt            time.Time      `json:"createdAt"`
	Amount               int64          `json:"amount"`
	BalanceBeforePayment int64          `json:"balanceBeforePayment"`
	BalanceAfterPayment  int64          `json:"balanceAfterPayment"`
	BalanceUUID          string         `json:"balanceUUID"`
	OrganizationUUID     string         `json:"organizationUUID"`
	VendorRecordID       string         `json:"vendorRecordID"`
	Status               WithdrawStatus `json:"status"`
	PreviousWithdrawHash string         `json:"previousWithdrawHash"`
}

func (w *Withdraw) SetBalance(balance *balance.Balance) {
	w.BalanceUUID = balance.UUID
}

func (w *Withdraw) SetAmount(amount int64) {}

func (w *Withdraw) SetOrganization(organization organization.Organization) {
	w.OrganizationUUID = organization.UUID
}

func (w *Withdraw) SetVendorRecord(uuid string) {
	w.VendorRecordID = uuid
}

func (w *Withdraw) SetSuccess() {
	w.Status = StatusSuccess
}

func (w *Withdraw) SetFailed() {
	w.Status = StatusFailed
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
