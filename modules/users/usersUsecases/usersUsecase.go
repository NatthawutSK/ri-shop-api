package usersUsecases

import (
	"fmt"

	"github.com/NatthawutSK/ri-shop/config"
	"github.com/NatthawutSK/ri-shop/modules/users"
	"github.com/NatthawutSK/ri-shop/modules/users/usersRepositories"
	riAuth "github.com/NatthawutSK/ri-shop/pkg/riauth"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface{
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
	GetPassport(req *users.UserCredential) (*users.UserPassport, error)
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

func (u *UserUsecase) GetPassport(req *users.UserCredential) (*users.UserPassport, error) {
	user, err := u.usersRepository.FindOneUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	// sign token
	accessToken, err := riAuth.NewRiAuth(riAuth.Access, u.cfg.Jwt(), &users.UserClaims{
		Id: user.Id,
		RoleId: user.RoleId,
	})
	refreshToken, err := riAuth.NewRiAuth(riAuth.Refresh, u.cfg.Jwt(), &users.UserClaims{
		Id: user.Id,
		RoleId: user.RoleId,
	})


	// set passport
	passport := &users.UserPassport{
		User: &users.User{
			Id: user.Id,
			Email: user.Email,
			Username: user.Username,
			RoleId: user.RoleId,
		},
		Token: &users.UserToken{
			AccessToken: accessToken.SignToken(),
			RefreshToken: refreshToken.SignToken(),
		},
	}

	if err := u.usersRepository.InsertOauth(passport); err != nil {
		return nil, err
	}
	return passport, nil

}