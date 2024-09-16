package service

type Service struct {
	User            *UserService
	Tokens          *TokensService
	Authentication  *AuthenticationService
	ActivationToken *ActivationTokenService
}

func New(userStorage UserStorage, refreshTokenStorage TokensRefreshTokenStorage, activationTokenStorage ActivationTokenStorage, tokenManager TokensManager, hasher Hasher) *Service {
	activationTokenService := NewActivationTokenService(activationTokenStorage)
	tokensService := NewTokensService(refreshTokenStorage, tokenManager)
	userService := NewUserService(userStorage, tokensService, activationTokenService, hasher)
	authenticationService := NewAuthenticationService(userService, tokensService)

	return &Service{
		User:            userService,
		Tokens:          tokensService,
		Authentication:  authenticationService,
		ActivationToken: activationTokenService,
	}
}
