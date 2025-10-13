package balance

import "github.com/21strive/redifu"

type Repository struct {
	base     *redifu.Base[Account]
	timeline *redifu.Timeline[Account]
}

func (br *Repository) Create(account *Account) (err error) {
	return nil
}

func (br *Repository) Update(account *Account) (err error) {
	return nil
}

func (br *Repository) Delete(account *Account) (err error) {
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
