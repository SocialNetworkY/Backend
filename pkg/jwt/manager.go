package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Manager struct {
	secretKey        string
	tokenDuration    time.Duration
	refreshSecretKey string
	refreshDuration  time.Duration
}

func NewManager(accessTokenSecretKey string, accessTokenDuration time.Duration, refreshTokenSecretKey string, refreshTokenDuration time.Duration) *Manager {
	return &Manager{
		secretKey:        accessTokenSecretKey,
		tokenDuration:    accessTokenDuration,
		refreshSecretKey: refreshTokenSecretKey,
		refreshDuration:  refreshTokenDuration,
	}
}

type UserClaims struct {
	jwt.StandardClaims
	UserID uint `json:"user_id"`
}

func (manager *Manager) Generate(userID uint) (accessToken string, refreshToken string, err error) {
	// Generate access token
	accessToken, err = manager.generateToken(userID, manager.secretKey, manager.tokenDuration)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshToken, err = manager.generateToken(userID, manager.refreshSecretKey, manager.refreshDuration)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (manager *Manager) generateToken(userID uint, secretKey string, duration time.Duration) (string, error) {
	claims := &UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
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
				return nil, errors.New("unexpected token signing method")
			}
			return []byte(secretKey), nil
		},
	)
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token")
	}

	return claims.UserID, nil
}
