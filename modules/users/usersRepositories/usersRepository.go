package usersRepositories

import (
	"context"
	"fmt"
	"time"

	"github.com/NatthawutSK/ri-shop/modules/users"
	"github.com/NatthawutSK/ri-shop/modules/users/usersPatterns"
	"github.com/jmoiron/sqlx"
)

type IUsersRepository interface{
	InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error)
	FindOneUserByEmail(email string) (*users.UserCredentialCheck, error)
	InsertOauth(req *users.UserPassport) error
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


func (r *UsersRepository) FindOneUserByEmail(email string) (*users.UserCredentialCheck, error) {
	query := `
	SELECT
		"id",
		"email",
		"password",
		"username",
		"role_id"
	FROM "users"
	WHERE "email" = $1;`
	user := new(users.UserCredentialCheck)
	if err := r.db.Get(user, query, email); err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}


func (r *UsersRepository) InsertOauth(req *users.UserPassport) error{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
	INSERT INTO "oauth" (
		"user_id",
		"refresh_token",
		"access_token"
	)
	VALUES ($1, $2, $3)
		RETURNING "id";`

	if err := r.db.QueryRowContext(
		ctx,
		query,
		req.User.Id,
		req.Token.RefreshToken,
		req.Token.AccessToken,
	).Scan(&req.Token.Id); err != nil {
		return fmt.Errorf("insert oauth failed: %v", err)
	}
	return nil
}