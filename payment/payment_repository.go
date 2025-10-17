package payment

import (
	"database/sql"
	"github.com/21strive/redifu"
	"github.com/redis/go-redis/v9"
	"paystore/balance"
	"paystore/config"
	"paystore/organization"
)

var firstPartSelectQuery = `SELECT p.uuid, p.randid, p.created_at, p.updated_at, p.amount, p.balance_before_payment, p.balance_after_payment, p.balance_uuid, p.organization_uuid, p.hash,`
var findLatestPaymentQuery = firstPartSelectQuery + `FROM payment p WHERE p.balance_uuid = $1 ORDER BY created_at DESC LIMIT 1;`
var findPaymentByUUIDQuery = firstPartSelectQuery + `FROM payment p WHERE p.uuid = $1;`

type RepositoryClient interface {
	Create(tx *sql.Tx, payment *Payment, balance *balance.Balance) error
	Update(tx *sql.Tx, payment *Payment, balance *balance.Balance) error
	FindLatestPayment(balance *balance.Balance) (*Payment, error)
	FindByUUID(uuid string) (*Payment, error)
	SeedPartialByBalance(subtraction int64, lastRandId string, balance *balance.Balance) error
	SeedAll() error
}

type Repository struct {
	readDB                  *sql.DB
	base                    *redifu.Base[*Payment]
	timelineByAccount       *redifu.Timeline[*Payment]
	timelineByAccountSeeder *redifu.TimelineSeeder[*Payment]
	Vendor                  config.Vendor
	Organization            *organization.Organization
	findLatestPaymentStmt   *sql.Stmt
	findPaymentByUUIDStmt   *sql.Stmt
}

func (br *Repository) Create(tx *sql.Tx, payment *Payment, balance *balance.Balance) error {
	if payment.BalanceUUID != balance.UUID {
		return UnmatchBalance
	}

	query := `
		INSERT INTO payment (
			uuid, randid, created_at, updated_at,
			amount, balance_before_payment, balance_after_payment,
			balance_uuid, organization_uuid, vendor_record_id, status, hash
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := tx.Exec(query,
		payment.GetUUID(),
		payment.GetRandId(),
		payment.GetCreatedAt(),
		payment.GetUpdatedAt(),
		payment.Amount,
		payment.BalanceBeforePayment,
		payment.BalanceAfterPayment,
		payment.BalanceUUID,
		payment.OrganizationUUID,
		payment.VendorRecordID,
		payment.Status,
		payment.Hash,
	)
	if err != nil {
		return err
	}

	errSet := br.base.Set(payment)
	if errSet != nil {
		return errSet
	}

	errSet = br.timelineByAccount.AddItem(payment, []string{br.Organization.GetRandId(), balance.GetRandId()})
	if errSet != nil {
		return errSet
	}

	return err
}

func (br *Repository) Update(tx *sql.Tx, payment *Payment, balance *balance.Balance) error {
	if payment.BalanceUUID != balance.UUID {
		return UnmatchBalance
	}

	query := `UPDATE payment SET updated_at = $1, organization_uuid = $2, vendor_record_id = $3, status = $4, hash = $5 WHERE uuid = $6`
	_, errExec := tx.Exec(query, payment.GetUpdatedAt(), payment.OrganizationUUID, payment.VendorRecordID, payment.Status, payment.Hash, payment.GetUUID())
	if errExec != nil {
		return errExec
	}

	br.base.Set(payment)
	return nil
}

func (br *Repository) JoinBuilder() string {
	var finalQuery string

	finalQuery += firstPartSelectQuery
	finalQuery += ` `

	if br.Vendor != nil {
		fields := br.Vendor.GetFields()
		for _, field := range fields {
			finalQuery += br.Vendor.GetAlias() + "." + field + ", "
		}

		finalQuery += `FROM payment p`
		finalQuery += ` `
		finalQuery += `LEFT JOIN ` + br.Vendor.GetTableName() + ` ` + br.Vendor.GetAlias()
		finalQuery += ` `
	} else {
		finalQuery += `FROM payment p`
	}

	return finalQuery
}

func (br *Repository) FindLatestPayment(balance *balance.Balance) (*Payment, error) {
	payment, err := br.PaymentRowScanner(br.findLatestPaymentStmt.QueryRow(balance.GetUUID()))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return payment, nil
}

func (br *Repository) FindByUUID(uuid string) (*Payment, error) {
	payment, err := br.PaymentRowScanner(br.findPaymentByUUIDStmt.QueryRow(uuid))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return payment, nil
}

func (br *Repository) SeedPartialByBalance(subtraction int64, lastRandId string, balance *balance.Balance) error {
	joinedQuery := br.JoinBuilder()

	rowQuery := joinedQuery + " WHERE p.randid = $1"
	firstPageQuery := joinedQuery + " WHERE p.balance_uuid = $1 ORDER BY created_at DESC"
	nextPageQuery := joinedQuery + " WHERE p.balance_uuid = $1 AND created_at < $2 ORDER BY created_at DESC"

	return br.timelineByAccountSeeder.SeedPartial(rowQuery, firstPageQuery, nextPageQuery, br.PaymentRowScanner, br.PaymentRowsScanner, []interface{}{balance.GetUUID()}, subtraction, lastRandId, []string{balance.GetRandId()})
}

func (br *Repository) SeedAll() error {
	return nil
}

func NewRepository(readDB *sql.DB, redis redis.UniversalClient, vendor config.Vendor, organization *organization.Organization, config *config.App) (*Repository, error) {
	var err error

	if vendor == nil {
		return nil, VendorRequired
	}
	if organization == nil {
		return nil, OrganizationRequired
	}
	if config == nil {
		return nil, ConfigRequired
	}

	basePayment := redifu.NewBase[*Payment](redis, "payment:%s", config.RecordAge)
	timelineByAccount := redifu.NewTimeline[*Payment](redis, basePayment, "payment:organization:%s:account:%s", config.ItemPerPage, redifu.Descending, config.PaginationAge)
	timelineByAccountSeeder := redifu.NewTimelineSeeder[*Payment](readDB, basePayment, timelineByAccount)

	findLatestPaymentStmt, err := readDB.Prepare(findLatestPaymentQuery)
	if err != nil {
		panic(err)
	}
	findPaymentByUUIDStmt, err := readDB.Prepare(findPaymentByUUIDQuery)
	if err != nil {
		panic(err)
	}

	return &Repository{
		base:                    basePayment,
		timelineByAccount:       timelineByAccount,
		timelineByAccountSeeder: timelineByAccountSeeder,
		Vendor:                  vendor,
		Organization:            organization,
		findLatestPaymentStmt:   findLatestPaymentStmt,
		findPaymentByUUIDStmt:   findPaymentByUUIDStmt,
	}, nil
}

func (br *Repository) PaymentRowScanner(row *sql.Row) (*Payment, error) {
	var scanDests []interface{}

	payment := NewPayment()
	scanDests = append(scanDests, payment.ScanDestinations()...)

	if br.Vendor != nil {
		scanDests = append(scanDests, br.Vendor.GetScanDestinations()...)
	}

	err := row.Scan(scanDests...)
	return payment, err
}

func (br *Repository) PaymentRowsScanner(rows *sql.Rows) (*Payment, error) {
	var scanDests []interface{}

	payment := NewPayment()
	scanDests = append(scanDests, payment.ScanDestinations()...)

	if br.Vendor != nil {
		scanDests = append(scanDests, br.Vendor.GetScanDestinations()...)
	}

	err := rows.Scan(scanDests...)
	return payment, err
}
