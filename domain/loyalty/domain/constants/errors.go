package constants

import "errors"

var (
	// ErrUsuarioNoEncontrado se devuelve cuando no se encuentra un usuario en la base de datos
	ErrUserNotFound      = errors.New("user not found")
	ErrPointsInsuficient = errors.New("Insuficient points")
)
