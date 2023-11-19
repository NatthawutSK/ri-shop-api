package usersRepositories

import (
	"github.com/NatthawutSK/ri-shop/modules/users"
	"github.com/NatthawutSK/ri-shop/modules/users/usersPatterns"
	"github.com/jmoiron/sqlx"
)

type IUsersRepository interface{
	InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error)
}

type UsersRepository struct {
	db *sqlx.DB
}


func UsersRepositoryHandler(db *sqlx.DB) IUsersRepository {
	return &UsersRepository{
		db: db,
	}
}

func (r *UsersRepository) InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error) {
	result := usersPatterns.InsertUser(r.db, req, isAdmin)

	var err error
	if isAdmin {
		result, err = result.Admin()
		if err != nil {
			return nil, err
		}
	} else {
		result, err = result.Customer()
		if err != nil {
			return nil, err
		}
	}

	user ,err := result.Result()
	if err != nil {
		return nil, err
	}
	return user, nil
}