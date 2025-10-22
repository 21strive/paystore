package command

import (
	"database/sql"
	"github.com/redis/go-redis/v9"
	"paystore/config"
	"paystore/lib/balance"
	"paystore/lib/def"
	"paystore/lib/model"
	"paystore/lib/organization"
	"paystore/lib/payment"
)

type OrganizationClient struct {
	organizationRepository organization.RepositoryClient
}

type PaystoreClient struct {
	writeDB                *sql.DB
	balanceRepository      balance.RepositoryClient
	paymentRepository      payment.RepositoryClient
	organizationRepository organization.RepositoryClient
}

func (ps *PaystoreClient) CreateBalance(externalID string, currency string, organizationSlug string) (*model.Balance, error) {
	organizationFromDB, errFind := ps.organizationRepository.FindBySlug(organizationSlug)
	if errFind != nil {
		return nil, errFind
	}

	newBalance := model.NewBalance()
	newBalance.OrganizationUUID = organizationFromDB.GetUUID()
	newBalance.Currency = currency
	newBalance.ExternalID = externalID
	newBalance.Active = true

	errCreate := ps.balanceRepository.Create(newBalance)
	if errCreate != nil {
		return nil, errCreate
	}

	return newBalance, nil
}

func (ps *PaystoreClient) CreatePayment(accountUUID string, amount int64, vendorRecordID string) (*model.Payment, error) {
	balanceFromDB, errFind := ps.balanceRepository.FindByUUID(accountUUID)
	if errFind != nil {
		return nil, errFind
	}

	organizationFromDB, errFind := ps.organizationRepository.FindByUUID(balanceFromDB.OrganizationUUID)
	if errFind != nil {
		return nil, errFind
	}

	newPayment := model.NewPayment()
	newPayment.Amount = amount
	newPayment.VendorRecordID = vendorRecordID
	newPayment.BalanceUUID = balanceFromDB.GetUUID()
	newPayment.OrganizationUUID = organizationFromDB.GetUUID()

	tx, errInitTx := ps.writeDB.Begin()
	if errInitTx != nil {
		return nil, errInitTx
	}
	defer tx.Rollback()

	errCreatePayment := ps.paymentRepository.Create(newPayment, balanceFromDB, organizationFromDB)
	if errCreatePayment != nil {
		return nil, errCreatePayment
	}

	return newPayment, nil
}

func (ps *PaystoreClient) FinalizedPayment(accountUUID string, paymentUUID string, paymentStatus def.PaymentStatus) error {
	balanceFromDB, errFind := ps.balanceRepository.FindByUUID(accountUUID)
	if errFind != nil {
		return errFind
	}

	paymentFromDB, errFind := ps.paymentRepository.FindByUUID(paymentUUID)
	if errFind != nil {
		return errFind
	}

	updateBalance := false
	if paymentStatus == def.PaymentStatusFailed {
		paymentFromDB.SetFailed()
	} else if paymentStatus == def.PaymentStatusPaid {
		paymentFromDB.SetPaid()
		balanceFromDB.LastReceive = paymentFromDB.GetCreatedAt()
		balanceFromDB.Collect(paymentFromDB.Amount)
		updateBalance = true
	}

	tx, errInitTx := ps.writeDB.Begin()
	if errInitTx != nil {
		return errInitTx
	}
	defer tx.Rollback()

	errUpdatePayment := ps.paymentRepository.Update(tx, paymentFromDB)
	if errUpdatePayment != nil {
		return errUpdatePayment
	}

	if updateBalance {
		errUpdateBalance := ps.balanceRepository.Update(tx, balanceFromDB)
		if errUpdateBalance != nil {
			return errUpdateBalance
		}
	}

	errCommit := tx.Commit()
	if errCommit != nil {
		return errCommit
	}

	return nil
}

func (ps *PaystoreClient) CreateWithdraw(amount string) error {

	//errWithdraw := balanceFromDB.Withdraw(request.Amount)
	//if errWithdraw != nil {
	//	return errWithdraw
	//}

	return nil
}

func (ps *PaystoreClient) FinalizedWithdraw(accountUUID string, withdrawUUID string, withdrawStatus def.WithdrawStatus) error {
	balanceFromDB, errFind := ps.balanceRepository.FindByUUID(accountUUID)
	if errFind != nil {
		return errFind
	}
}

type PaymentSeeder struct {
	ps *PaystoreClient
}

func (psr *PaymentSeeder) ByBalance(subtraction int64, lastRandId string, balanceUUID string) error {
	balanceFromDB, errFind := psr.ps.balanceRepository.FindByUUID(balanceUUID)
	if errFind != nil {
		return errFind
	}
	if balanceFromDB == nil {
		return def.BalanceNotFound
	}

	return psr.ps.paymentRepository.SeedPartialByBalance(subtraction, lastRandId, balanceFromDB)
}

// TODO: Payment seeder by Organization

func (ps *PaystoreClient) SeedPayment() *PaymentSeeder {
	return &PaymentSeeder{ps: ps}
}

func New(writeDB *sql.DB, readDB *sql.DB, redis redis.UniversalClient,
	config *config.App) *PaystoreClient {
	var errInit error

	balanceRepo := balance.NewRepository(writeDB, readDB, redis, config)
	paymentRepo, errInit := payment.NewRepository(readDB, redis, config)
	if errInit != nil {
		panic(errInit)
	}
	if errInit != nil {
		panic(errInit)
	}

	return Client(writeDB, balanceRepo, paymentRepo)
}

func Client(writeDB *sql.DB, balanceRepository balance.RepositoryClient,
	paymentRepository payment.RepositoryClient) *PaystoreClient {
	return &PaystoreClient{
		writeDB:           writeDB,
		balanceRepository: balanceRepository,
		paymentRepository: paymentRepository,
	}
}
