package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sebsvt/financial-service/logs"
	"github.com/sebsvt/financial-service/repository"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/account"
	"github.com/stripe/stripe-go/v79/accountlink"
)

type financialAccountService struct {
	repo repository.FinancialAccountRepository
}

func NewFinancialAccountService(repo repository.FinancialAccountRepository) FinancialAccountService {
	return financialAccountService{repo: repo}
}

// GetFinancialFromOrganisation implements FinancialAccountService.
func (srv financialAccountService) GetFinancialFromOrganisation(organisation_id string) (*FinancialAccountResponse, error) {
	fnc_acc, err := srv.repo.FromOrganisationID(organisation_id)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	return (*FinancialAccountResponse)(fnc_acc), nil
}

// SetUpFinancialAccount implements FinancialAccountService.
func (srv financialAccountService) SetUpFinancialAccount(organisation_id string) (string, error) {
	stripe.Key = "testingkeywaithingforrealkeyfromenv"

	// Create a new Stripe account
	params := &stripe.AccountParams{
		Type:         stripe.String(string(stripe.AccountTypeExpress)),
		BusinessType: stripe.String(string(stripe.AccountBusinessTypeIndividual)),
	}
	stripeAccount, err := account.New(params)
	if err != nil {
		logs.Error(fmt.Sprintf("Failed to create Stripe account: %v", err))
		return "", err
	}

	// Create an account link for onboarding
	linkParams := &stripe.AccountLinkParams{
		Account:    stripe.String(stripeAccount.ID),
		RefreshURL: stripe.String("https://aislena-mock-domain.com/refresh"),
		ReturnURL:  stripe.String("https://aiselena-mock-domain.com/success"),
		Type:       stripe.String("account_onboarding"),
	}
	accountLink, err := accountlink.New(linkParams)
	if err != nil {
		logs.Error(fmt.Sprintf("Failed to create account link: %v", err))
		return "", err
	}

	// Create a new FinancialAccount record to save in database
	financialAccount := repository.FinancialAccount{
		ID:              uuid.New().String(),
		StripeAccountID: stripeAccount.ID,
		OrganisationID:  organisation_id,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Save the FinancialAccount to the database
	if err := srv.repo.Save(&financialAccount); err != nil {
		logs.Error(fmt.Sprintf("Failed to save financial account: %v", err))
		return "", err
	}

	return accountLink.URL, nil
}
