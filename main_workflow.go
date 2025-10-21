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
	"paystore/lib/request"
	vendorModel "paystore/user/vendors/model"
	vendorRepo "paystore/user/vendors/repository"
)

type OrganizationClient struct {
	organizationRepository organization.RepositoryClient
}

//func (oc *OrganizationClient) Register(request def.CreateOrganizationRequest) error {
//	newOrganization := model.NewOrganization()
//	newOrganization.Name = request.Name
//	newOrganization.Slug = request.Slug
//
//	organizationFromDB, errFind := oc.organizationRepository.FindBySlug(newOrganization.Slug)
//	if errFind != nil {
//		return errFind
//	}
//	if organizationFromDB != nil {
//		return nil
//	}
//
//	return oc.organizationRepository.Create(newOrganization)
//}
//
//type OrganizationFinder struct {
//	oc *OrganizationClient
//}
//
//func (of *OrganizationFinder) ByUUID(uuid string) (*model.Organization, error) {
//	return of.oc.organizationRepository.FindByUUID(uuid)
//}
//
//func (of *OrganizationFinder) BySlug(slug string) (*model.Organization, error) {
//	return of.oc.organizationRepository.FindBySlug(slug)
//}
//
//func (oc *OrganizationClient) Find() *OrganizationFinder {
//	return &OrganizationFinder{oc: oc}
//}
//
//func NewOrganizationClient(writeDB *sql.DB, readDB *sql.DB,
//	redis redis.UniversalClient, config *config.App) *OrganizationClient {
//
//	organizationRepository := organization.NewRepository(writeDB, readDB, redis, config)
//	return &OrganizationClient{
//		organizationRepository: organizationRepository,
//	}
//}

type PaystoreClient struct {
	writeDB                *sql.DB
	balanceRepository      balance.RepositoryClient
	paymentRepository      payment.RepositoryClient
	vendorRepository       vendorRepo.VendorRepositoryClient
	organizationRepository organization.RepositoryClient
}

func (ps *PaystoreClient) ReceivePayment(request request.ReceivePaymentRequest, vendorItem *vendorModel.Vendor) (*model.Payment, error) {
	balanceFromDB, errFind := ps.balanceRepository.FindByUUID(request.AccountUUID)
	if errFind != nil {
		return nil, errFind
	}

	organizationFromDB, errFind := ps.organizationRepository.FindByUUID(balanceFromDB.OrganizationUUID)
	if errFind != nil {
		return nil, errFind
	}

	newPayment := model.NewPayment()
	newPayment.Amount = request.Amount
	newPayment.VendorRecordID = request.VendorRecordID
	newPayment.BalanceUUID = balanceFromDB.GetUUID()
	newPayment.OrganizationUUID = organizationFromDB.GetUUID()

	tx, errInitTx := ps.writeDB.Begin()
	if errInitTx != nil {
		return nil, errInitTx
	}
	defer tx.Rollback()

	errCreatePayment := ps.paymentRepository.Create(tx, newPayment, balanceFromDB, organizationFromDB)
	if errCreatePayment != nil {
		return nil, errCreatePayment
	}

	errCreateVendor := ps.vendorRepository.Create(tx, vendorItem)
	if errCreateVendor != nil {
		return nil, errCreateVendor
	}

	balanceFromDB.LastReceive = newPayment.GetCreatedAt()
	balanceFromDB.Collect(newPayment.Amount)
	errUpdateBalance := ps.balanceRepository.Update(tx, balanceFromDB)
	if errUpdateBalance != nil {
		return nil, errUpdateBalance
	}

	errCommit := tx.Commit()
	if errCommit != nil {
		return nil, errCommit
	}

	return newPayment, nil
}

func (ps *PaystoreClient) CreateBalance(request request.CreateBalanceRequest) (*model.Balance, error) {
	organizationFromDB, errFind := ps.organizationRepository.FindBySlug(request.OrganizationSlug)
	if errFind != nil {
		return nil, errFind
	}

	newBalance := model.NewBalance()
	newBalance.OrganizationUUID = organizationFromDB.GetUUID()
	newBalance.Currency = request.Currency
	newBalance.Active = true

	errCreate := ps.balanceRepository.Create(newBalance)
	if errCreate != nil {
		return nil, errCreate
	}

	return newBalance, nil
}

func (ps *PaystoreClient) CreateOrganization(request request.CreateOrganizationRequest) (*model.Organization, error) {
	organizationFromDB, errFind := ps.organizationRepository.FindBySlug(request.Slug)
	if errFind != nil {
		return nil, errFind
	}
	if organizationFromDB != nil {
		return nil, def.DuplicateSlug
	}

	organizationFromDB, errFind = ps.organizationRepository.FindByName(request.Name)
	if errFind != nil {
		return nil, errFind
	}
	if organizationFromDB != nil {
		return nil, def.DuplicateName
	}

	newOrganization := model.NewOrganization()
	newOrganization.Name = request.Name
	newOrganization.Slug = request.Slug

	errCreate := ps.organizationRepository.Create(newOrganization)
	if errCreate != nil {
		return nil, errCreate
	}

	return newOrganization, nil
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
	config *config.App) *PaystoreClient {
	var errInit error

	balanceRepo := balance.NewRepository(writeDB, readDB, redis, config)
	vendorRepo := vendorRepo.NewVendorRepository()
	paymentRepo, errInit := payment.NewRepository(readDB, redis, vendorRepo, config)
	if errInit != nil {
		panic(errInit)
	}
	if errInit != nil {
		panic(errInit)
	}

	return Client(writeDB, balanceRepo, paymentRepo, vendorRepo)
}

func Client(writeDB *sql.DB, balanceRepository balance.RepositoryClient,
	paymentRepository payment.RepositoryClient, vendorRepository vendorRepo.VendorRepositoryClient) *PaystoreClient {
	return &PaystoreClient{
		writeDB:           writeDB,
		balanceRepository: balanceRepository,
		paymentRepository: paymentRepository,
		vendorRepository:  vendorRepository,
	}
}
