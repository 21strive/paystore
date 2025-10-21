package def

import "errors"

var InsufficientFunds = errors.New("Insufficient funds")
var AccountNotFound = errors.New("Account not found")
