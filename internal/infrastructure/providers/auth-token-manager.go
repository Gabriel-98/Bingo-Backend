package providers

import (
	"errors"
	"fmt"
	"github.com/gabriel-98/bingo-backend/internal/application/types"
	"github.com/gabriel-98/bingo-backend/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

// An AuthTokenManager creates and validates access and refresh tokens.
type AuthTokenManager struct {
	refreshTokenDuration time.Duration
	refreshSigningKey string
	accessTokenDuration time.Duration
	accessSigningKey string
	issuer string
}

// NewAuthTokenManager creates a new AuthTokenManager using the configuration set
// for the token durations (access and refresh), the signing key (also for each token
// type), and the token issuer.
func NewAuthTokenManager(config config.AuthConfig) *AuthTokenManager {
	return &AuthTokenManager{
		refreshTokenDuration: config.RefreshToken.Duration,
		refreshSigningKey: config.RefreshToken.SigningKey,
		accessTokenDuration: config.AccessToken.Duration,
		accessSigningKey: config.AccessToken.SigningKey,
		issuer: config.Issuer,
	}
}

// newToken creates a new JWT token encapsulating an entities.UserAuthData and the claims
// 'exp', 'iat' and 'iss'.
// The token signature is generated using the signing method HS256.
func (s *AuthTokenManager) newToken(tokenDuration time.Duration, signingKey string, userAuthData types.UserAuthData) (string, error) {
	// Creates a jwt.Token containing the header and payload data.
	claims := AuthTokenPayload{
		UserAuthData: userAuthData,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenDuration)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer: s.issuer,
		},
	}
	jwttoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generates the token signature from the header and payload, and returns the token
	// string.
	token, err := jwttoken.SignedString([]byte(signingKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %w", err)
	}
	return token, nil
}

// validateToken validates a JWT token, which includes checking the existence of the
// claims 'exp' and 'iss', validating the signature, and that the token has not yet
// expired.
func (s *AuthTokenManager) validateToken(signingKey string, token string) (*types.UserAuthData, error) {
	// Converts the token string to a jwt.Token and validates it.
	jwttoken, err := jwt.ParseWithClaims(
		token,
		&AuthTokenPayload{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(signingKey), nil
		},
		jwt.WithExpirationRequired(),
		//jwt.WithIssuedAt(), // This method does not work.
		jwt.WithIssuer(s.issuer),
	)
	switch {
		case errors.Is(err, jwt.ErrInvalidKey):
			return nil, fmt.Errorf("token manager error: %w", err)
		case errors.Is(err, jwt.ErrInvalidKeyType):
			return nil, fmt.Errorf("token manager error: %w", err)
		case errors.Is(err, jwt.ErrHashUnavailable):
			return nil, fmt.Errorf("token manager error: %w", err)
		case errors.Is(err, jwt.ErrTokenExpired):
			// Identifies the user ID of the token owner, the token expiration time, and
			// the current time. Errors are ignored because the token's signature has
			// already been validated and the expiration time has been set.
			payload := AuthTokenPayload{}
			jwt.NewParser().ParseUnverified(token, &payload) // Returned data and error are ignored.
			userId := strconv.FormatInt(payload.UserAuthData.UserId, 10)

			expirationTimeND, _ := jwttoken.Claims.GetExpirationTime() // Error is ignored.
			expirationTime := expirationTimeND.Time
			return nil, fmt.Errorf("expired token error: user-id=%d, expiration-time=%s, current-time=%s",
				userId, expirationTime, time.Now())
		case err != nil:
			return nil, fmt.Errorf("invalid token: %w",err)
	}

	// Returns an entities.UserAuthData containing the token's custom fields.
	tokenPayload, ok := jwttoken.Claims.(*AuthTokenPayload)
	if !ok {
		return nil, fmt.Errorf(
			"failed to validate JWT token: invalid claims type: expected (%T), found (%T)",
			&AuthTokenPayload{},
			jwttoken.Claims,
		)
	}
	return &tokenPayload.UserAuthData, nil
}

