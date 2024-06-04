package domain

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/victorspringer/backend-coding-challenge/lib/log"
)

// Client implements the Authenticator interface.
type Client struct {
	logger                 *log.Logger
	userServiceClient      UserServiceClient
	refreshTokenRepository Repository
	accessTokenRepository  Repository
	flowRepository         Repository
	issuer                 string
	expiration             map[string]time.Duration
	jwtKey                 *rsa.PrivateKey
}

// Claims represents JWT claims data structure.
type Claims struct {
	Name  string `json:"name,omitempty"`
	Level Level  `json:"level"`
	*jwt.RegisteredClaims
}

// Level represents level field type for claims data structure.
type Level string

// Valid Level values as constants.
const (
	AdminLevel     Level = "admin"
	UserLevel      Level = "user"
	AnonymousLevel Level = "anonymous"
)

// Expiration map valid keys.
const (
	AccessTokenExpiration       string = "accessTokenExpiration"
	AnonymousExpiration         string = "anonymousExpiration"
	ShortRefreshTokenExpiration string = "shortRefreshTokenExpiration"
	LongRefreshTokenExpiration  string = "longRefreshTokenExpiration"
)

// Common errors.
var (
	ErrInternalError = errors.New("internal error")
	ErrUnauthorized  = errors.New("unauthorized")
)

// NewClientSingleton returns a new instance of authentication Client.
func NewClient(
	logger *log.Logger,
	userServiceClient UserServiceClient,
	refreshTokenRepository,
	accessTokenRepository,
	flowRepository Repository,
	issuer string,
	expiration map[string]time.Duration,
	jwtKey *rsa.PrivateKey,
) *Client {
	return &Client{
		logger:                 logger,
		userServiceClient:      userServiceClient,
		refreshTokenRepository: refreshTokenRepository,
		accessTokenRepository:  accessTokenRepository,
		flowRepository:         flowRepository,
		issuer:                 issuer,
		expiration:             expiration,
		jwtKey:                 jwtKey,
	}
}

func (c *Client) GenerateAnonymousTokens(userID string, flow FlowType) (*Tokens, error) {
	claims := &Claims{
		Level: AdminLevel,
		RegisteredClaims: &jwt.RegisteredClaims{
			Subject: userID,
		},
	}
	return c.generateTokens(claims, flow)
}

func (c *Client) GenerateUserTokens(username, password string, flow FlowType) (*Tokens, error) {
	user, err := c.userServiceClient.CheckCredentials(username, password)
	if err != nil {
		if err == errUserNotFound || err == errBadRequest {
			c.logger.Info("user service: not found or bad credentials", log.String("username", username), log.Error(err))
			return nil, ErrUnauthorized
		}
		c.logger.Error("user service returned error", log.String("username", username), log.Error(err))
		return nil, ErrInternalError
	}

	claims := &Claims{
		Level: Level(user.Level),
		Name:  user.Name,
		RegisteredClaims: &jwt.RegisteredClaims{
			Subject: username,
		},
	}

	return c.generateTokens(claims, flow)
}

