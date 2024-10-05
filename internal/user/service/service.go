package service

type Service struct {
	User *UserService
	Ban  *BanService
}

func New(userStorage UserStorage, banStorage BanStorage, is ImageStorage, ag AuthGateway) *Service {
	userService := NewUserService(userStorage, is, ag)
	banService := NewBanService(banStorage)

	return &Service{
		User: userService,
		Ban:  banService,
	}
}
