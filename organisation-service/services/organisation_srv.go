package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sebsvt/organisation-service/logs"
	"github.com/sebsvt/organisation-service/repository"
)

type organisationService struct {
	org_repo    repository.OrganisationRepository
	member_repo repository.OrganisationMemberRepository
}

func NewOrganisationService(org_repo repository.OrganisationRepository, member_repo repository.OrganisationMemberRepository) OrganisationService {
	return organisationService{org_repo: org_repo, member_repo: member_repo}
}

// CreateNewOrganisation implements OrganisationService.
func (srv organisationService) CreateNewOrganisation(creator string, newOrganisation OrgansationCreated) (string, error) {
	// check if domain is already in used
	org_from_domain, err := srv.org_repo.FromDomain(newOrganisation.Domain)
	if err != nil && err != sql.ErrNoRows {
		logs.Error(err)
		return "", err
	}
	fmt.Println(org_from_domain)
	if org_from_domain != nil {
		return "", ErrDomainAlreadyinUsed
	}

	// create a new organisation
	now := time.Now()
	new_id := srv.org_repo.NextIdentity()
	if err := srv.org_repo.Save(&repository.Organisation{
		OrganisationID: new_id,
		Name:           newOrganisation.Name,
		Domain:         newOrganisation.Domain,
		Logo:           newOrganisation.Logo,
		CreatedBy:      creator,
		UpdatedAt:      now,
		CreatedAt:      now,
		DeletedAt:      nil,
	}); err != nil {
		logs.Error(err)
		return "", err
	}
	// add new member
	if err := srv.member_repo.Save(&repository.OrganisationMember{
		MemberID:       srv.member_repo.NextIdentity(),
		OrganisationID: new_id,
		UserID:         creator,
		Role:           repository.Admin,
		CreatedAt:      now,
		UpdatedAt:      now,
		DeletedAt:      nil,
	}); err != nil {
		logs.Error(err)
		return "", err
	}
	return new_id, nil
}

// GetOrganisationByDomain implements OrganisationService.
func (srv organisationService) GetOrganisationByDomain(domain string) (*OrganisationResponse, error) {
	org, err := srv.org_repo.FromDomain(domain)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	return (*OrganisationResponse)(org), nil
}

func (srv organisationService) GetAllOrganisationMember(organisation_id string, user_id string) ([]OrganisationMember, error) {
	var org_member []OrganisationMember
	include_member := false
	org_member_from_db, err := srv.member_repo.FromOrganisationID(organisation_id)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	for _, member := range org_member_from_db {
		if member.UserID == user_id {
			include_member = true
		}
		org_member = append(org_member, OrganisationMember{
			MemberID:       member.MemberID,
			OrganisationID: member.OrganisationID,
			UserID:         member.UserID,
			Role:           string(member.Role),
			CreatedAt:      member.CreatedAt,
			UpdatedAt:      member.UpdatedAt,
			DeletedAt:      member.DeletedAt,
		})
	}
	if !include_member {
		return nil, ErrPermissionDenied
	}
	return org_member, nil
}

func (srv organisationService) GetAllOrganisationFromUserID(user_id string) ([]OrganisationResponse, error) {
	var organisations []OrganisationResponse
	organisation_member, err := srv.member_repo.FromUserID(user_id)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	for _, org_mem := range organisation_member {
		org, err := srv.org_repo.FromID(org_mem.OrganisationID)
		if err != nil {
			return nil, err
		}
		organisations = append(organisations, OrganisationResponse(*org))
	}
	return organisations, nil
}
