package definition

import "errors"

type CommonTransaction interface {
	GetUUID() string
}
type TransactionType string

const (
	TypePayment  TransactionType = "payment"
	TypeWithdraw TransactionType = "withdraw"
)

var InsufficientFunds = errors.New("Insufficient funds")
var UnmatchBalance = errors.New("The payment owner must match the account balance.")
