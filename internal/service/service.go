package service

type Service struct {
	User           *UserService
	Token          *TokenService
	Authentication *AuthenticationService
}

func New(userStorage UserStorage, refreshTokenStorage RefreshTokenStorage, tokenManager TokenManager, hasher Hasher) *Service {
	tokenService := NewTokenService(refreshTokenStorage, tokenManager)
	userService := NewUserService(userStorage, tokenService, hasher)
	authenticationService := NewAuthenticationService(userService, tokenService)

	return &Service{
		User:           userService,
		Token:          tokenService,
		Authentication: authenticationService,
	}
}
