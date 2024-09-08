package data

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type AuthCodeRequest struct {
	ID         int       `db:"id"`
	UserID     int       `db:"user_id"`
	AuthCode   uuid.UUID `db:"auth_code"`
	ValidUntil time.Time `db:"valid_until"`
	Consumed   int       `db:"consumed"`
}

func (arc *AuthCodeRequest) CreateAuthCodeRequest(ctx context.Context, logger slog.Logger, connection DbActions, userID int) error {
	row := connection.QueryRow(
		ctx,
		"INSERT INTO authmantledb.in_auth_code_requests (id, user_id) VALUES (nextval('authmantledb.in_auth_code_requests_id_seq'), $1) RETURNING *",
		userID,
	)
	logger.DebugContext(ctx, "Auth code row was inserted")
	err := row.Scan(
		&arc.ID,
		&arc.UserID,
		&arc.AuthCode,
		&arc.ValidUntil,
		&arc.Consumed,
	)
	if err != nil {
		return err
	}
	logger.DebugContext(ctx, "Auth code row was scanned without errors")
	return nil
}
func (arc *AuthCodeRequest) ConsumeAuthCodeRequest(ctx context.Context, logger slog.Logger, connection DbActions) error {
	_, err := connection.Exec(
		ctx,
		"UPDATE authmantledb.in_auth_code_requests SET consumed = 1 WHERE auth_code = $1",
		arc.AuthCode,
	)
	logger.DebugContext(ctx, "Auth code query was executed")
	if err != nil {
		return err
	}
	logger.DebugContext(ctx, "Auth code query was executed without errors")
	return nil
}

func (arc *AuthCodeRequest) GetAuthCodeRequest(ctx context.Context, logger slog.Logger, connection DbActions, code string) error {
	row := connection.QueryRow(
		ctx,
		"SELECT * FROM authmantledb.in_auth_code_requests WHERE auth_code = $1",
		code,
	)
	logger.DebugContext(ctx, "Auth code row was queried")
	err := row.Scan(
		&arc.ID,
		&arc.UserID,
		&arc.AuthCode,
		&arc.ValidUntil,
		&arc.Consumed,
	)
	if err != nil {
		return err
	}
	logger.DebugContext(ctx, "Auth code row was scanned without errors")
	return nil
}
