package orders

import "errors"

var (
	// ErrInvalidStatus indicated an invalid order status
	ErrInvalidStatus = errors.New("invalid order status")

	// ErrConflict indicates that entity already exists.
	ErrConflict = errors.New("entity already exists")

	// ErrMalformedEntity indicates a malformed entity specification
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrCreateEntity indicates error in creating entity or entities
	ErrCreateEntity = errors.New("failed to create entity in the db")

	// ErrViewEntity indicates error in viewing entity or entities
	ErrViewEntity = errors.New("view entity failed")

	// ErrUpdateEntity indicates error in updating entity or entities
	ErrUpdateEntity = errors.New("update entity failed")

	// ErrRemoveEntity indicates error in removing entity
	ErrRemoveEntity = errors.New("failed to remove entity")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("entity not found")
)