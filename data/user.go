package data

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type User struct {
	ID           int       `db:"id"`
	Username     string    `db:"username"`
	RoleID       int       `db:"role_id"`
	Email        string    `db:"email"`
	Password     string    `db:"password"`
	FirstName    string    `db:"first_name"`
	LastName     string    `db:"last_name"`
	PhoneNumber  string    `db:"phone_number"`
	City         string    `db:"city"`
	State        *string   `db:"state"`
	CountryID    int       `db:"country_id"`
	ShareID      *int      `db:"share_id"`
	UpdatedAt    time.Time `db:"updated_at"`
	UpdatedBy    string    `db:"updated_by"`
	RegisteredAt time.Time `db:"registered_at"`
	RegisteredBy string    `db:"registered_by"`
}

type AuthCodeRequest struct {
	ID         int       `db:"id"`
	UserID     int       `db:"user_id"`
	AuthCode   uuid.UUID `db:"auth_code"`
	ValidUntil time.Time `db:"valid_until"`
}

type UserNotFoundError struct{}

func (e *UserNotFoundError) Error() string {
	return "User not found"
}

func (user *User) GetUser(connection *pgxpool.Conn, username string) error {
	row := connection.QueryRow(
		context.Background(),
		"SELECT * FROM authmantledb.us_user us WHERE us.username = $1",
		username,
	)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.RoleID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.City,
		&user.State,
		&user.CountryID,
		&user.ShareID,
		&user.UpdatedAt,
		&user.UpdatedBy,
		&user.RegisteredAt,
		&user.RegisteredBy,
	)
	if err != nil {
		return &UserNotFoundError{}
	}
	return nil
}

func (arc *AuthCodeRequest) CreateAuthCodeRequest(connection *pgxpool.Conn, userID int) error {
	row := connection.QueryRow(
		context.Background(),
		"INSERT INTO authmantledb.in_auth_code_requests (id, user_id) VALUES (nextval('authmantledb.in_auth_code_requests_id_seq'), $1) RETURNING *",
		userID,
	)
	err := row.Scan(
		&arc.ID,
		&arc.UserID,
		&arc.AuthCode,
		&arc.ValidUntil,
	)
	if err != nil {
		return err
	}
	return nil
}

func (arc *AuthCodeRequest) GetAuthCodeRequest(connection *pgxpool.Conn, code string) error {
	row := connection.QueryRow(
		context.Background(),
		"SELECT * FROM authmantledb.in_auth_code_requests WHERE auth_code = $1",
		code,
	)
	err := row.Scan(
		&arc.ID,
		&arc.UserID,
		&arc.AuthCode,
		&arc.ValidUntil,
	)
	if err != nil {
		return err
	}
	return nil
}

func CheckRedirectURI(connection *pgxpool.Conn, uri string) (bool, error) {
	row := connection.QueryRow(
		context.Background(),
		"SELECT u.redirect_uri FROM authmantledb.in_supp_auth_allowed_redirects u WHERE redirect_uri = $1",
		uri,
	)
	var redir string
	err := row.Scan(&redir)
	if err != nil {
		return false, err
	}
	return redir == uri, nil
}
