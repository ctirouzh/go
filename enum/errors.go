package enum

import "errors"

var (
	// ErrNotRegisteredYet is returned when the given enum type is not registered yet.
	// Register the enum values to fix this error.
	ErrNotRegisteredYet = errors.New("not registered yet")
	// ErrInvalidValue is returned when the given value is not one of the registered values of the given enum type.
	ErrInvalidValue = errors.New("[Enum] invalid value")
)
