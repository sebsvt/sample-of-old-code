package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// OrganisationMemberRepository handles OrganisationMember persistence.
type organisationMemberRepository struct {
	db *sqlx.DB
}

// NewOrganisationMemberRepository initializes a new OrganisationMemberRepository.
func NewOrganisationMemberRepository(db *sqlx.DB) OrganisationMemberRepository {
	return organisationMemberRepository{db: db}
}

// NextIdentity generates a new unique identifier for an organisation member.
func (repo organisationMemberRepository) NextIdentity() string {
	return "member_" + uuid.New().String()
}

// FromOrganisationID retrieves all active members for a specified organisation.
func (repo organisationMemberRepository) FromOrganisationID(organisationID string) ([]OrganisationMember, error) {
	var members []OrganisationMember
	query := `
		SELECT member_id, organisation_id, user_id, role, created_at, updated_at, deleted_at
		FROM organisation_members
		WHERE organisation_id = $1 AND deleted_at IS NULL`
	err := repo.db.Select(&members, query, organisationID)
	if err != nil {
		return nil, err
	}
	return members, nil
}

// Save inserts a new organisation member or updates an existing one using UPSERT.
func (repo organisationMemberRepository) Save(member *OrganisationMember) error {
	query := `
		INSERT INTO organisation_members (member_id, organisation_id, user_id, role, created_at, updated_at)
		VALUES (:member_id, :organisation_id, :user_id, :role, :created_at, :updated_at)
		ON CONFLICT (member_id)
		DO UPDATE SET
			organisation_id = EXCLUDED.organisation_id,
			user_id = EXCLUDED.user_id,
			role = EXCLUDED.role,
			updated_at = EXCLUDED.updated_at`
	_, err := repo.db.NamedExec(query, member)
	return err
}

// Delete sets the DeletedAt timestamp for a member by ID.
func (repo organisationMemberRepository) Delete(memberID string) error {
	now := time.Now()
	query := `
		UPDATE organisation_members
		SET deleted_at = $1
		WHERE member_id = $2`
	result, err := repo.db.Exec(query, now, memberID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no member found with the given ID")
	}
	return nil
}

// FromID retrieves a member by their unique member ID.
func (repo organisationMemberRepository) FromID(memberID string) (*OrganisationMember, error) {
	var member OrganisationMember
	query := `
		SELECT member_id, organisation_id, user_id, role, created_at, updated_at, deleted_at
		FROM organisation_members
		WHERE member_id = $1`
	err := repo.db.Get(&member, query, memberID)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// FromUserID implements OrganisationMemberRepository.
func (repo organisationMemberRepository) FromUserID(user_id string) ([]OrganisationMember, error) {
	var members []OrganisationMember

	query := `
		SELECT member_id, organisation_id, user_id, role, created_at, updated_at, deleted_at
		FROM organisation_members
		WHERE user_id = $1`

	if err := repo.db.Select(&members, query, user_id); err != nil {
		return nil, err
	}

	return members, nil
}
