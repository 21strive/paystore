package withdraw

type WithdrawStatus string

const (
	StatusPending WithdrawStatus = "pending"
	StatusSuccess WithdrawStatus = "success"
	StatusFailed  WithdrawStatus = "failed"
)
