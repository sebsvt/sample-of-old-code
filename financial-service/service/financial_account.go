package service

import "time"

type FinancialAccountResponse struct {
	ID              string     `json:"id"`
	StripeAccountID string     `json:"stripe_account_id"`
	OrganisationID  string     `json:"organisation_id"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_aqt"`
}

type FinancialAccountService interface {
	SetUpFinancialAccount(organisation_id string) (string, error)
	GetFinancialFromOrganisation(organisation_id string) (*FinancialAccountResponse, error)
}
