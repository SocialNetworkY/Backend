package service

type Service struct {
	User *UserService
	Ban  *BanService
}

func New(userRepo UserRepo, banRepo BanRepo, is ImageStorage, ag AuthGateway, pg PostGateway, rg ReportGateway) *Service {
	userService := NewUserService(userRepo, is, ag, pg, rg)
	banService := NewBanService(banRepo, pg)

	return &Service{
		User: userService,
		Ban:  banService,
	}
}
