package repository

import "time"

var (
	Admin  Role = "admin"
	Owner  Role = "owner"
	Member Role = "member"
	Guest  Role = "guest"
)

type Role string

type OrganisationMember struct {
	MemberID       string     `db:"member_id"`
	OrganisationID string     `db:"organisation_id"`
	UserID         string     `db:"user_id"`
	Role           Role       `db:"role"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at"`
}

type OrganisationMemberRepository interface {
	NextIdentity() string
	FromOrganisationID(organisationID string) ([]OrganisationMember, error)
	FromUserID(user_id string) ([]OrganisationMember, error)
	Save(*OrganisationMember) error
	Delete(organisationID string) error
}
