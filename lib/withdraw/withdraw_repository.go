package withdraw

import (
	"database/sql"
	"github.com/21strive/redifu"
	"paystore/lib/model"
)

type RepositoryClient interface {
	Create(tx *sql.Tx, withdraw *model.Withdraw) error
}

type Repository struct {
	base                    *redifu.Base[*model.Withdraw]
	timelineByBalance       *redifu.Timeline[*model.Withdraw]
	timelineSeederByBalance *redifu.TimelineSeeder[*model.Withdraw]
}

func (r *Repository) Create(tx *sql.Tx, withdraw *model.Withdraw, balance *model.Balance, organization *model.Organization) error {
	query := `
		INSERT INTO withdraw (uuid, randid, created_at, updated_at, amount, balance_before_payment, 
		balance_after_payment, balance_uuid, organization_uuid, vendor_record_id, status, hash) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := tx.Exec(query, withdraw.GetUUID(), withdraw.GetRandId(), withdraw.GetCreatedAt(),
		withdraw.GetUpdatedAt(), withdraw.Amount, withdraw.BalanceBeforePayment, withdraw.BalanceAfterPayment,
		withdraw.BalanceUUID, withdraw.OrganizationUUID, withdraw.VendorRecordID, withdraw.Status, withdraw.Hash)
	if err != nil {
		return err
	}

	errSet := r.base.Set(withdraw)
	if errSet != nil {
		return errSet
	}

	r.timelineByBalance.AddItem(withdraw, []string{organization.GetRandId(), balance.GetRandId()})
	return nil
}
