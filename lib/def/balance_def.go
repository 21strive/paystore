package def

import "errors"

var InsufficientFunds = errors.New("Insufficient funds")
var BalanceNotFound = errors.New("Account not found")
