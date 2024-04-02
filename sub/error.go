package sub

import "errors"

var (
	ErrIncorrectValueInterval = errors.New("interval value is incorrect")
	ErrInvalidConn            = errors.New("invalid connection")
	ErrInvalidAttackType      = errors.New("invalid a type of attack")
	ErrNoArguments            = errors.New("no arguments specified")
	ErrNoArgumentValue        = errors.New("no argument value specified")
)
