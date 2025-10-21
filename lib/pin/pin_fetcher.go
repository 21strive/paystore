package pin

import (
	"github.com/21strive/redifu"
	"github.com/redis/go-redis/v9"
	"paystore/lib/model"
)

type Fetcher struct {
	base *redifu.Base[*model.Pin]
}

func (f *Fetcher) FetchByBalance(balanceUUID string) (*model.Pin, error) {
	pin, errFetch := f.base.Get(balanceUUID)
	if errFetch != nil {
		if errFetch == redis.Nil {
			return nil, nil
		}
	}

	return pin, nil
}

func (f *Fetcher) IsBlankByBalance(balanceUUID string) (bool, error) {
	isBlank, errCheck := f.base.IsBlank(balanceUUID)
	if errCheck != nil {
		return false, errCheck
	}

	return isBlank, nil
}
