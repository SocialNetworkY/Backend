package service

type Service struct {
	User *UserService
	Ban  *BanService
}

func New(userStorage UserStorage, banStorage BanStorage, ag AuthGateway) *Service {
	userService := NewUserService(userStorage, ag)
	banService := NewBanService(banStorage)

	return &Service{
		User: userService,
		Ban:  banService,
	}
}
