package util

import (
	"errors"

	"github.com/lib/pq"
)

func IsUniqueViolation(err error) bool {
	var pgErr *pq.Error
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
