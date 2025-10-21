package organization

import (
	"database/sql"
	"github.com/21strive/redifu"
	"github.com/redis/go-redis/v9"
	"paystore/config"
	"paystore/lib/model"
)

var createOrganizationQuery = `INSERT INTO organization (name, slug, avatar) VALUES ($1, $2, $3)`
var updateOrganizationQuery = `UPDATE organization SET name = $1, slug = $2, avatar = $3 WHERE slug = $4`
var findOrganizationByUUIDQuery = `SELECT name, slug, avatar FROM organization WHERE uuid = $1`
var findOrganizationBySlugQuery = `SELECT name, slug, avatar FROM organization WHERE slug = $1`

type RepositoryClient interface {
	Create(organization *model.Organization) error
	Update(organization *model.Organization) error
	FindByUUID(uuid string) (*model.Organization, error)
	FindBySlug(slug string) (*model.Organization, error)
}

type Repository struct {
	writeDB                    *sql.DB
	readDB                     *sql.DB
	base                       *redifu.Base[*model.Organization]
	createOrganizationStmt     *sql.Stmt
	updateOrganizationStmt     *sql.Stmt
	findOrganizationByUUIDStmt *sql.Stmt
	findOrganizationBySlugStmt *sql.Stmt
}

func (or *Repository) Create(organization *model.Organization) error {
	_, err := or.createOrganizationStmt.Exec(organization.Name, organization.Slug, organization.Avatar)
	return err
}

func (or *Repository) Update(organization *model.Organization) error {
	_, err := or.updateOrganizationStmt.Exec(organization.Name, organization.Slug, organization.Avatar, organization.Slug)
	return err
}

func (or *Repository) FindByUUID(uuid string) (*model.Organization, error) {
	row, errScan := OrganizationRowScanner(or.findOrganizationByUUIDStmt.QueryRow(uuid))
	if errScan != nil {
		return nil, errScan
	}

	errSet := or.base.Set(row)
	if errSet != nil {
		return nil, errSet
	}
	return row, nil
}

func (or *Repository) FindBySlug(slug string) (*model.Organization, error) {
	row, errScan := OrganizationRowScanner(or.findOrganizationBySlugStmt.QueryRow(slug))
	if errScan != nil {
		return nil, errScan
	}

	errSet := or.base.Set(row)
	if errSet != nil {
		return nil, errSet
	}
	return row, nil
}

func OrganizationRowScanner(row *sql.Row) (*model.Organization, error) {
	org := model.NewOrganization()
	err := row.Scan(&org.Name, &org.Slug, &org.Avatar)
	if err != nil {
		return nil, err
	}

	return org, nil
}

func NewRepository(writeDB *sql.DB, readDB *sql.DB, redis redis.UniversalClient, config *config.App) *Repository {
	base := redifu.NewBase[*model.Organization](redis, "organization:%s", config.RecordAge)

	createOrganizationStmt, err := writeDB.Prepare(createOrganizationQuery)
	if err != nil {
		panic(err)
	}
	updateOrganizationStmt, err := writeDB.Prepare(updateOrganizationQuery)
	if err != nil {
		panic(err)
	}
	findOrganizationByUUIDStmt, err := readDB.Prepare(findOrganizationByUUIDQuery)
	if err != nil {
		panic(err)
	}
	findOrganizationBySlugStmt, err := readDB.Prepare(findOrganizationBySlugQuery)
	if err != nil {
		panic(err)
	}

	organizationRepo := &Repository{
		writeDB:                    writeDB,
		readDB:                     readDB,
		base:                       base,
		createOrganizationStmt:     createOrganizationStmt,
		updateOrganizationStmt:     updateOrganizationStmt,
		findOrganizationByUUIDStmt: findOrganizationByUUIDStmt,
		findOrganizationBySlugStmt: findOrganizationBySlugStmt,
	}

	return organizationRepo
}
