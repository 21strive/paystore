package balance

import (
	"github.com/21strive/redifu"
	"paystore/definition"
	"time"
)

type Account struct {
	*redifu.Record
	Balance              int64
	LastIncome           time.Time
	LastWithdraw         time.Time
	IncomeAccumulation   int64
	WithdrawAccumulation int64
	Currency             string
	ExternalID           string
}

func (ac *Account) SetCurrency(currency string) {
	ac.Currency = currency
}

func (ac *Account) SetExternalID(externalID string) {
	ac.ExternalID = externalID
}

func (ac *Account) Collect(amount int64) {
	if amount > 0 {
		ac.Balance += amount
		ac.IncomeAccumulation += amount
	}
}

func (ac *Account) Withdraw(amount int64) error {
	if amount > ac.Balance {
		return definition.InsufficientFunds
	}

	ac.Balance -= amount
	ac.WithdrawAccumulation += amount
	return nil
}

func NewAccount() *Account {
	account := &Account{}
	redifu.InitRecord(account)
	return account
}
