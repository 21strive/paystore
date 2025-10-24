package organization

import (
	"github.com/21strive/redifu"
)

type Organization struct {
	*redifu.Record
	Name         string
	Slug         string
	FeesConstant int64
	FeesType     FeesType
}

func (o *Organization) SetName(name string) {
	o.Name = name
}

func (o *Organization) SetSlug(slug string) {
	o.Slug = slug
}

func (o *Organization) SetPaymentFees(feesConstant int64, feesType FeesType) {
	o.FeesConstant = feesConstant
	o.FeesType = feesType
}

func NewOrganization() *Organization {
	organization := &Organization{}
	redifu.InitRecord(organization)
	organization.FeesType = Fixed
	organization.FeesConstant = 0
	return organization
}
