package data

import (
	"context"
	"log/slog"
)

type Audience struct {
	Id       int    `db:"id"`
	Audience string `db:"audience"`
}
type GrantType struct {
	Id        int    `db:"id"`
	GrantType string `db:"grant_type"`
}
type Scope struct {
	Id    int    `db:"id"`
	Scope string `db:"scope"`
}
type Claim struct {
	Id    int    `db:"id"`
	Claim string `db:"claim"`
}
type SubjectType struct {
	Id          int    `db:"id"`
	SubjectType string `db:"subject_type"`
}
type Redirect struct {
	Id          int    `db:"id"`
	RedirectURI string `db:"redirect_uri"`
}

func CheckRedirectURI(ctx context.Context, logger slog.Logger, connection DbActions, uri string) (bool, error) {
	row := connection.QueryRow(
		ctx,
		"SELECT u.redirect_uri FROM authmantledb.in_supp_auth_allowed_redirects u WHERE redirect_uri = $1",
		uri,
	)
	logger.DebugContext(ctx, "Redirect row was queried")
	var redir string
	err := row.Scan(&redir)
	if err != nil {
		return false, err
	}
	logger.DebugContext(ctx, "Redirect row was scanned without errors")
	return redir == uri, nil
}
