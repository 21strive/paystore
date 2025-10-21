package repository

import (
	"database/sql"
	"github.com/21strive/redifu"
	"paystore/user/vendorspec/model"
)

type VendorRepository struct {
	base *redifu.Base[*model.Vendor]
}

func (v *VendorRepository) GetBase() *redifu.Base[*model.Vendor] {
	return v.base
}

func (v *VendorRepository) Create(tx *sql.Tx, vendor *model.Vendor) error {

	return nil
}

func NewVendorRepository() *VendorRepository {
	return &VendorRepository{}
}
