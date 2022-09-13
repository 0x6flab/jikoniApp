package auth

import (
	context "context"
	"time"

	"github.com/0x6flab/jikoniApp/BackendApp/internal/errors"
	"github.com/golang-jwt/jwt/v4"
)

var secret = []byte("?~m:E.wT&3|o:}(uti=~(:L@C/M.C")

var _ PolicyService = (*service)(nil)

type claims struct {
	jwt.StandardClaims
	PhoneNumber *string `json:"phone_number,omitempty"`
	Type        *string `json:"type,omitempty"`
}

type service struct {
	policies      PolicyRepository
	loginDuration time.Duration
}

// New instantiates the auth service implementation.
func New(policies PolicyRepository, duration time.Duration) PolicyService {
	return &service{
		policies:      policies,
		loginDuration: duration,
	}
}

func (svc service) Issue(ctx context.Context, req IssueToken) (string, error) {
	claims := claims{
		StandardClaims: jwt.StandardClaims{
			Subject:   req.PhoneNumber,
			Issuer:    "jikoni-auth",
			IssuedAt:  time.Now().Local().Unix(),
			ExpiresAt: time.Now().Add(svc.loginDuration).Local().Unix(),
		},
		Type: &req.Type,
	}
	if req.ID != "" {
		claims.Id = req.ID
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func (svc service) Identify(ctx context.Context, token string) (User, error) {
	c := claims{}
	_, _ = jwt.ParseWithClaims(token, &c, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrAuthentication
		}
		return []byte(secret), nil
	})

	return c.toUser(), nil
}

func (svc service) Authorize(ctx context.Context, p Policy) error {
	panic("Unimplemented")
}

func (svc service) AddPolicy(ctx context.Context, token string, p Policy) error {
	panic("Unimplemented")
}

func (svc service) ListPolicy(ctx context.Context, token string, pm PageMetadata) (PolicyPage, error) {

	panic("Unimplemented")
}

func (svc service) DeletePolicy(ctx context.Context, token string, p Policy) error {
	panic("Unimplemented")
}

func (c claims) toUser() User {
	user := User{
		ID:          c.Id,
		PhoneNumber: c.Subject,
		Type:        *c.Type,
		CreatedAt:   time.Unix(c.IssuedAt, 0).UTC(),
	}

	return user
}
