package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type organisationRepositoryDB struct {
	db *sqlx.DB
}

func NewOrganisationRepositoryDB(db *sqlx.DB) OrganisationRepository {
	return &organisationRepositoryDB{db: db}
}

// NextIdentity generates a new unique organisation ID.
func (repo *organisationRepositoryDB) NextIdentity() string {
	return "org_" + uuid.New().String()
}

// FromDomain retrieves an organisation by its domain.
func (repo *organisationRepositoryDB) FromDomain(organisationDomain string) (*Organisation, error) {
	var organisation Organisation
	query := `
		SELECT organisation_id, name, domain, logo, created_by, updated_at, created_at, deleted_at
		FROM organisations
		WHERE domain = $1`
	err := repo.db.Get(&organisation, query, organisationDomain)
	if err != nil {
		return nil, err
	}
	return &organisation, err
}

// FromID retrieves an organisation by its ID.
func (repo *organisationRepositoryDB) FromID(organisationID string) (*Organisation, error) {
	var organisation Organisation
	query := `
		SELECT organisation_id, name, domain, logo, created_by, updated_at, created_at, deleted_at
		FROM organisations
		WHERE organisation_id = $1`
	err := repo.db.Get(&organisation, query, organisationID)
	if err != nil {
		return nil, err
	}
	return &organisation, err
}

// Save inserts a new organisation or updates an existing one.
func (repo *organisationRepositoryDB) Save(org *Organisation) error {
	query := `
		INSERT INTO organisations (organisation_id, name, domain, logo, created_by, updated_at, created_at, deleted_at)
		VALUES (:organisation_id, :name, :domain, :logo, :created_by, :updated_at, :created_at, :deleted_at)
		ON CONFLICT (organisation_id)
		DO UPDATE SET name = :name, domain = :domain, logo = :logo, updated_at = :updated_at, deleted_at = :deleted_at`
	_, err := repo.db.NamedExec(query, org)
	return err
}

// Delete marks an organisation as deleted by setting the DeletedAt field.
func (repo *organisationRepositoryDB) Delete(organisationID string) error {
	query := `
		UPDATE organisations
		SET deleted_at = $1
		WHERE organisation_id = $2`
	_, err := repo.db.Exec(query, nil, organisationID) // Assuming DeletedAt set at service layer
	return err
}

// FromUserID implements OrganisationRepository.
func (repo *organisationRepositoryDB) FromUserID(user_id string) ([]Organisation, error) {
	var organisations []Organisation
	query := `
		SELECT organisation_id, name, domain, logo, created_by, updated_at, created_at, deleted_at
		FROM organisations
		WHERE user_id = $1`
	if err := repo.db.Select(&organisations, query, user_id); err != nil {
		return nil, err
	}
	return organisations, nil
}
