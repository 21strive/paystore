package payment

import (
	"fmt"
	"github.com/21strive/redifu"
	"paystore/helper"
	"testing"
	"time"
)

type Xendit struct {
	*redifu.Record
	Amount        int64  `db:"amount"`
	PaymentMethod string `db:"payment_method"`
}

func (x *Xendit) GetAlias() string {
	return "x"
}

func (x *Xendit) GetTableName() string {
	return "xendit"
}

func (x *Xendit) GetFields() []string {
	return helper.FetchColumns(x)
}

func (x *Xendit) GetScanDestinations() []interface{} {
	return []interface{}{
		&x.UUID,
		&x.RandId,
		&x.CreatedAt,
		&x.UpdatedAt,
		&x.Amount,
		&x.PaymentMethod,
	}
}

func TestQueryJoin(t *testing.T) {
	dummyXendit := &Xendit{}
	redifu.InitRecord(dummyXendit)
	dummyXendit.SetUUID()
	dummyXendit.SetRandId()
	dummyXendit.SetCreatedAt(time.Now())
	dummyXendit.SetUpdatedAt(time.Now())
	dummyXendit.Amount = 10000
	dummyXendit.PaymentMethod = "credit_card"

	repository, _ := NewRepository(nil, nil, dummyXendit, nil, nil)
	finalQuery := repository.JoinBuilder()
	fmt.Print(finalQuery)
}
