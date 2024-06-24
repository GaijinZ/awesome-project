package repositories

import (
	"awesomeProject/internal/models"
	"awesomeProject/pkg/database"
)

type AuthRepository interface {
	Register(user *models.Auth) error
	Login(auth *models.Auth) (*models.UserResponse, error)
}

type AuthRepositoryImpl struct {
	db database.Database
}

func NewAuthRepositoryImpl(db database.Database) AuthRepository {
	return &AuthRepositoryImpl{
		db: db,
	}
}

func (a *AuthRepositoryImpl) Register(user *models.Auth) error {
	_, err := a.db.Exec(AddCustomer, &user.Username, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthRepositoryImpl) Login(auth *models.Auth) (*models.UserResponse, error) {
	var user models.UserResponse

	err := a.db.QueryRow(GetUserByEmail, auth.Email).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return &user, err
	}

	return &user, nil
}
