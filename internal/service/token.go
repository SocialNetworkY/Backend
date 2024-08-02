package service

type (
	RefreshTokenStorage interface {
		SetToken(userID uint, refreshToken string) error
	}

	TokenManager interface {
		Generate(userID uint) (accessToken string, refreshToken string, err error)
		Verify(accessToken string) (userID uint, err error)
		VerifyRefreshToken(refreshToken string) (userID uint, err error)
	}
)

type TokenService struct {
	storage      RefreshTokenStorage
	tokenManager TokenManager
}

// TODO: Implement functions Verify
func NewTokenService(refreshTokenStorage RefreshTokenStorage, tokenManager TokenManager) *TokenService {
	return &TokenService{
		storage:      refreshTokenStorage,
		tokenManager: tokenManager,
	}
}

func (ts *TokenService) Generate(userID uint) (string, string, error) {
	accessToken, refreshToken, err := ts.tokenManager.Generate(userID)
	if err != nil {
		return "", "", err
	}

	if err := ts.storage.SetToken(userID, accessToken); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
