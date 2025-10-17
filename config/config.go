package config

import "time"

type App struct {
	ItemPerPage      int64
	RecordAge        time.Duration
	PaginationAge    time.Duration
	OrganizationName string
	OrganizationSlug string
}

func DefaultConfig(orgName string, orgSlug string) *App {
	return &App{
		ItemPerPage:   50,
		RecordAge:     time.Hour * 12,
		PaginationAge: time.Hour * 24,
	}
}

type Vendor interface {
	GetAlias() string
	GetTableName() string
	GetFields() []string
	GetScanDestinations() []interface{}
}
