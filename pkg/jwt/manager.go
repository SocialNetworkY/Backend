package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Manager struct {
	secretKey        string
	tokenDuration    time.Duration
	refreshSecretKey string
	refreshDuration  time.Duration
}

func NewManager(tokenDuration, refreshDuration time.Duration, secret, refreshSecret string) *Manager {
	return &Manager{
		secretKey:        secret,
		tokenDuration:    tokenDuration,
		refreshSecretKey: refreshSecret,
		refreshDuration:  refreshDuration,
	}
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID uint `json:"user_id"`
}

func (manager *Manager) Generate(userID uint) (accessToken string, refreshToken string, err error) {
	// Generate access token
	accessToken, err = manager.generateToken(userID, manager.secretKey, manager.tokenDuration)
	if err != nil {
		return "", "", ErrGenerateAccessToken
	}

	// Generate refresh token
	refreshToken, err = manager.generateToken(userID, manager.refreshSecretKey, manager.refreshDuration)
	if err != nil {
		return "", "", ErrGenerateRefreshToken
	}

	return accessToken, refreshToken, nil
}

func (manager *Manager) generateToken(userID uint, secretKey string, duration time.Duration) (string, error) {
	claims := &UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func (manager *Manager) Verify(accessToken string) (userID uint, err error) {
	return manager.verifyToken(accessToken, manager.secretKey)
}

func (manager *Manager) VerifyRefreshToken(refreshToken string) (userID uint, err error) {
	return manager.verifyToken(refreshToken, manager.refreshSecretKey)
}

func (manager *Manager) verifyToken(tokenString, secretKey string) (userID uint, err error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrUnexpectedSigningMethod
			}
			return []byte(secretKey), nil
		},
	)
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return 0, ErrInvalidToken
	}

	return claims.UserID, nil
}
