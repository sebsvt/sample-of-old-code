package repository

import "time"

type Organisation struct {
	OrganisationID string     `db:"organisation_id"`
	Name           string     `db:"name"`
	Domain         string     `db:"domain"`
	Logo           string     `db:"logo"`
	CreatedBy      string     `db:"created_by"`
	UpdatedAt      time.Time  `db:"updated_at"`
	CreatedAt      time.Time  `db:"created_at"`
	DeletedAt      *time.Time `db:"deleted_at"`
}

type OrganisationRepository interface {
	NextIdentity() string
	FromUserID(user_id string) ([]Organisation, error)
	FromID(orgnisationID string) (*Organisation, error)
	FromDomain(organisationDomain string) (*Organisation, error)
	// save function could be create and update in the same function
	Save(org *Organisation) error
	Delete(organisationID string) error
}
