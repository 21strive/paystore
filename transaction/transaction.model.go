package transaction

import (
	"github.com/21strive/redifu"
	"paystore"
	"paystore/balance"
)

type Transaction struct {
	*redifu.Record
	Type        TransactionType
	RecordUUID  string
	BalanceUUID string
}

func (t *Transaction) SetType(transactionType TransactionType) {
	t.Type = transactionType
}

func (t *Transaction) SetRecord(transaction main.CommonTransaction) {
	t.RecordUUID = transaction.GetUUID()
}

func (t *Transaction) SetBalance(balance balance.Balance) {
	t.BalanceUUID = balance.UUID
}

func NewTransaction() *Transaction {
	transaction := &Transaction{}
	redifu.InitRecord(transaction)
	return transaction
}
