package storage

import "errors"

var ErrNotExists = errors.New("doesn't exist")

func IsNotExists(err error) bool {
	return err == ErrNotExists
}
