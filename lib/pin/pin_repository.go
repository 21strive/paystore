package pin

import (
	"database/sql"
	"github.com/21strive/redifu"
	"github.com/redis/go-redis/v9"
	"paystore/lib/model"
	"time"
)

type Repository struct {
	base *redifu.Base[*model.Pin]
}

func (r *Repository) Create(tx sql.Tx, pin *model.Pin) error {
	query := `INSERT INTO pin (uuid, randid, created_at, updated_at, pin, balance_uuid) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := tx.Exec(query, pin.GetUUID(), pin.GetRandId(), pin.GetCreatedAt(), pin.GetUpdatedAt(), pin.PIN, pin.BalanceUUID)
	if err != nil {
		return err
	}

	errSet := r.base.Set(pin)
	if errSet != nil {
		return errSet
	}

	return nil
}

func (r *Repository) Update(tx sql.Tx, pin *model.Pin) error {
	query := `UPDATE pin SET updated_at = $1, pin = $2 WHERE uuid = $3`
	_, errExec := tx.Exec(query, pin.GetUpdatedAt(), pin.PIN, pin.GetUUID())
	if errExec != nil {
		return errExec
	}

	return nil
}

func NewPINRepository(redis redis.UniversalClient, recordAge time.Duration) *Repository {
	base := redifu.NewBase[*model.Pin](redis, "pin:%s", recordAge)
	return &Repository{
		base: base,
	}
}
