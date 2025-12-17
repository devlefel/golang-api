package domain

import "errors"

var (
	ErrDeviceNotFound       = errors.New("device not found")
	ErrDeviceAlreadyExists  = errors.New("device already exists")
	ErrInvalidDeviceState   = errors.New("invalid device state")
	ErrImmutableField       = errors.New("field cannot be updated")
	ErrDeviceInUse          = errors.New("device is in use")
)
