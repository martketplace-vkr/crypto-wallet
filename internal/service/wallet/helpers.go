package wallet

import (
	"database/sql"
	"strings"
)

func isNotFound(err error) bool {
	return err == sql.ErrNoRows || errorsIs(err, ErrNotFound)
}

func isUniqueRace(err error) bool {
	if err == nil {
		return false
	}

	msg := err.Error()
	return strings.Contains(msg, "duplicate key") || strings.Contains(msg, "unique constraint")
}

func errorsIs(err error, target error) bool {
	type causer interface{ Unwrap() error }

	for err != nil {
		if err == target {
			return true
		}

		u, ok := err.(causer)
		if !ok {
			return false
		}
		err = u.Unwrap()
	}

	return false
}
