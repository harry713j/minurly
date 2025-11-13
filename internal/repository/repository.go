package repository

type Repository struct {
	Auth *AuthRepo
	User *UserRepo
	Url  *UrlRepo
}

func NewRepository(auth *AuthRepo, user *UserRepo, url *UrlRepo) *Repository {
	return &Repository{
		Auth: auth,
		User: user,
		Url:  url,
	}
}
