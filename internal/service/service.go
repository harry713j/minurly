package service

type Service struct {
	Auth *AuthService
	User *UserService
	Url  *UrlService
}

func NewService(auth *AuthService, user *UserService, url *UrlService) *Service {
	return &Service{
		Auth: auth,
		User: user,
		Url:  url,
	}
}
