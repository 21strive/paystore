package transaction

import (
	"database/sql"
	"github.com/21strive/redifu"
	"paystore/lib/model"
)

type RepositoryClient interface {
	Create()
	SeedPartial()
}

type Repository struct {
	base                    *redifu.Base[*model.Transaction]
	timelineByBalance       *redifu.Timeline[*model.Transaction]
	timelineSeederByBalance *redifu.TimelineSeeder[*model.Transaction]
}

func (r *Repository) Create(tx *sql.Tx, transaction *model.Transaction) error {
	query := `INSERT INTO transaction (uuid, randid, created_at, updated_at, transaction_type, record_uuid, balance_uuid) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, errExec := tx.Exec(query, transaction.GetUUID(), transaction.GetRandId(), transaction.GetCreatedAt(),
		transaction.GetUpdatedAt(), transaction.Type, transaction.RecordUUID, transaction.BalanceUUID)
	if errExec != nil {
		return errExec
	}

	errSet := r.base.Set(transaction)
	if errSet != nil {
		return errSet
	}

	r.timelineByBalance.AddItem(transaction, []string{transaction.BalanceUUID})
	return nil
}
