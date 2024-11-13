package user

import (
	"go/adv-demo/pkg/db"
)

type UserRepository struct {
	Database *db.Db
}

func NewUserRepository(db *db.Db) *UserRepository{
	return &UserRepository{
		Database: db,
	}
}

func (repo *UserRepository) Create(user *User) (*User, error){
	result := repo.Database.DB.Create(user)
	if result.Error != nil{
		return nil, result.Error
	}
	return user, nil
}

func (repo *UserRepository) FindByEmail(email string) (*User, error){
	var user User
	result := repo.Database.DB.First(&user, "email = ?", email)
	if result.Error != nil{
		return nil, result.Error
	}
	return &user, nil
}