// generateTokens generates authentication tokens.
func (c *Client) generateTokens(claims *Claims, flow FlowType) (*Tokens, error) {
	if claims.RegisteredClaims == nil {
		claims.RegisteredClaims = &jwt.RegisteredClaims{}
	}
	claims.Issuer = c.issuer
	claims.IssuedAt = jwt.NewNumericDate(time.Now())

	claims.ExpiresAt = jwt.NewNumericDate(claims.IssuedAt.Add(c.expiration[AccessTokenExpiration]))
	expiresIn := c.expiration[AccessTokenExpiration]
	if claims.Level == AnonymousLevel {
		claims.ExpiresAt = jwt.NewNumericDate(claims.IssuedAt.Add(c.expiration[AnonymousExpiration]))
		expiresIn = c.expiration[AnonymousExpiration]
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	accessToken, err := token.SignedString(c.jwtKey)
	if err != nil {
		c.logger.Error("failed to generate access token", log.Error(err))
		return nil, ErrInternalError
	}

	if claims.Level == AnonymousLevel {
		return &Tokens{
			AccessToken: accessToken,
			ExpiresIn:   int64(expiresIn.Seconds()),
		}, nil
	}

	refreshToken := uuid.New().String()
	var repositoryExpiration time.Duration
	if flow == RememberMeFlow {
		repositoryExpiration = time.Duration(c.expiration[LongRefreshTokenExpiration] * 1e9)
	} else {
		repositoryExpiration = time.Duration(c.expiration[ShortRefreshTokenExpiration] * 1e9)
	}

	ctx := context.Background()

	err = c.accessTokenRepository.Set(ctx, refreshToken, accessToken, repositoryExpiration)
	if err != nil {
		c.logger.Error(
			"failed to save access token",
			log.String("key", refreshToken),
			log.String("access_token", accessToken),
			log.Error(err),
		)
		return nil, ErrInternalError
	}

	refreshTokenKey := fmt.Sprintf("%s-%d", claims.Subject, claims.IssuedAt.Unix())
	err = c.refreshTokenRepository.Set(ctx, refreshTokenKey, refreshToken, repositoryExpiration)
	if err != nil {
		c.logger.Error(
			"failed to save refresh token",
			log.String("key", refreshTokenKey),
			log.String("access_token", refreshToken),
			log.Error(err),
		)
		err := c.accessTokenRepository.Del(ctx, refreshToken) // cleanup prev step
		if err != nil {
			c.logger.Warn("failed to delete from access token repository", log.String("key", refreshToken), log.Error(err))
		}
		return nil, ErrInternalError
	}

	err = c.flowRepository.Set(ctx, refreshTokenKey, string(flow), repositoryExpiration)
	if err != nil {
		c.logger.Error(
			"failed to save flow",
			log.String("key", refreshTokenKey),
			log.String("flow", string(flow)),
			log.Error(err),
		)
		err := c.accessTokenRepository.Del(ctx, refreshToken) // cleanup prev step
		if err != nil {
			c.logger.Warn("failed to delete from access token repository", log.String("key", refreshToken), log.Error(err))
		}
		err = c.refreshTokenRepository.Del(ctx, refreshTokenKey) // cleanup prev step
		if err != nil {
			c.logger.Warn("failed to delete from refresh token repository", log.String("key", refreshTokenKey), log.Error(err))
		}
		return nil, ErrInternalError
	}

	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(expiresIn.Seconds()),
	}, nil
}

// Revoke revokes an access token registered for a given refresh token.
func (c *Client) Revoke(accessToken string) error {
	claims, err := c.decryptAccessToken(accessToken)
	if err != nil {
		c.logger.Warn("failed to decrypt access token", log.String("access_token", accessToken), log.Error(err))
		return ErrUnauthorized
	}

	ctx := context.Background()

	refreshTokenKey := fmt.Sprintf("%s-%d", claims.Subject, claims.IssuedAt.Unix())
	refreshToken, err := c.refreshTokenRepository.Get(ctx, refreshTokenKey)
	if err != nil {
		c.logger.Debug("failed to get from refresh token repository", log.String("key", refreshTokenKey), log.Error(err))
		return ErrUnauthorized
	}
	err = c.accessTokenRepository.Del(ctx, refreshToken)
	if err != nil {
		c.logger.Warn("failed to delete from access token repository", log.String("key", refreshToken), log.Error(err))
	}
	err = c.refreshTokenRepository.Del(ctx, refreshTokenKey)
	if err != nil {
		c.logger.Warn("failed to delete from refresh token repository", log.String("key", refreshTokenKey), log.Error(err))
	}
	err = c.flowRepository.Del(ctx, refreshTokenKey)
	if err != nil {
		c.logger.Warn("failed to delete from flow repository", log.String("key", refreshTokenKey), log.Error(err))
	}

	return nil
}

