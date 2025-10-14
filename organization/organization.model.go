package organization

import (
	"github.com/21strive/redifu"
	"paystore/balance"
)

type Organization struct {
	*redifu.Record
	Name        string
	Avatar      string
	BalanceUUID string
}

func (o *Organization) SetName(name string) {
	o.Name = name
}

func (o *Organization) SetAvatar(avatar string) {
	o.Avatar = avatar
}

func (o *Organization) SetBalance(balance *balance.Account) {
	o.BalanceUUID = balance.UUID
}

func NewOrganization() *Organization {
	organization := &Organization{}
	redifu.InitRecord(organization)
	return organization
}
