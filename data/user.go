package data

import (
	"context"
	"time"
)

type User struct {
	ID           int       `db:"id"`
	Username     string    `db:"username"`
	RealmID      int       `db:"realm_id"`
	RoleID       int       `db:"role_id"`
	Email        string    `db:"email"`
	Password     string    `db:"password"`
	FirstName    string    `db:"first_name"`
	LastName     string    `db:"last_name"`
	PhoneNumber  string    `db:"phone_number"`
	City         string    `db:"city"`
	State        *string   `db:"state"`
	CountryID    int       `db:"country_id"`
	UpdatedAt    time.Time `db:"updated_at"`
	UpdatedBy    string    `db:"updated_by"`
	RegisteredAt time.Time `db:"registered_at"`
	RegisteredBy string    `db:"registered_by"`
}

type UserNotFoundError struct{}

func (e *UserNotFoundError) Error() string {
	return "User not found"
}

func (user *User) GetUser(ctx context.Context, connection DbActions, username string, realmId int) error {
	row := connection.QueryRow(
		ctx,
		"SELECT us.* FROM authmantledb.user us WHERE us.username = $1 AND us.realm_id = $2",
		username,
		realmId,
	)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.RealmID,
		&user.RoleID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.City,
		&user.State,
		&user.CountryID,
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
