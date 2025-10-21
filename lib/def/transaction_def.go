package def

type TransactionType string

const (
	TypePayment  TransactionType = "payment"
	TypeWithdraw TransactionType = "withdraw"
)
