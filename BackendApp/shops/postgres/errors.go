package postgres

import (
	"github.com/0x6flab/jikoniApp/BackendApp/internal/errors"
	"github.com/jackc/pgconn"
	"go.uber.org/multierr"
)

// Postgres error codes:
// https://www.postgresql.org/docs/current/errcodes-appendix.html
const (
	errDuplicate  = "23505" // unique_violation
	errTruncation = "22001" // string_data_right_truncation
	errFK         = "23503" // foreign_key_violation
	errInvalid    = "22P02" // invalid_text_representation
)

func handleError(err, wrapper error) error {
	pqErr, ok := err.(*pgconn.PgError)
	if ok {
		switch pqErr.Code {
		case errDuplicate:
			return multierr.Combine(errors.ErrConflict, err)
		case errInvalid, errTruncation:
			return multierr.Combine(errors.ErrMalformedEntity, err)
		case errFK:
			return multierr.Combine(errors.ErrCreateEntity, err)
		}
	}
	return multierr.Combine(wrapper, err)
}
