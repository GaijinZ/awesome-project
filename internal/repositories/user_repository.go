package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"awesomeProject/internal/models"
	"awesomeProject/pkg/database"
)

type UserRepository interface {
	GetUserByUsername(id string) (*models.UserResponse, error)
	GetAllUsers() ([]models.UserResponse, error)
	UpdateUser(user *models.User) error
	CreateUser(user *models.User) error
	DeleteUser(id string) error
}

type UserRepositoryImpl struct {
	db database.Database
}

func NewUserRepository(db database.Database) UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (u *UserRepositoryImpl) GetUserByUsername(name string) (*models.UserResponse, error) {
	var userResponse models.UserResponse

	err := u.db.QueryRow(GetUserByUsername, name).Scan(&userResponse.Username, &userResponse.Email, &userResponse.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &userResponse, nil
}

func (u *UserRepositoryImpl) GetAllUsers() ([]models.UserResponse, error) {
	var users []models.UserResponse

	rows, err := u.db.Query(GetAllUsers)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.UserResponse

		err = rows.Scan(&user.Email, &user.Role)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (u *UserRepositoryImpl) UpdateUser(user *models.User) error {
	_, err := u.db.Exec(UpdateUser, user.ID, user.Username, user.Email, user.Role)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (u *UserRepositoryImpl) CreateUser(user *models.User) error {
	exists, err := checkUserExists(user.Email, u.db)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	if exists {
		return fmt.Errorf("user already exists")
	}

	_, err = u.db.Exec(AddCustomer, user.Username, user.Email, user.Password, user.Role)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (u *UserRepositoryImpl) DeleteUser(id string) error {
	_, err := u.db.Exec(DeleteUser, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func checkUserExists(userEmail string, db database.Database) (bool, error) {
	var exists bool

	err := db.QueryRow(CheckUserExists, userEmail).Scan(&exists)
	if err != nil {
		return exists, fmt.Errorf("error checking if user exists: %w", err)
	}

	return exists, nil
}

func getUserByEmail(email string, db database.Database) (*models.UserResponse, error) {
	var userResponse models.UserResponse

	return &userResponse, nil
}
