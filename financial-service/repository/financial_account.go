package repository

import "time"

type FinancialAccount struct {
	ID              string     `db:"id"`
	StripeAccountID string     `db:"stripe_account_id"`
	OrganisationID  string     `db:"organisation_id"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at"`
}

type FinancialAccountRepository interface {
	FromOrganisationID(string) (*FinancialAccount, error)
	Save(*FinancialAccount) error
}
