package users

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       string `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Username string `db:"username" json:"username"`
	RoleId   int `db:"role_id" json:"role_id"`
}

type UserRegisterReq struct {
	Email    string `db:"email" json:"email" form:"email"`
	Password string `db:"password" json:"password" form:"password"`
	Username string `db:"username" json:"username" form:"username"`
}

type UserCredentialCheck struct {
	Id 	 string `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
	Username string `db:"username" json:"username"`
	RoleId   int `db:"role_id" json:"role_id"`
}

type UserCredential struct {
	Email string `db:"email" json:"email" form:"email"`
	Password string `db:"password" json:"password" form:"password"`
}

func (obj *UserRegisterReq) BcryptHashing() error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(obj.Password), 10)
	if err != nil {
		return fmt.Errorf("hash password failed: %v", err)
	}
	obj.Password = string(hashPassword)
	return nil
}

func (obj *UserRegisterReq) IsEmail() bool {
	match, err := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, obj.Email)
	if err != nil {
		return false
	}
	return match
	
}


type UserPassport struct {
	User *User `json:"user"`
	Token *UserToken `json:"token"`
}

type UserToken struct {
	Id string `db:"id" json:"id"`
	AccessToken string `db:"access_token" json:"access_token"`
	RefreshToken string `db:"refresh_token" json:"refresh_token"`
}


type UserClaims struct {
	Id string `json:"id" db:"id"`
	RoleId int `json:"role" db:"role"`
}