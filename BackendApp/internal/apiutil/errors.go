package apiutil

import "github.com/0x6flab/jikoniApp/BackendApp/internal/errors"

// ErrorRes represents the HTTP error response body.
type ErrorRes struct {
	Err string `json:"error"`
}

// Errors defined in this file are used by the LoggingErrorEncoder decorator
// to distinguish and log API request validation errors and avoid that service
// errors are logged twice.
var (
	// ErrAuthentication indicates failure occurred while authenticating the entity.
	ErrAuthentication = errors.New("failed to perform authentication over the entity")

	// ErrAuthorization indicates failure occurred while authorizing the entity.
	ErrAuthorization = errors.New("failed to perform authorization over the entity")

	// ErrBearerToken indicates missing or invalid bearer user token.
	ErrBearerToken = errors.New("missing or invalid bearer user token")

	// ErrMissingID indicates missing entity ID.
	ErrMissingID = errors.New("missing entity id")

	// ErrInvalidAuthKey indicates invalid auth key.
	ErrInvalidAuthKey = errors.New("invalid auth key")

	// ErrNameSize indicates that name size exceeds the max.
	ErrNameSize = errors.New("invalid name size")

	// ErrLimitSize indicates that an invalid limit.
	ErrLimitSize = errors.New("invalid limit size")

	// ErrOffsetSize indicates an invalid offset.
	ErrOffsetSize = errors.New("invalid offset size")

	// ErrInvalidQueryParams indicates invalid query parameters
	ErrInvalidQueryParams = errors.New("invalid query parameters")

	// ErrNotFoundParam indicates that the parameter was not found in the query
	ErrNotFoundParam = errors.New("parameter not found in the query")

	// ErrMalformedEntity indicates a malformed entity specification
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("entity not found")

	// ErrConflict indicates that entity already exists.
	ErrConflict = errors.New("entity already exists")

	// ErrCreateEntity indicates error in creating entity or entities
	ErrCreateEntity = errors.New("failed to create entity in the db")

	// ErrViewEntity indicates error in viewing entity or entities
	ErrViewEntity = errors.New("view entity failed")

	// ErrUpdateEntity indicates error in updating entity or entities
	ErrUpdateEntity = errors.New("update entity failed")

	// ErrRemoveEntity indicates error in removing entity
	ErrRemoveEntity = errors.New("failed to remove entity")

	// ErrScanMetadata indicates problem with metadata in db
	ErrScanMetadata = errors.New("failed to scan metadata in db")

	// ErrUnsupportedContentType indicates unacceptable or lack of Content-Type
	ErrUnsupportedContentType = errors.New("unsupported content type")
)
