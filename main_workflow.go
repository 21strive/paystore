package paystore

import (
	"database/sql"
	"github.com/redis/go-redis/v9"
	"paystore/config"
	"paystore/lib/balance"
	"paystore/lib/def"
	"paystore/lib/model"
	"paystore/lib/organization"
	"paystore/lib/payment"
	vendorRepo "paystore/user/vendorspec/repository"
)

type OrganizationClient struct {
	organizationRepository organization.RepositoryClient
}

func (oc *OrganizationClient) Register(request def.CreateOrganizationRequest) error {
	newOrganization := model.NewOrganization()
	newOrganization.Name = request.Name
	newOrganization.Slug = request.Slug

	organizationFromDB, errFind := oc.organizationRepository.FindBySlug(newOrganization.Slug)
	if errFind != nil {
		return errFind
	}
	if organizationFromDB != nil {
		return nil
	}

	return oc.organizationRepository.Create(newOrganization)
}

type OrganizationFinder struct {
	oc *OrganizationClient
}

func (of *OrganizationFinder) ByUUID(uuid string) (*model.Organization, error) {
	return of.oc.organizationRepository.FindByUUID(uuid)
}

func (of *OrganizationFinder) BySlug(slug string) (*model.Organization, error) {
	return of.oc.organizationRepository.FindBySlug(slug)
}

func (oc *OrganizationClient) Find() *OrganizationFinder {
	return &OrganizationFinder{oc: oc}
}

func NewOrganizationClient(writeDB *sql.DB, readDB *sql.DB,
	redis redis.UniversalClient, config *config.App) *OrganizationClient {

	organizationRepository := organization.NewRepository(writeDB, readDB, redis, config)
	return &OrganizationClient{
		organizationRepository: organizationRepository,
	}
}

type PaystoreClient struct {
	writeDB           *sql.DB
	balanceRepository balance.RepositoryClient
	paymentRepository payment.RepositoryClient
}

func (ps *PaystoreClient) ReceivePayment(request def.ReceivePaymentRequest, selectedOrganization *model.Organization) error {
	balanceFromDB, errFind := ps.balanceRepository.FindByUUID(request.AccountUUID)
	if errFind != nil {
		return errFind
	}
	if balanceFromDB == nil {
		return def.AccountNotFound
	}

	if balanceFromDB.OrganizationUUID != selectedOrganization.GetUUID() {
		return def.OrganizationMismatch
	}

	newPayment := model.NewPayment()
	newPayment.Amount = request.Amount
	newPayment.VendorRecordID = request.VendorRecordID
	newPayment.BalanceUUID = balanceFromDB.GetUUID()
	newPayment.OrganizationUUID = selectedOrganization.GetUUID()

	tx, errInitTx := ps.writeDB.Begin()
	if errInitTx != nil {
		return errInitTx
	}
	defer tx.Rollback()

	errCreatePayment := ps.paymentRepository.Create(tx, newPayment, balanceFromDB)
	if errCreatePayment != nil {
		return errCreatePayment
	}

	balanceFromDB.LastReceive = newPayment.GetCreatedAt()
	balanceFromDB.Collect(newPayment.Amount)
	errUpdateBalance := ps.balanceRepository.Update(tx, balanceFromDB)
	if errUpdateBalance != nil {
		return errUpdateBalance
	}

	errCommit := tx.Commit()
	if errCommit != nil {
		return errCommit
	}

	return nil
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
		return def.AccountNotFound
	}

	return psr.ps.paymentRepository.SeedPartialByBalance(subtraction, lastRandId, balanceFromDB)
}

// TODO: Payment seeder by Organization

func (ps *PaystoreClient) SeedPayment(subtraction int64, lastRandId string, balanceUUID string) error {
	balanceFromDB, errFind := ps.balanceRepository.FindByUUID(balanceUUID)
	if errFind != nil {
		return errFind
	}

	errSeed := ps.paymentRepository.SeedPartialByBalance(subtraction, lastRandId, balanceFromDB)
	if errSeed != nil {
		return errSeed
	}

	return nil
}

func (ps *PaystoreClient) InitWithdraw(amount int64, balance *model.Balance) error {
	errWithdraw := balance.Withdraw(amount)
	if errWithdraw != nil {
		return errWithdraw
	}

	return nil
}

func New(writeDB *sql.DB, readDB *sql.DB, redis redis.UniversalClient,
	config *config.App, vendorConfig *config.Vendor, organization *model.Organization) *PaystoreClient {
	var errInit error

	balanceRepo := balance.NewRepository(readDB, redis, config)
	vendorRepo := vendorRepo.NewVendorRepository()
	paymentRepo, errInit := payment.NewRepository(readDB, redis, vendorConfig, vendorRepo, organization, config)
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
