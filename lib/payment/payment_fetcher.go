package payment

import (
	"github.com/21strive/redifu"
	"github.com/redis/go-redis/v9"
	"paystore/config"
	"paystore/lib/model"
)

type FetcherClient interface {
	FetchByBalance(lastRandId []string, balanceRandId string) ([]*model.Payment, string, string, error)
	IsBlankByBalance(balanceRandId string) (bool, error)
	GetItemPerPage() int64
}

type Fetcher struct {
	base              *redifu.Base[*model.Payment]
	timelineByBalance *redifu.Timeline[*model.Payment]
}

func (f *Fetcher) FetchByBalance(lastRandId []string, balanceRandId string) ([]*model.Payment, string, string, error) {
	return f.timelineByBalance.Fetch([]string{balanceRandId}, lastRandId, nil, nil)
}

func (f *Fetcher) IsBlankByBalance(balanceRandId string) (bool, error) {
	isBlank, errCheck := f.timelineByBalance.IsBlankPage([]string{balanceRandId})
	if errCheck != nil {
		return false, errCheck
	}

	return isBlank, nil
}

func (f *Fetcher) GetItemPerPage() int64 {
	return f.timelineByBalance.GetItemPerPage()
}

func NewFetcher(redis redis.UniversalClient, config *config.App) *Fetcher {
	base := redifu.NewBase[*model.Payment](redis, "payment:%s", config.RecordAge)
	timelineByBalance := redifu.NewTimeline[*model.Payment](redis, base, "payment-balance:%s", config.ItemPerPage, redifu.Descending, config.PaginationAge)
	return &Fetcher{
		base:              base,
		timelineByBalance: timelineByBalance,
	}
}
