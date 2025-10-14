package transaction

import (
	"github.com/21strive/redifu"
	"paystore/balance"
	"paystore/definition"
)

type Transaction struct {
	*redifu.Record
	Type        definition.TransactionType
	RecordUUID  string
	BalanceUUID string
}

func (t *Transaction) SetType(transactionType definition.TransactionType) {
	t.Type = transactionType
}

func (t *Transaction) SetRecord(transaction definition.CommonTransaction) {
	t.RecordUUID = transaction.GetUUID()
}

func (t *Transaction) SetBalance(balance balance.Account) {
	t.BalanceUUID = balance.UUID
}

func NewTransaction() *Transaction {
	transaction := &Transaction{}
	redifu.InitRecord(transaction)
	return transaction
}
