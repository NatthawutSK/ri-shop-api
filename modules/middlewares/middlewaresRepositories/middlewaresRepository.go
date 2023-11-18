package middlewaresRepositories

import "github.com/jmoiron/sqlx"

type IMiddlewaresRepository interface {
}

type middlewaresRepositories struct {
	db *sqlx.DB

}

func MiddlewaresRepository(db *sqlx.DB) IMiddlewaresRepository {
	return &middlewaresRepositories{
		db: db,
	}
}