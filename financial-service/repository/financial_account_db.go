package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type financialAccountRepositoryDB struct {
	db *sqlx.DB
}

func NewFinancialAccountRepositoryDB(db *sqlx.DB) FinancialAccountRepository {
	return financialAccountRepositoryDB{db: db}
}

// FromOrganisationID implements FinancialAccountRepository.
func (repo financialAccountRepositoryDB) FromOrganisationID(organisation_id string) (*FinancialAccount, error) {
	var financial_account FinancialAccount
	query := `
		SELECT * from financial_accounts WHERE organisation_id=$1
	`
	if err := repo.db.Get(&financial_account, query, organisation_id); err != nil {
		return nil, err
	}
	return &financial_account, nil
}

// Save implements FinancialAccountRepository.
func (repo financialAccountRepositoryDB) Save(account *FinancialAccount) error {
	query := `
		INSERT INTO financial_accounts (id, stripe_account_id, organisation_id, created_at, updated_at)
		VALUES (:id, :stripe_account_id, :organisation_id, :created_at, :updated_at)
		ON CONFLICT (organisation_id)
		DO UPDATE SET
			stripe_account_id = EXCLUDED.stripe_account_id,
			updated_at = EXCLUDED.updated_at
	`

	// Set timestamps for created and updated times
	if account.CreatedAt.IsZero() {
		account.CreatedAt = time.Now()
	}
	account.UpdatedAt = time.Now()

	// Use NamedExec to map struct fields to query parameters
	_, err := repo.db.NamedExec(query, account)
	if err != nil {
		return err
	}
	return nil
}
