package repo

type UserRepo interface {
	GetByUsername(username string) bool
}

type UserRepoImpl struct{}

func NewUserRepo() UserRepo {
	return UserRepoImpl{}
}

func (a UserRepoImpl) GetByUsername(username string) bool {
	registerUser := getRegisterUser()
	_, exists := registerUser[username]
	if !exists {
		return false
	}
	return true
}

func getRegisterUser() map[string]string {
	users := make(map[string]string)
	users["narik"] = "narik"
	users["santo"] = "santo"
	return users
}
