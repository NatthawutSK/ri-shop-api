package usersUsecases

import (
	"fmt"

	"github.com/NatthawutSK/ri-shop/config"
	"github.com/NatthawutSK/ri-shop/modules/users"
	"github.com/NatthawutSK/ri-shop/modules/users/usersRepositories"
	riAuth "github.com/NatthawutSK/ri-shop/pkg/riauth"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
	InsertAdmin(req *users.UserRegisterReq) (*users.UserPassport, error)
	GetPassport(req *users.UserCredential) (*users.UserPassport, error)
	RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error)
	DeleteOauth(oauthId string) error
	GetUserProfile(userId string) (*users.User, error)
}

type UserUsecase struct {
	cfg             config.IConfig
	usersRepository usersRepositories.IUsersRepository
}

func UserUsecaseHandler(usersRepository usersRepositories.IUsersRepository, cfg config.IConfig) IUserUsecase {
	return &UserUsecase{
		usersRepository: usersRepository,
		cfg:             cfg,
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

func (u *UserUsecase) InsertAdmin(req *users.UserRegisterReq) (*users.UserPassport, error) {
	//hashing password
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}
	//insert user
	result, err := u.usersRepository.InsertUser(req, true)
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
	accessToken, err1 := riAuth.NewRiAuth(riAuth.Access, u.cfg.Jwt(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})
	if err1 != nil {
		return nil, err
	}
	refreshToken, err2 := riAuth.NewRiAuth(riAuth.Refresh, u.cfg.Jwt(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})
	if err2 != nil {
		return nil, err
	}

	// set passport
	passport := &users.UserPassport{
		User: &users.User{
			Id:       user.Id,
			Email:    user.Email,
			Username: user.Username,
			RoleId:   user.RoleId,
		},
		Token: &users.UserToken{
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken.SignToken(),
		},
	}

	if err := u.usersRepository.InsertOauth(passport); err != nil {
		return nil, err
	}
	return passport, nil

}

// use for refresh token
func (u *UserUsecase) RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error) {
	claims, err := riAuth.ParseToken(u.cfg.Jwt(), req.RefreshToken)
	if err != nil {
		return nil, err
	}

	//check oauth
	oauth, err := u.usersRepository.FindOneOauth(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	//find profile
	profile, err := u.usersRepository.GetProfile(oauth.UserId)
	if err != nil {
		return nil, err
	}

	newClaims := &users.UserClaims{
		Id:     profile.Id,
		RoleId: profile.RoleId,
	}

	accessToken, err := riAuth.NewRiAuth(riAuth.Access, u.cfg.Jwt(), newClaims)
	if err != nil {
		return nil, err
	}
	refreshToken := riAuth.RepeatToken(u.cfg.Jwt(), newClaims, claims.ExpiresAt.Unix())

	passport := &users.UserPassport{
		User: profile,
		Token: &users.UserToken{
			Id:           oauth.Id,
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken,
		},
	}

	if err := u.usersRepository.UpdateOauth(passport.Token); err != nil {
		return nil, err
	}

	return passport, nil

}

func (u *UserUsecase) DeleteOauth(oauthId string) error {
	if err := u.usersRepository.DeleteOauth(oauthId); err != nil {
		return err
	}
	return nil

}

func (u *UserUsecase) GetUserProfile(userId string) (*users.User, error) {
	profile, err := u.usersRepository.GetProfile(userId)
	if err != nil {
		return nil, err
	}
	return profile, nil

}
