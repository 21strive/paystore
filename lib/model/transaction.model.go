package model

import (
	"github.com/21strive/redifu"
	"paystore/lib/def"
)

type CommonTransaction interface {
	GetUUID() string
}

type Transaction struct {
	*redifu.Record
	TransactionType def.TransactionType `json:"transcationType"`
	RecordUUID      string              `json:"recordUUID"`
	BalanceUUID     string              `json:"balanceUUID"`
}

func (t *Transaction) SetType(transactionType def.TransactionType) {
	t.TransactionType = transactionType
}

func (t *Transaction) SetRecord(transaction CommonTransaction) {
	t.RecordUUID = transaction.GetUUID()
}

func (t *Transaction) SetBalance(balance *Balance) {
	t.BalanceUUID = balance.UUID
}

func NewTransaction() *Transaction {
	transaction := &Transaction{}
	redifu.InitRecord(transaction)
	return transaction
}
