package organization

import (
	"github.com/21strive/redifu"
)

type Organization struct {
	*redifu.Record
	Name   string
	Slug   string
	Avatar string
}

func (o *Organization) SetName(name string) {
	o.Name = name
}

func (o *Organization) SetSlug(slug string) {
	o.Slug = slug
}

func (o *Organization) SetAvatar(avatar string) {
	o.Avatar = avatar
}

func NewOrganization() *Organization {
	organization := &Organization{}
	redifu.InitRecord(organization)
	return organization
}
