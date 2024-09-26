package service

type Service struct {
	User *UserService
}

func New(userStorage UserStorage, ag AuthGateway) *Service {
	userService := NewUserService(userStorage, ag)

	return &Service{
		User: userService,
	}
}
