package service

type Service struct {
	User  *UserService
	Token *TokenService
}

func New(userStorage UserStorage, refreshTokenStorage RefreshTokenStorage, tokenManager TokenManager, hasher Hasher) *Service {
	tokenService := NewTokenService(refreshTokenStorage, tokenManager)
	userService := NewUserService(userStorage, tokenService, hasher)

	return &Service{
		User:  userService,
		Token: tokenService,
	}
}
