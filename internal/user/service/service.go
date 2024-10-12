package service

type Service struct {
	User *UserService
	Ban  *BanService
}

func New(userRepo UserRepo, banRepo BanRepo, is ImageStorage, ag AuthGateway) *Service {
	userService := NewUserService(userRepo, is, ag)
	banService := NewBanService(banRepo)

	return &Service{
		User: userService,
		Ban:  banService,
	}
}
