package services

import (
	"errors"
	"time"
)

var (
	ErrDomainAlreadyinUsed = errors.New("domain already in used")
	ErrPermissionDenied    = errors.New("permission denied")
)

type OrgansationCreated struct {
	Name   string `json:"name"`
	Domain string `json:"domain"`
	Logo   string `json:"logo"`
}

type OrganisationResponse struct {
	OrganisationID string     `json:"organisation_id"`
	Name           string     `json:"name"`
	Domain         string     `json:"domain"`
	Logo           string     `json:"logo"`
	CreatedBy      string     `json:"created_by"`
	UpdatedAt      time.Time  `json:"updated_at"`
	CreatedAt      time.Time  `json:"created_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}

type OrganisationMember struct {
	MemberID       string     `json:"member_id"`
	OrganisationID string     `json:"organisation_id"`
	UserID         string     `json:"user_id"`
	Role           string     `json:"role"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}

type OrganisationService interface {
	CreateNewOrganisation(creator string, newOrganisation OrgansationCreated) (string, error)
	GetOrganisationByDomain(domain string) (*OrganisationResponse, error)
	GetAllOrganisationMember(organisation_id string, user_id string) ([]OrganisationMember, error)
	GetAllOrganisationFromUserID(user_id string) ([]OrganisationResponse, error)
}
