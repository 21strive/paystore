package def

import "errors"

type WithdrawStatus string

const (
	StatusPending WithdrawStatus = "pending"
	StatusSuccess WithdrawStatus = "success"
	StatusFailed  WithdrawStatus = "failed"
)

var WithdrawNotFound = errors.New("Withdraw not found")
