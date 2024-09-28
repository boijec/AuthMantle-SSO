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
type Realm struct {
	Id           int       `db:"id"`
	Name         string    `db:"name"`
	UpdatedAt    time.Time `db:"updated_at"`
	UpdatedBy    string    `db:"updated_by"`
	RegisteredAt time.Time `db:"registered_at"`
	RegisteredBy string    `db:"registered_by"`
}

type RealmCacheObject struct {
	RealmId      *int
	RealmName    string
	Redirects    []Redirect
	SubjectTypes []SubjectType
	GrantTypes   []GrantType
	Scopes       []Scope
	Claims       []Claim
	Audience     []Audience
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

// TODO Figure out how to malloc the settings strings without a count query prior to every query

func (rco *RealmCacheObject) GetRealmSettings(ctx context.Context, connection DbActions, realmName string) error {
	realmRow := connection.QueryRow(
		ctx,
		"SELECT r.id, r.name FROM authmantledb.realm r WHERE r.name = $1",
		realmName,
	)
	err := realmRow.Scan(
		&rco.RealmId,
		&rco.RealmName,
	)
	if err != nil {
		return err
	}

	return nil
}