// Refresh refreshes logged-in/anonymous user authentication token.
func (c *Client) Refresh(refreshToken string) (*Tokens, error) {
	ctx := context.Background()

	accessToken, err := c.accessTokenRepository.Get(ctx, refreshToken)
	if err != nil {
		c.logger.Debug("failed to get from access token repository", log.String("key", refreshToken), log.Error(err))
		return nil, ErrUnauthorized
	}

	claims, err := c.decryptAccessToken(accessToken)
	if err != nil {
		c.logger.Warn("failed to decrypt access token", log.String("access_token", accessToken), log.Error(err))
		return nil, ErrUnauthorized
	}

	refreshTokenKey := fmt.Sprintf("%s-%d", claims.Subject, claims.IssuedAt.Unix())
	rt, err := c.refreshTokenRepository.Get(ctx, refreshTokenKey)
	if err != nil {
		c.logger.Debug("failed to get from access token repository", log.String("key", refreshTokenKey), log.Error(err))
		return nil, ErrUnauthorized
	}

	if rt != refreshToken {
		c.logger.Debug("tokens do not match", log.String("rt", rt), log.String("refresh_token", refreshTokenKey))
		return nil, ErrUnauthorized
	}

	flow, err := c.flowRepository.Get(ctx, refreshTokenKey)
	if err != nil {
		c.logger.Debug("failed to get from flow repository", log.String("key", refreshTokenKey), log.Error(err))
		return nil, ErrUnauthorized
	}

	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.ExpiresAt = jwt.NewNumericDate(claims.IssuedAt.Add(c.expiration[AccessTokenExpiration]))
	expiresIn := c.expiration[AccessTokenExpiration]
	newRefreshToken := uuid.New().String()
	var repositoryExpiration time.Duration
	if (FlowType(flow)) == RememberMeFlow {
		repositoryExpiration = time.Duration(c.expiration[LongRefreshTokenExpiration] * 1e9)
	} else {
		repositoryExpiration = time.Duration(c.expiration[ShortRefreshTokenExpiration] * 1e9)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	newAccessToken, err := token.SignedString(c.jwtKey)
	if err != nil {
		c.logger.Error("failed to generate access token", log.Error(err))
		return nil, ErrInternalError
	}

	//save new refresh token
	err = c.accessTokenRepository.Set(ctx, newRefreshToken, newAccessToken, repositoryExpiration)
	if err != nil {
		c.logger.Error(
			"failed to save to access token repository",
			log.String("key", refreshToken),
			log.String("access_token", accessToken),
			log.Error(err),
		)
		return nil, ErrInternalError
	}

	//del prev refresh token expiration
	err = c.refreshTokenRepository.Del(ctx, refreshTokenKey)
	if err != nil {
		c.logger.Warn("failed to delete from refresh token repository", log.String("key", refreshTokenKey), log.Error(err))
	}
	//del prev flow
	err = c.flowRepository.Del(ctx, refreshTokenKey)
	if err != nil {
		c.logger.Warn("failed to delete from flow repository", log.String("key", refreshTokenKey), log.Error(err))
	}

	refreshTokenKey = fmt.Sprintf("%s-%d", claims.Subject, claims.IssuedAt.Unix())
	//save new refresh token expirtion
	err = c.refreshTokenRepository.Set(ctx, refreshTokenKey, newRefreshToken, repositoryExpiration)
	if err != nil {
		c.logger.Error(
			"failed to save to access token repository",
			log.String("key", refreshToken),
			log.String("access_token", accessToken),
			log.Error(err),
		)
		err := c.accessTokenRepository.Del(ctx, newRefreshToken) //cleanup prev step
		if err != nil {
			c.logger.Warn("failed to delete from access token repository", log.String("key", newRefreshToken), log.Error(err))
		}
		return nil, ErrInternalError
	}

	//del prev access token
	err = c.accessTokenRepository.Del(ctx, refreshToken)
	if err != nil {
		c.logger.Warn("failed to delete from access token repository", log.String("key", refreshToken), log.Error(err))
	}

	//save new flow
	err = c.flowRepository.Set(ctx, refreshTokenKey, flow, repositoryExpiration)
	if err != nil {
		c.logger.Error("failed to save to flow repository", log.String("key", refreshTokenKey), log.String("flow", flow), log.Error(err))
		err := c.accessTokenRepository.Del(ctx, newRefreshToken) //cleanup prev step
		if err != nil {
			c.logger.Warn("failed to delete from access token repository", log.String("key", refreshToken), log.Error(err))
		}
		err = c.refreshTokenRepository.Del(ctx, refreshTokenKey) //cleanup prev step
		if err != nil {
			c.logger.Warn("failed to delete from refresh token repository", log.String("key", refreshTokenKey), log.Error(err))
		}
		return nil, ErrInternalError
	}

	c.logger.Debug(fmt.Sprintf("token %s refreshed to %s", refreshToken, newRefreshToken))

	return &Tokens{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(expiresIn.Seconds()),
	}, nil
}

// ValidateAccessToken checks if logged-in user authentication token exists.
func (c *Client) ValidateAccessToken(accessToken string) error {
	claims, err := c.decryptAccessToken(accessToken)
	if err != nil {
		c.logger.Warn("failed to decrypt access token", log.String("access_token", accessToken), log.Error(err))
		return ErrUnauthorized
	}

	ctx := context.Background()

	refreshTokenKey := fmt.Sprintf("%s-%d", claims.Subject, claims.IssuedAt.Unix())
	refreshToken, err := c.refreshTokenRepository.Get(ctx, refreshTokenKey)
	if err != nil {
		c.logger.Debug("failed to get from refresh token repository", log.String("key", refreshTokenKey), log.Error(err))
		return ErrUnauthorized
	}

	_, err = c.accessTokenRepository.Get(ctx, refreshToken)
	if err != nil {
		c.logger.Debug("failed to get from access token repository", log.String("key", refreshToken), log.Error(err))
		return ErrUnauthorized
	}

	return nil
}

func (c *Client) decryptAccessToken(accessToken string) (*Claims, error) {
	claims := &Claims{}

	parser := jwt.NewParser(jwt.WithoutClaimsValidation())

	_, err := parser.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return &c.jwtKey.PublicKey, nil
	})
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (c *Client) JWTKey() *rsa.PrivateKey {
	return c.jwtKey
}
