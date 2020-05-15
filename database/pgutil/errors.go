package pgutil

import "strings"

const DuplicateKeyErrorPrefix = "ERROR #23505 duplicate key value violates unique constraint"

func ErrIsDuplicateKey(err error) bool {
	if strings.Index(err.Error(), DuplicateKeyErrorPrefix) == 0 {
		return true
	}
	return false
}

func NilifyErrDuplicateKey(err error) error {
	if err != nil && ErrIsDuplicateKey(err) {
		return nil
	}
	return err
}
