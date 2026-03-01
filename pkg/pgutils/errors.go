package pgutils

import (
	"errors"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrInvalidType         = errors.New("Invalid Input type")
	ErrNoRowsAffected      = errors.New("No rows affected")
	ErrForeignKeyViolation = errors.New("Foreign Key Violation")
	ErrUniqueViolation     = errors.New("Unique Key Violation")
	ErrNotNullViolation    = errors.New("Not Null Violation")
)

func CheckError(err error, w http.ResponseWriter) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	var status int

	switch {
	case errors.Is(err, ErrNoRowsAffected):
		status = http.StatusNotFound
	case errors.Is(err, pgx.ErrNoRows):
		status = http.StatusNotFound
	case errors.Is(err, pgx.ErrTooManyRows):
		status = http.StatusInternalServerError
	case errors.As(err, &pgErr):
		switch pgErr.Code {
		case "22P02":
			status = http.StatusUnprocessableEntity
			err = ErrInvalidType
		case "23502":
			status = http.StatusConflict
			err = ErrNotNullViolation
		case "23503":
			status = http.StatusConflict
			err = ErrForeignKeyViolation
		case "23505":
			status = http.StatusConflict
			err = ErrUniqueViolation
		default:
			log.Println(err)
			status = http.StatusInternalServerError
		}
	default:
		log.Println(err)
		status = http.StatusInternalServerError
	}

	if w != nil {
		w.WriteHeader(status)
	}

	return err
}
