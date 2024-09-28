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
	RealmId          *int
	RealmName        string
	Redirects        []string
	SubjectTypes     []string
	GrantTypes       []string
	Scopes           []string
	Claims           []string
	Audience         []string
	ResponseTypes    []string
	TokenSigningAlgs []string
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

func (rco *RealmCacheObject) GetRealmSettings(ctx context.Context, connection DbActions, realmId int) error {
	realmRow := connection.QueryRow(
		ctx,
		"SELECT r.id, r.name FROM authmantledb.realm r WHERE r.id = $1",
		realmId,
	)
	err := realmRow.Scan(
		&rco.RealmId,
		&rco.RealmName,
	)
	if err != nil {
		return err
	}
	redirects, err := connection.Query(
		ctx,
		"SELECT r.redirect_uri FROM authmantledb.supp_auth_allowed_redirects r WHERE r.realm_id = $1",
		realmId,
	)
	defer redirects.Close()
	if err != nil {
		return err
	}
	for redirects.Next() {
		var redir string
		err = redirects.Scan(
			&redir,
		)
		if err != nil {
			return err
		}
		rco.Redirects = append(rco.Redirects, redir)
	}
	subjectTypes, err := connection.Query(
		ctx,
		"SELECT s.subject_type FROM authmantledb.supp_auth_subject_types s WHERE s.realm_id = $1",
		realmId,
	)
	defer subjectTypes.Close()
	if err != nil {
		return err
	}
	for subjectTypes.Next() {
		var st string
		err = subjectTypes.Scan(
			&st,
		)
		if err != nil {
			return err
		}
		rco.SubjectTypes = append(rco.SubjectTypes, st)
	}
	grantTypes, err := connection.Query(
		ctx,
		"SELECT g.grant_type FROM authmantledb.supp_auth_grant_types g WHERE g.realm_id = $1",
		realmId,
	)
	defer grantTypes.Close()
	if err != nil {
		return err
	}
	for grantTypes.Next() {
		var gt string
		err = grantTypes.Scan(
			&gt,
		)
		if err != nil {
			return err
		}
		rco.GrantTypes = append(rco.GrantTypes, gt)
	}
	scopes, err := connection.Query(
		ctx,
		"SELECT s.scope FROM authmantledb.supp_auth_scopes s WHERE s.realm_id = $1",
		realmId,
	)
	defer scopes.Close()
	if err != nil {
		return err
	}
	for scopes.Next() {
		var sc string
		err = scopes.Scan(
			&sc,
		)
		if err != nil {
			return err
		}
		rco.Scopes = append(rco.Scopes, sc)
	}
	claims, err := connection.Query(
		ctx,
		"SELECT c.claim FROM authmantledb.supp_auth_claims c WHERE c.realm_id = $1",
		realmId,
	)
	defer claims.Close()
	if err != nil {
		return err
	}
	for claims.Next() {
		var cl string
		err = claims.Scan(
			&cl,
		)
		if err != nil {
			return err
		}
		rco.Claims = append(rco.Claims, cl)
	}
	audience, err := connection.Query(
		ctx,
		"SELECT a.audience FROM authmantledb.supp_auth_audience a WHERE a.realm_id = $1",
		realmId,
	)
	defer audience.Close()
	if err != nil {
		return err
	}
	for audience.Next() {
		var aud string
		err = audience.Scan(
			&aud,
		)
		if err != nil {
			return err
		}
		rco.Audience = append(rco.Audience, aud)
	}
	responseTypes, err := connection.Query(
		ctx,
		"SELECT r.response_type FROM authmantledb.supp_auth_response_types r WHERE r.realm_id = $1",
		realmId,
	)
	defer responseTypes.Close()
	if err != nil {
		return err
	}
	for responseTypes.Next() {
		var rt string
		err = responseTypes.Scan(
			&rt,
		)
		if err != nil {
			return err
		}
		rco.ResponseTypes = append(rco.ResponseTypes, rt)
	}
	tokenSigningAlgs, err := connection.Query(
		ctx,
		"SELECT t.signing_alg FROM authmantledb.supp_id_token_signing_alg_values t WHERE t.realm_id = $1",
		realmId,
	)
	defer tokenSigningAlgs.Close()
	if err != nil {
		return err
	}
	for tokenSigningAlgs.Next() {
		var tsa string
		err = tokenSigningAlgs.Scan(
			&tsa,
		)
		if err != nil {
			return err
		}
		rco.TokenSigningAlgs = append(rco.TokenSigningAlgs, tsa)
	}

	return nil
}
