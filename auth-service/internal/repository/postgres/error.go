package repositories

import (
	"auth-service/internal/utils"

	"github.com/jackc/pgx/v5/pgconn"
)

func ClassifyDBerror(err error) *utils.APIError {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		switch pgErr.Code {
		case "23505":
			return utils.NewAPIError(409, "Resource already exists", "")
		default:
			return utils.NewAPIError(500, "Database error", "")
		}
	}
	return utils.NewAPIError(500, "Internal server error", "")
}