// NewRefreshToken creates a new refresh token (Long-lived token).
// The returned token will contain the custom data provided by userAuthData and the
// following claims: Expiration time, Issued at, and Issuer.
// This method creates tokens with the duration and signing key defined for refresh
// tokens.
func (s *AuthTokenManager) NewRefreshToken(userAuthData types.UserAuthData) (string, error) {
	return s.newToken(s.refreshTokenDuration, s.refreshSigningKey, userAuthData)
}

// NewAccessToken creates a new access token (Short-lived token).
// The returned token will contain the custom data provided by userAuthData and the
// following claims: Expiration time, Issued at, and Issuer.
// This method creates tokens with the duration and signing key defined for access
// tokens.
func (s *AuthTokenManager) NewAccessToken(userAuthData types.UserAuthData) (string, error) {
	return s.newToken(s.accessTokenDuration, s.accessSigningKey, userAuthData)
}

// ValidateRefreshToken validates a JWT token, which includes checking the existence of
// the claims: Expiration time and Issuer, validating the signature, and that the token
// has not yet expired.
// This method uses the refresh signing key for validation.
func (s *AuthTokenManager) ValidateRefreshToken(tokenString string) (*types.UserAuthData, error) {
	return s.validateToken(s.refreshSigningKey, tokenString)
}

// ValidateAccessToken validates a JWT token, which includes checking the existence of
// the claims: Expiration time and Issuer, validating the signature, and that the token
// has not yet expired.
// This method uses the access signing key for validation.
func (s *AuthTokenManager) ValidateAccessToken(tokenString string) (*types.UserAuthData, error) {
	return s.validateToken(s.accessSigningKey, tokenString)
}

// IssuedAtAndExpiresAt retrieves the issue and expiration time of a JWT token.
// The token signature is not validated, therefore, a non-nil error is returned only
// because the token is malformed or the signing method is not available or specified.
func (s *AuthTokenManager) IssuedAtAndExpiresAt(token string) (time.Time, time.Time, error) {
	// Converts the token string to a jwt.Token.
	jwttoken, _, err := jwt.NewParser().ParseUnverified(token, &AuthTokenPayload{})
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	// Retrieves the token payload, which contains claims and custom fields.
	tokenPayload, ok := jwttoken.Claims.(*AuthTokenPayload)
	if !ok {
		return time.Time{}, time.Time{}, fmt.Errorf(
			"invalid claims type: expected (%T), found (%T)",
			&AuthTokenPayload{},
			jwttoken.Claims,
		)
	}

	// Retrieves the issue and expiration time.
	issuedAt, err := tokenPayload.GetIssuedAt()
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	expirationTime, err := tokenPayload.GetExpirationTime()
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return issuedAt.Time, expirationTime.Time, nil
}

// An AuthTokenPayload contains the data to be stored in the authentication tokens.
// Additionally, AuthTokenPayload implements the jwt.Claims interface.
type AuthTokenPayload struct {
	// UserAuthData contains the custom fields of the token.
	types.UserAuthData

	// RegisteredClaims implements jwt.Claims, and its claims are described below:
	//
	// - 'exp': Expiration time
	//   UNIX timestamp indicating the date and time after which the token will become
	//   invalid.
	//   Access function: GetExpirationTime()
	//
	// - 'iat': Issued at
	//   UNIX timestamp indicating the date and time when the token was issued.
	//   Access function: GetIssuedAt()
	//
	// - 'nbf': Not before
	//   UNIX timestamp indicating the date and time after which the token will be valid.
	//   Access function: GetNotBefore()
	//
	// - 'iss': Issuer
	//   String indicating the issuer of the token (e.g., the application or service).
	//   Access function: GetIssuer()
	//
	// - 'sub': Subject
	//   String indicating the token recipient (e.g., user ID of the token owner).
	//   Access function: GetSubject()
	//
	// - 'aud': Audience
	//   One or more strings indicating the applications or services for which the token
	//   is intended.
	//   Access function: GetAudience()
	jwt.RegisteredClaims
}