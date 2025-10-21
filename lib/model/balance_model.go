package model

import (
	"github.com/21strive/redifu"
	"paystore/lib/def"
	"time"
)

type Balance struct {
	*redifu.Record
	Balance              int64
	LastReceive          time.Time
	LastWithdraw         time.Time
	IncomeAccumulation   int64
	WithdrawAccumulation int64
	Currency             string
	Active               bool
	OwnerID              string
	OrganizationUUID     string
}

func (ac *Balance) SetCurrency(currency string) {
	ac.Currency = currency
}

func (ac *Balance) SetOwner(ownerID string) {
	ac.OwnerID = ownerID
}

func (ac *Balance) SetOrganization(organization Organization) {
	ac.OrganizationUUID = organization.GetUUID()
}

func (ac *Balance) Deactivate() {
	ac.Active = false
}

func (ac *Balance) Collect(amount int64) {
	if amount > 0 {
		ac.Balance += amount
		ac.IncomeAccumulation += amount
	}
}

func (ac *Balance) Withdraw(amount int64) error {
	if amount > ac.Balance {
		return def.InsufficientFunds
	}

	ac.Balance -= amount
	ac.WithdrawAccumulation += amount
	return nil
}

func (ac *Balance) ScanDestinations() []interface{} {
	return []interface{}{
		&ac.UUID,
		&ac.RandId,
		&ac.CreatedAt,
		&ac.UpdatedAt,
		&ac.Balance,
		&ac.LastReceive,
		&ac.LastWithdraw,
		&ac.IncomeAccumulation,
		&ac.WithdrawAccumulation,
		&ac.Currency,
		&ac.Active,
		&ac.OwnerID,
		&ac.OrganizationUUID,
	}
}

func NewBalance() *Balance {
	account := &Balance{}
	redifu.InitRecord(account)

	account.Active = true
	return account
}
