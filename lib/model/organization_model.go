package model

import (
	"github.com/21strive/redifu"
)

type Organization struct {
	*redifu.Record
	Name string
	Slug string
}

func (o *Organization) SetName(name string) {
	o.Name = name
}

func (o *Organization) SetSlug(slug string) {
	o.Slug = slug
}

func NewOrganization() *Organization {
	organization := &Organization{}
	redifu.InitRecord(organization)
	return organization
}
