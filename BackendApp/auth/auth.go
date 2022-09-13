package auth

import (
	context "context"
	"time"
)

type IssueToken struct {
	ID          string `json:"id,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
}

type User struct {
	ID          string    `json:"id,omitempty"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	Name        string    `json:"name,omitempty"`
	Type        string    `json:"type,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

type Policy struct {
	User      string    `json:"user,omitempty"`
	Shop      string    `json:"shop,omitempty"`
	Action    string    `json:"action,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type PageMetadata struct {
	Total  uint64
	Offset uint64
	Limit  uint64
	User   string
	Shop   string
	Action string
}

// PolicyPage contains a page of policies.
type PolicyPage struct {
	PageMetadata
	Policies []Policy
}

// PolicyRepository specifies an account persistence API.
type PolicyRepository interface {
	// Save creates a policy for the given User, Returns a non-nil
	// error in case of failures.
	Save(ctx context.Context, p Policy) error

	// Retrieve retrieves policy for a given input.
	Retrieve(ctx context.Context, pm PageMetadata) (PolicyPage, error)

	// Delete deletes the policy
	Delete(ctx context.Context, p Policy) error
}

// PolicyService represents a authorization service. It exposes
// functionalities through `auth` to perform authorization.
type PolicyService interface {
	// Issue issues a new Key, returning its token value alongside.
	Issue(ctx context.Context, req IssueToken) (string, error)

	// Identify validates token token. If token is valid, content
	// is returned. If token is invalid, or invocation failed for some
	// other reason, non-nil error value is returned in response.
	Identify(ctx context.Context, token string) (User, error)

	// Authorize checks authorization of the given `user`. Basically,
	// Authorize verifies that Is `user` allowed to `relation` on
	// `shop`. Authorize returns a non-nil error if the user has
	// no relation on the shop (which simply means the operation is
	// denied).
	Authorize(ctx context.Context, p Policy) error

	// AddPolicy creates a policy for the given user, so that, after
	// AddPolicy, `user` has a `relation` on `shop`. Returns a non-nil
	// error in case of failures.
	AddPolicy(ctx context.Context, token string, p Policy) error

	// ListPolicy lists policies based on the given policy structure.
	ListPolicy(ctx context.Context, token string, pm PageMetadata) (PolicyPage, error)

	// DeletePolicy removes a policy.
	DeletePolicy(ctx context.Context, token string, p Policy) error
}
