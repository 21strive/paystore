package lib

import "github.com/21strive/redifu"

type XenditPayment struct {
	*redifu.SQLItem
	Amount int64
	Hash   string
}
type Payment struct {
	*redifu.SQLItem
	Amount              int64  `json:"amount"`
	Hash                string `json:"hash"`
	BalanceAfterPayment int64  `json:"balanceAfterPayment"`
	IsPaid              bool   `json:"IsPaid"`
	BalanceUUID         string `json:"BalanceUUID"`
	AccountUUID         string `json:"AccountUUID"`
	Xendit              *XenditPayment
}
