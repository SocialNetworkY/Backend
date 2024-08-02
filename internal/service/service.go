package service

type Service struct {
	User  *UserService
	Token *TokenService
}

func New(userStorage UserStorage, refreshTokenStorage RefreshTokenStorage) *Service {
	tokenService := NewTokenService(refreshTokenStorage)
	userService := NewUserService(userStorage, tokenService)

	return &Service{
		User:  userService,
		Token: tokenService,
	}
}
