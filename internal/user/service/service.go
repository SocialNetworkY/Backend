package service

type Service struct {
	User *UserService
	Ban  *BanService
}

func New(userRepo UserRepo, banRepo BanRepo, is ImageStorage, ag AuthGateway, pg PostGateway) *Service {
	userService := NewUserService(userRepo, is, ag, pg)
	banService := NewBanService(banRepo, pg)

	return &Service{
		User: userService,
		Ban:  banService,
	}
}
