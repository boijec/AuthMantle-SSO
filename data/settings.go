package data

import (
	"context"
	"time"
)

type Audience struct {
	Id           int       `db:"id"`
	Audience     string    `db:"audience"`
	UpdatedAt    time.Time `db:"updated_at"`
	UpdatedBy    string    `db:"updated_by"`
	RegisteredAt time.Time `db:"registered_at"`
	RegisteredBy string    `db:"registered_by"`
}
type GrantType struct {
	Id           int       `db:"id"`
	GrantType    string    `db:"grant_type"`
	UpdatedAt    time.Time `db:"updated_at"`
	UpdatedBy    string    `db:"updated_by"`
	RegisteredAt time.Time `db:"registered_at"`
	RegisteredBy string    `db:"registered_by"`
}
type Scope struct {
	Id           int       `db:"id"`
	Scope        string    `db:"scope"`
	UpdatedAt    time.Time `db:"updated_at"`
	UpdatedBy    string    `db:"updated_by"`
	RegisteredAt time.Time `db:"registered_at"`
	RegisteredBy string    `db:"registered_by"`
}
type Claim struct {
	Id           int       `db:"id"`
	Claim        string    `db:"claim"`
	UpdatedAt    time.Time `db:"updated_at"`
	UpdatedBy    string    `db:"updated_by"`
	RegisteredAt time.Time `db:"registered_at"`
	RegisteredBy string    `db:"registered_by"`
}
type SubjectType struct {
	Id           int       `db:"id"`
	SubjectType  string    `db:"subject_type"`
	UpdatedAt    time.Time `db:"updated_at"`
	UpdatedBy    string    `db:"updated_by"`
	RegisteredAt time.Time `db:"registered_at"`
	RegisteredBy string    `db:"registered_by"`
}
type Redirect struct {
	Id           int       `db:"id"`
	RedirectURI  string    `db:"redirect_uri"`
	UpdatedAt    time.Time `db:"updated_at"`
	UpdatedBy    string    `db:"updated_by"`
	RegisteredAt time.Time `db:"registered_at"`
	RegisteredBy string    `db:"registered_by"`
}

func CheckRedirectURI(ctx context.Context, connection DbActions, uri string) (bool, error) {
	row := connection.QueryRow(
		ctx,
		"SELECT u.redirect_uri FROM authmantledb.supp_auth_allowed_redirects u WHERE redirect_uri = $1",
		uri,
	)
	var redir string
	err := row.Scan(&redir)
	if err != nil {
		return false, err
	}
	return redir == uri, nil
}
