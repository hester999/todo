package repo

import (
	"github.com/jmoiron/sqlx"
)

type UserRepoImpl struct {
	db *sqlx.DB
}

func NewUserRepoImpl(db *sqlx.DB) *UserRepoImpl {
	return &UserRepoImpl{
		db: db,
	}
}

//func (u *UserRepoImpl) CreateUser(user entity.User) (entity.User, error) {
//	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3)`
//	newUser := struct {
//		Id   string `db:"id"`
//		Name string `db:"name"`
//		TgId int64  `db:"messanger_id"`
//	}(user)
//
//}
