package errors

// Errors defined in this file are used by the LoggingErrorEncoder decorator
// to distinguish and log API request validation errors and avoid that service
// errors are logged twice.
var (
	// ErrAuthentication indicates failure occurred while authenticating the entity.
	ErrAuthentication = New("failed to perform authentication over the entity")

	// ErrAuthorization indicates failure occurred while authorizing the entity.
	ErrAuthorization = New("failed to perform authorization over the entity")

	// ErrBearerToken indicates missing or invalid bearer user token.
	ErrBearerToken = New("missing or invalid bearer user token")

	// ErrMissingID indicates missing entity ID.
	ErrMissingID = New("missing entity id")

	// ErrInvalidAuthKey indicates invalid auth key.
	ErrInvalidAuthKey = New("invalid auth key")

	// ErrNameSize indicates that name size exceeds the max.
	ErrNameSize = New("invalid name size")

	// ErrLimitSize indicates that an invalid limit.
	ErrLimitSize = New("invalid limit size")

	// ErrOffsetSize indicates an invalid offset.
	ErrOffsetSize = New("invalid offset size")

	// ErrInvalidQueryParams indicates invalid query parameters
	ErrInvalidQueryParams = New("invalid query parameters")

	// ErrNotFoundParam indicates that the parameter was not found in the query
	ErrNotFoundParam = New("parameter not found in the query")

	// ErrMalformedEntity indicates a malformed entity specification
	ErrMalformedEntity = New("malformed entity specification")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = New("entity not found")

	// ErrConflict indicates that entity already exists.
	ErrConflict = New("entity already exists")

	// ErrCreateEntity indicates error in creating entity or entities
	ErrCreateEntity = New("failed to create entity in the db")

	// ErrViewEntity indicates error in viewing entity or entities
	ErrViewEntity = New("view entity failed")

	// ErrUpdateEntity indicates error in updating entity or entities
	ErrUpdateEntity = New("update entity failed")

	// ErrRemoveEntity indicates error in removing entity
	ErrRemoveEntity = New("failed to remove entity")

	// ErrScanMetadata indicates problem with metadata in db
	ErrScanMetadata = New("failed to scan metadata in db")

	// ErrUnsupportedContentType indicates unacceptable or lack of Content-Type
	ErrUnsupportedContentType = New("unsupported content type")
)
