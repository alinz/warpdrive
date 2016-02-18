package service

type Service struct {
	userService
}

func New() *Service {
	return &Service{}
}
