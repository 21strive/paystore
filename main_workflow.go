package main

import (
	"database/sql"
	"github.com/redis/go-redis/v9"
	"paystore/balance"
	"paystore/config"
	"paystore/organization"
	"paystore/payment"
)

type OrganizationClient struct {
	organizationRepository organization.RepositoryClient
}

func (oc *OrganizationClient) Register(request organization.CreateOrganizationRequest) error {
	newOrganization := organization.NewOrganization()
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

func (of *OrganizationFinder) ByUUID(uuid string) (*organization.Organization, error) {
	return of.oc.organizationRepository.FindByUUID(uuid)
}

func (of *OrganizationFinder) BySlug(slug string) (*organization.Organization, error) {
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

func (ps *PaystoreClient) ReceivePayment(request payment.ReceivePaymentRequest, selectedOrganization *organization.Organization) error {
	balanceFromDB, errFind := ps.balanceRepository.FindByUUID(request.AccountUUID)
	if errFind != nil {
		return errFind
	}
	if balanceFromDB == nil {
		return balance.AccountNotFound
	}

	if balanceFromDB.OrganizationUUID != selectedOrganization.GetUUID() {
		return organization.OrganizationMismatch
	}

	newPayment := payment.NewPayment()
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
		return balance.AccountNotFound
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

func New(writeDB *sql.DB, readDB *sql.DB, redis redis.UniversalClient,
	config *config.App, vendor config.Vendor, organization *organization.Organization) *PaystoreClient {
	var errInit error

	balanceRepo := balance.NewRepository(readDB, redis, config)
	paymentRepo, errInit := payment.NewRepository(readDB, redis, vendor, organization, config)
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
