package payment

import (
	"database/sql"
	"github.com/21strive/redifu"
	"github.com/redis/go-redis/v9"
	"paystore/config"
	"paystore/lib/def"
	"paystore/lib/model"
	model2 "paystore/user/vendorspec/model"
	vendorspec "paystore/user/vendorspec/repository"
)

var firstPartSelectQuery = `SELECT p.uuid, p.randid, p.created_at, p.updated_at, p.amount, p.balance_before_payment, p.balance_after_payment, p.balance_uuid, p.organization_uuid, p.hash,`
var findLatestPaymentQuery = firstPartSelectQuery + `FROM payment p WHERE p.balance_uuid = $1 ORDER BY created_at DESC LIMIT 1;`
var findPaymentByUUIDQuery = firstPartSelectQuery + `FROM payment p WHERE p.uuid = $1;`

type RepositoryClient interface {
	Create(tx *sql.Tx, payment *model.Payment, balance *model.Balance) error
	Update(tx *sql.Tx, payment *model.Payment, balance *model.Balance) error
	FindLatestPayment(balance *model.Balance) (*model.Payment, error)
	FindByUUID(uuid string) (*model.Payment, error)
	SeedPartialByBalance(subtraction int64, lastRandId string, balance *model.Balance) error
	SeedAll() error
}

type Repository struct {
	readDB                  *sql.DB
	base                    *redifu.Base[*model.Payment]
	timelineByAccount       *redifu.Timeline[*model.Payment]
	timelineByAccountSeeder *redifu.TimelineSeeder[*model.Payment]
	Vendor                  *config.Vendor
	Organization            *model.Organization
	findLatestPaymentStmt   *sql.Stmt
	findPaymentByUUIDStmt   *sql.Stmt
}

func (br *Repository) Create(tx *sql.Tx, payment *model.Payment, balance *model.Balance) error {
	if payment.BalanceUUID != balance.UUID {
		return def.UnmatchBalance
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

func (br *Repository) Update(tx *sql.Tx, payment *model.Payment, balance *model.Balance) error {
	if payment.BalanceUUID != balance.UUID {
		return def.UnmatchBalance
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
			finalQuery += br.Vendor.GetVendorTableAlias() + "." + field + ", "
		}

		finalQuery += `FROM payment p`
		finalQuery += ` `
		finalQuery += `LEFT JOIN ` + br.Vendor.GetVendorTableName() + ` ` + br.Vendor.GetVendorTableAlias()
		finalQuery += ` `
	} else {
		finalQuery += `FROM payment p`
	}

	return finalQuery
}

func (br *Repository) FindLatestPayment(balance *model.Balance) (*model.Payment, error) {
	payment, err := PaymentRowScanner(br.findLatestPaymentStmt.QueryRow(balance.GetUUID()))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return payment, nil
}

func (br *Repository) FindByUUID(uuid string) (*model.Payment, error) {
	payment, err := PaymentRowScanner(br.findPaymentByUUIDStmt.QueryRow(uuid))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return payment, nil
}

func (br *Repository) SeedPartialByBalance(subtraction int64, lastRandId string, balance *model.Balance) error {
	joinedQuery := br.JoinBuilder()

	rowQuery := joinedQuery + " WHERE p.randid = $1"
	firstPageQuery := joinedQuery + " WHERE p.balance_uuid = $1 ORDER BY created_at DESC"
	nextPageQuery := joinedQuery + " WHERE p.balance_uuid = $1 AND created_at < $2 ORDER BY created_at DESC"

	return br.timelineByAccountSeeder.SeedPartialWithRelation(rowQuery, firstPageQuery, nextPageQuery, PaymentRowScanner, PaymentRowsScanner, []interface{}{balance.GetUUID()}, subtraction, lastRandId, []string{balance.GetRandId()})
}

func (br *Repository) SeedAll() error {
	return nil
}

func NewRepository(readDB *sql.DB, redis redis.UniversalClient, vendorConfig *config.Vendor, vendorRepo *vendorspec.VendorRepository, organization *model.Organization, paystoreConfig *config.App) (*Repository, error) {
	var err error

	if organization == nil {
		return nil, def.OrganizationRequired
	}
	if paystoreConfig == nil {
		return nil, def.ConfigRequired
	}

	vendorRelation := redifu.NewRelation[*model2.Vendor](vendorRepo.GetBase(), "Vendor", "VendorRandId")

	basePayment := redifu.NewBase[*model.Payment](redis, "payment:%s", paystoreConfig.RecordAge)
	timelineByAccount := redifu.NewTimeline[*model.Payment](redis, basePayment, "payment:organization:%s:account:%s", paystoreConfig.ItemPerPage, redifu.Descending, paystoreConfig.PaginationAge)
	timelineByAccount.AddRelation("vendor", vendorRelation)
	timelineByAccountSeeder := redifu.NewTimelineSeeder[*model.Payment](readDB, basePayment, timelineByAccount)

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
		Vendor:                  vendorConfig,
		Organization:            organization,
		findLatestPaymentStmt:   findLatestPaymentStmt,
		findPaymentByUUIDStmt:   findPaymentByUUIDStmt,
	}, nil
}

func PaymentRowScanner(row *sql.Row) (*model.Payment, error) {

	payment := model.NewPayment()

	err := row.Scan(
		&payment.UUID,
		&payment.RandId,
		&payment.CreatedAt,
		&payment.UpdatedAt,
		&payment.Amount,
		&payment.BalanceBeforePayment,
		&payment.BalanceAfterPayment,
		&payment.BalanceUUID,
		&payment.OrganizationUUID,
		&payment.VendorRecordID,
		&payment.Status,
		&payment.Hash,
	)

	return payment, err
}

func PaymentRowsScanner(rows *sql.Rows, relation map[string]redifu.Relation) (*model.Payment, error) {

	payment := model.NewPayment()
	paymentVendor := model2.NewVendor()

	var scanDestionations []interface{}

	scanDestionations = append(scanDestionations, payment.ScanDestinations()...)
	scanDestionations = append(scanDestionations, paymentVendor.ScanDestinations()...)

	err := rows.Scan(scanDestionations)

	if paymentVendor.UUID != "" {
		errSet := relation["vendor"].SetItem(paymentVendor)
		if errSet != nil {
			return nil, errSet
		}
	}

	return payment, err
}
