package data

import (
	"context"
	"log/slog"
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

type UserNotFoundError struct{}

func (e *UserNotFoundError) Error() string {
	return "User not found"
}

func (user *User) GetUser(ctx context.Context, logger slog.Logger, connection DbActions, username string) error {
	row := connection.QueryRow(
		ctx,
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
