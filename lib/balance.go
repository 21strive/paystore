package lib

import (
	"github.com/21strive/redifu"
	"time"
)

type Balance struct {
	*redifu.SQLItem
	Balance      int64
	LastIncome   time.Time
	LastWithdraw time.Time
}
