package data

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type AuthCodeRequest struct {
	ID           int       `db:"id"`
	UserID       int       `db:"user_id"`
	AuthCode     uuid.UUID `db:"auth_code"`
	ValidUntil   time.Time `db:"valid_until"`
	Consumed     int       `db:"consumed"`
	UpdatedAt    time.Time `db:"updated_at"`
	UpdatedBy    string    `db:"updated_by"`
	RegisteredAt time.Time `db:"registered_at"`
	RegisteredBy string    `db:"registered_by"`
}

type AuthCodeError struct {
	Reason string
}

func (e *AuthCodeError) Error() string {
	return e.Reason
}

func (arc *AuthCodeRequest) Validate() error {
	if arc.Consumed == 1 {
		return &AuthCodeError{"Auth code request has already been consumed"}
	}
	if time.Now().After(arc.ValidUntil) {
		return &AuthCodeError{"Auth code request has expired"}
	}
	return nil
}

func (arc *AuthCodeRequest) CreateAuthCodeRequest(ctx context.Context, connection DbActions, userID int) error {
	row := connection.QueryRow(
		ctx,
		"INSERT INTO authmantledb.auth_code_requests (id, user_id, updated_by, registered_by) VALUES (nextval('authmantledb.auth_code_requests_id_seq'), $1, $2, $2) RETURNING *",
		userID,
		"system:CreateAuthCodeRequest",
	)
	err := row.Scan(
		&arc.ID,
		&arc.UserID,
		&arc.AuthCode,
		&arc.ValidUntil,
		&arc.Consumed,
		&arc.UpdatedAt,
		&arc.UpdatedBy,
		&arc.RegisteredAt,
		&arc.RegisteredBy,
	)
	if err != nil {
		return err
	}
	return nil
}
func (arc *AuthCodeRequest) ConsumeAuthCodeRequest(ctx context.Context, connection DbActions) error {
	_, err := connection.Exec(
		ctx,
		"UPDATE authmantledb.auth_code_requests SET consumed = 1, updated_by = 'system:ConsumeAuthCodeRequest', updated_at = now() WHERE auth_code = $1",
		arc.AuthCode,
	)
	if err != nil {
		return err
	}
	return nil
}

func (arc *AuthCodeRequest) GetAuthCodeRequest(ctx context.Context, connection DbActions, code string) error {
	row := connection.QueryRow(
		ctx,
		"SELECT * FROM authmantledb.auth_code_requests WHERE auth_code = $1",
		code,
	)
	err := row.Scan(
		&arc.ID,
		&arc.UserID,
		&arc.AuthCode,
		&arc.ValidUntil,
		&arc.Consumed,
		&arc.UpdatedAt,
		&arc.UpdatedBy,
		&arc.RegisteredAt,
		&arc.RegisteredBy,
	)
	if err != nil {
		return err
	}
	return nil
}
