package model

import (
	"github.com/21strive/redifu"
	"paystore/lib/helper"
)

type Vendor struct {
	*redifu.Record
	// TODO: Fill attributes of your payment vendor here
}

func (v *Vendor) ScanDestinations() []interface{} {
	// TODO: Fill scan destionations of your Vendor item here
	return nil
}

func (v *Vendor) GetFields() []string {
	return helper.FetchColumns(v)
}

func NewVendor() *Vendor {
	vendor := &Vendor{}
	redifu.InitRecord(vendor)
	return vendor
}
