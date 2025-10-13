package balance

import (
	"database/sql"
	"github.com/21strive/redifu"
)

type Repository struct {
	base           *redifu.Base[Account]
	timeline       *redifu.Timeline[Account]
	timelineSeeder *redifu.TimelineSeeder[Account]
}

func (br *Repository) Create(account *Account) (err error) {
	return nil
}

func (br *Repository) Update(tx *sql.Tx, account *Account) (err error) {
	return nil
}

func (br *Repository) Delete(tx *sql.Tx, account *Account) (err error) {
	return nil
}

func (br *Repository) FindByUUID(uuid string) (account *Account, err error) {
	return nil, nil
}

func (br *Repository) FindByExternalID(externalID string) (account *Account, err error) {
	return nil, nil
}

func (br *Repository) SeedPartial() error {
	return nil
}

func NewManager() *Repository {
	return &Repository{}
}
