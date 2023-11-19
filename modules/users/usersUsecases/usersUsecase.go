package usersUsecases

import (
	"github.com/NatthawutSK/ri-shop/config"
	"github.com/NatthawutSK/ri-shop/modules/users"
	"github.com/NatthawutSK/ri-shop/modules/users/usersRepositories"
)

type IUserUsecase interface{
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
}

type UserUsecase struct {
	cfg             config.IConfig
	usersRepository usersRepositories.IUsersRepository
	
}

func UserUsecaseHandler(usersRepository usersRepositories.IUsersRepository, cfg config.IConfig) IUserUsecase {
	return &UserUsecase{
		usersRepository: usersRepository,
		cfg: cfg,
	}
}

func (u *UserUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error) {
	//hashing password
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}
	//insert user
	result, err := u.usersRepository.InsertUser(req, false)
	if err != nil {
		return nil, err
	}
	return result, nil
}