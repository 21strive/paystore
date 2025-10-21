package config

import (
	"paystore/lib/helper"
	"paystore/user/vendorspec/model"
	"time"
)

type App struct {
	ItemPerPage      int64
	RecordAge        time.Duration
	PaginationAge    time.Duration
	VendorTableAlias string
	VendorTableName  string
}

func DefaultConfig(orgName string, orgSlug string) *App {
	return &App{
		ItemPerPage:   50,
		RecordAge:     time.Hour * 12,
		PaginationAge: time.Hour * 24,
	}
}

type Vendor struct {
	VendorTableAlias string
	VendorTableName  string
	vendorSampleItem *model.Vendor
}

func (v *Vendor) GetVendorTableAlias() string {
	return v.VendorTableAlias
}

func (v *Vendor) GetVendorTableName() string {
	return v.VendorTableName
}

func (v *Vendor) GetFields() []string {
	return helper.FetchColumns(v.vendorSampleItem)
}

func NewVendorConfig(vendorTableAlias string, vendorTableName string) *Vendor {
	vendorSampleItem := model.NewVendor()
	return &Vendor{
		VendorTableAlias: vendorTableAlias,
		VendorTableName:  vendorTableName,
		vendorSampleItem: vendorSampleItem,
	}
}
