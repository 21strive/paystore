package config

import (
	"paystore/lib/helper"
	"paystore/user"
	"time"
)

type App struct {
	ItemPerPage      int64
	RecordAge        time.Duration
	PaginationAge    time.Duration
	VendorTableAlias string
	VendorTableName  string
	vendorSampleItem *user.Vendor
}

func (a *App) GetVendorTableAlias() string {
	return a.VendorTableAlias
}

func (a *App) GetVendorTableName() string {
	return a.VendorTableName
}

func (a *App) GetFields() []string {
	return helper.FetchColumns(a.vendorSampleItem)
}

func DefaultConfig(vendorTableName string) *App {
	var vendorTableAlias string
	if len(vendorTableName) > 0 {
		firstChar := string(vendorTableName[0])
		if firstChar == "p" {
			vendorTableAlias = "v"
		} else {
			vendorTableAlias = firstChar
		}
	}

	vendorSampleItem := user.NewVendor()

	return &App{
		ItemPerPage:      50,
		RecordAge:        time.Hour * 12,
		PaginationAge:    time.Hour * 24,
		VendorTableName:  vendorTableName,
		VendorTableAlias: vendorTableAlias,
		vendorSampleItem: vendorSampleItem,
	}
}
