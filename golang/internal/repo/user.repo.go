package repo

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}
func (ur *UserRepo) GetInfoUser() string {
	return "Get Info User"
}
