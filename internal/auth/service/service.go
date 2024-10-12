package service

type Service struct {
	User            *UserService
	Tokens          *TokensService
	Authentication  *AuthenticationService
	ActivationToken *ActivationTokenService
}

func New(userRepo UserRepo, refreshTokenRepo TokensRefreshTokenRepo, activationTokenRepo ActivationTokenRepo, tokenManager TokensManager, hasher Hasher, ug UserGateway) *Service {
	activationTokenService := NewActivationTokenService(activationTokenRepo)
	tokensService := NewTokensService(refreshTokenRepo, tokenManager)
	userService := NewUserService(userRepo, tokensService, activationTokenService, hasher, ug)
	authenticationService := NewAuthenticationService(userService, tokensService)

	return &Service{
		User:            userService,
		Tokens:          tokensService,
		Authentication:  authenticationService,
		ActivationToken: activationTokenService,
	}
}
