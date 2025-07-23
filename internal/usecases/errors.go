package usecases

import "errors"

var ErrEntityNotFound = errors.New("entity not found")
var ErrEntityAlreadyExists = errors.New("entity already exists")
var ErrInvalidDateFormat = errors.New("invalid date format")
var ErrInvalidUUID = errors.New("invalid UUID format")
