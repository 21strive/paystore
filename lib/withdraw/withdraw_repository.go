package withdraw

import (
	"database/sql"
	"github.com/21strive/redifu"
	"github.com/redis/go-redis/v9"
	"paystore/config"
	"paystore/lib/def"
	"paystore/lib/helper"
	"paystore/lib/model"
	vendorModel "paystore/user"
)

var firstPartSelectQuery = `SELECT w.uuid, w.randid, w.created_at, w.updated_at, w.amount, w.balance_before_payment, w.balance_after_payment, w.balance_uuid, w.organization_uuid, w.vendor_record_id, w.status, w.hash`
var findWithdrawByUUIDQuery = firstPartSelectQuery + `FROM withdraw w WHERE w.uuid = $1;`

type RepositoryClient interface {
	Create(tx *sql.Tx, withdraw *model.Withdraw, balance *model.Balance, organization *model.Organization) error
	Update(tx *sql.Tx, withdraw *model.Withdraw) error
	FindByUUID(uuid string) (*model.Withdraw, error)
	SeedPartialByBalance(subtraction int64, lastRandId string, balance *model.Balance) error
}

type Repository struct {
	base                    *redifu.Base[*model.Withdraw]
	timelineByBalance       *redifu.Timeline[*model.Withdraw]
	timelineSeederByBalance *redifu.TimelineSeeder[*model.Withdraw]
	findWithdrawByUUIDStmt  *sql.Stmt
	AppConfig               *config.App
}

func (r *Repository) Close() {
	r.findWithdrawByUUIDStmt.Close()
}

func (r *Repository) Create(tx *sql.Tx, withdraw *model.Withdraw, balance *model.Balance,
	organization *model.Organization) error {
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

func (r *Repository) Update(tx *sql.Tx, withdraw *model.Withdraw) error {
	query := `UPDATE withdraw SET updated_at = $1, status = $2, hash = $3 WHERE uuid = $4`
	_, errExec := tx.Exec(query, withdraw.GetUpdatedAt(), withdraw.Status, withdraw.Hash, withdraw.GetUUID())
	if errExec != nil {
		return errExec
	}

	return r.base.Set(withdraw)
}

func (r *Repository) FindByUUID(uuid string) (*model.Withdraw, error) {
	withdraw, err := WithdrawRowScanner(r.findWithdrawByUUIDStmt.QueryRow(uuid))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, def.WithdrawNotFound
		}
		return nil, err
	}

	return withdraw, nil
}

func (r *Repository) SeedPartialByBalance(subtraction int64, lastRandId string, balance *model.Balance) error {
	joinedQuery := helper.JoinBuilder(firstPartSelectQuery, def.TypeWithdraw, r.AppConfig)

	rowQuery := joinedQuery + " WHERE w.randid = $1"
	firstPageQuery := joinedQuery + " WHERE w.balance_uuid = $1 ORDER BY created_at DESC"
	nextPageQuery := joinedQuery + " WHERE w.balance_uuid = $1 AND created_at < $2 ORDER BY created_at DESC"

	return r.timelineSeederByBalance.SeedPartialWithRelation(
		rowQuery, firstPageQuery, nextPageQuery, WithdrawRowScanner, WithdrawRowsScanner,
		[]interface{}{balance.GetUUID()}, subtraction, lastRandId, []string{balance.GetRandId()})
}

func WithdrawRowScanner(row *sql.Row) (*model.Withdraw, error) {
	withdraw := model.NewWithdraw()
	err := row.Scan(withdraw.ScanDestinations()...)
	return withdraw, err
}

func WithdrawRowsScanner(rows *sql.Rows, relation map[string]redifu.Relation) (*model.Withdraw, error) {
	withdraw := model.NewWithdraw()
	withdrawVendor := vendorModel.NewWithdrawVendor()

	var scanDestinations []interface{}
	scanDestinations = append(scanDestinations, withdraw.ScanDestinations()...)
	scanDestinations = append(scanDestinations, withdrawVendor.ScanDestionations()...)

	err := rows.Scan(scanDestinations...)
	if err != nil {
		return nil, err
	}

	if withdrawVendor.UUID != "" {
		errSet := relation["vendor"].SetItem(withdrawVendor)
		if errSet != nil {
			return nil, errSet
		}
	}

	return withdraw, nil
}

func NewRepository(readDB *sql.DB, redis redis.UniversalClient, config *config.App) *Repository {
	base := redifu.NewBase[*model.Withdraw](redis, "withdraw:%s", config.RecordAge)
	timelineByBalance := redifu.NewTimeline[*model.Withdraw](redis, base,
		"withdraw:organization:%s:balance:%s", config.ItemPerPage, redifu.Descending, config.PaginationAge)
	timelineSeederByBalance := redifu.NewTimelineSeeder[*model.Withdraw](nil, base, timelineByBalance)

	findWithdrawByUUIDStmt, err := readDB.Prepare(findWithdrawByUUIDQuery)
	if err != nil {
		panic(err)
	}

	return &Repository{
		base:                    base,
		timelineByBalance:       timelineByBalance,
		timelineSeederByBalance: timelineSeederByBalance,
		findWithdrawByUUIDStmt:  findWithdrawByUUIDStmt,
		AppConfig:               config,
	}
}
