package payment

import (
	"database/sql"
	"github.com/21strive/redifu"
	"paystore/balance"
)

type Repository struct {
	base           *redifu.Base[Payment]
	timeline       *redifu.Timeline[Payment]
	timelineSeeder *redifu.TimelineSeeder[Payment]
	sorted         *redifu.Sorted[Payment]
	sortedSeeder   *redifu.SortedSeeder[Payment]
}

func (br *Repository) Create(tx *sql.Tx, payment *Payment) (err error) {
	return nil
}

func (br *Repository) FindLatestPayment(balance *balance.Account) (payment *Payment, err error) {
	return nil, nil
}

func (br *Repository) SeedPartial() error {
	return nil
}

func (br *Repository) SeedAll() error {
	return nil
}

func NewManager() *Repository {
	return &Repository{}
}
