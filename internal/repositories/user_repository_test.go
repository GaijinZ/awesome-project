package repositories

import (
	"awesomeProject/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("User Repository", func() {
	var (
		db           *sql.DB
		mock         sqlmock.Sqlmock
		repo         UserRepository
		user         *models.User
		userResponse *models.UserResponse
		users        []models.UserResponse
		err          error
	)

	BeforeEach(func() {
		db, mock, err = sqlmock.New()
		Expect(err).Should(BeNil())

		repo = NewUserRepository(db)
		user = &models.User{
			ID:       "1",
			Username: "test user",
			Email:    "testuser@example.com",
			Password: "password",
			Role:     "tester",
		}

		userResponse = &models.UserResponse{}
		users = []models.UserResponse{}
	})

	AfterEach(func() {
		db.Close()
	})

	Describe("GetUserByUsername", func() {
		It("should return user response successfully", func() {
			rows := sqlmock.NewRows([]string{"username", "email", "role"}).
				AddRow("username test", "test@example.com", "admin")

			mock.ExpectQuery(regexp.QuoteMeta("SELECT username, email, role FROM customer WHERE username = $1")).
				WithArgs(user.Username).
				WillReturnRows(rows)

			userResponse, err = repo.GetUserByUsername(user.Username)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(userResponse).ShouldNot(BeNil())
			Expect(userResponse.Email).Should(Equal("test@example.com"))
			Expect(userResponse.Role).Should(Equal("admin"))
		})
		It("should return error when user not found", func() {
			mock.ExpectQuery(regexp.QuoteMeta("SELECT username, email, role FROM customer WHERE username = $1")).
				WithArgs(user.Username).
				WillReturnError(sql.ErrNoRows)

			userResponse, err = repo.GetUserByUsername(user.Username)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("user not found"))
			Expect(userResponse).Should(BeNil())
		})
		It("should return error on query failure", func() {
			mock.ExpectQuery(regexp.QuoteMeta("SELECT username, email, role FROM customer WHERE username = $1")).
				WithArgs(user.Username).
				WillReturnError(errors.New("query error"))

			userResponse, err = repo.GetUserByUsername(user.Username)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("failed to get user: query error"))
			Expect(userResponse).Should(BeNil())
		})
	})

	Describe("GetAllUsers", func() {
		It("should return all users successfully", func() {
			rows := sqlmock.NewRows([]string{"email", "role"}).
				AddRow("user1@example.com", "admin").
				AddRow("user2@example.com", "user")

			mock.ExpectQuery(regexp.QuoteMeta("SELECT email, role FROM customer")).
				WillReturnRows(rows)

			users, err = repo.GetAllUsers()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(users).Should(HaveLen(2))
			Expect(users[0].Email).Should(Equal("user1@example.com"))
			Expect(users[0].Role).Should(Equal("admin"))
			Expect(users[1].Email).Should(Equal("user2@example.com"))
			Expect(users[1].Role).Should(Equal("user"))
		})
		It("should return error on query failure", func() {
			mock.ExpectQuery(regexp.QuoteMeta("SELECT email, role FROM customer")).
				WillReturnError(errors.New("query error"))

			users, err = repo.GetAllUsers()
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("failed to get all users: query error"))
			Expect(users).Should(BeNil())
		})
		It("should return error on scan failure", func() {
			rows := sqlmock.NewRows([]string{"email", "role"}).
				AddRow("user1@example.com", "admin").
				AddRow("user2@example.com", nil)

			mock.ExpectQuery(regexp.QuoteMeta("SELECT email, role FROM customer")).
				WillReturnRows(rows)

			users, err = repo.GetAllUsers()
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("failed to scan user"))
			Expect(users).Should(BeNil())
		})
	})

	Describe("UpdateUser", func() {
		It("should update user successfully", func() {
			user := &models.User{
				ID:       "1",
				Username: "updated_username",
				Email:    "updated_email@example.com",
				Password: "new_password",
				Role:     "user",
			}

			mock.ExpectExec(regexp.QuoteMeta("UPDATE customer SET username = $2, email = $3, role = $4 WHERE id = $1")).
				WithArgs(user.ID, user.Username, user.Email, user.Role).
				WillReturnResult(sqlmock.NewResult(0, 1))

			err := repo.UpdateUser(user)
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should return error on database failure", func() {
			user := &models.User{
				ID:       "1",
				Username: "updated_username",
				Email:    "updated_email@example.com",
				Password: "new_password",
				Role:     "user",
			}

			mock.ExpectExec(regexp.QuoteMeta("UPDATE customer SET username = $2, email = $3, role = $4 WHERE id = $1")).
				WithArgs(user.ID, user.Username, user.Email, user.Role).
				WillReturnError(fmt.Errorf("database error"))

			err := repo.UpdateUser(user)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("failed to update user"))
		})
	})

	Describe("CreateUser", func() {
		It("should create user successfully", func() {
			user := &models.User{
				Username: "test_user",
				Email:    "test_user@example.com",
				Password: "password",
				Role:     "user",
			}

			mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS (SELECT 1 FROM customer WHERE email = $1)")).
				WithArgs(user.Email).
				WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO customer (username, email, password, role) VALUES ($1, $2, $3, $4)")).
				WithArgs(user.Username, user.Email, user.Password, user.Role).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := repo.CreateUser(user)
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should return error user already exists", func() {
			user := &models.User{
				Username: "existing_user",
				Email:    "existing_user@example.com",
				Password: "password",
				Role:     "user",
			}

			mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS (SELECT 1 FROM customer WHERE email = $1)")).
				WithArgs(user.Email).
				WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

			err := repo.CreateUser(user)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("user already exists"))
		})

		It("should return error on database failure", func() {
			user := &models.User{
				Username: "test_user",
				Email:    "test_user@example.com",
				Password: "password",
				Role:     "user",
			}

			mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS (SELECT 1 FROM customer WHERE email = $1)")).
				WithArgs(user.Email).
				WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO customer (username, email, password, role) VALUES ($1, $2, $3, $4)")).
				WithArgs(user.Username, user.Email, user.Password, user.Role).
				WillReturnError(fmt.Errorf("database error"))

			err := repo.CreateUser(user)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("failed to create user"))
		})
	})

	Describe("DeleteUser", func() {
		It("should delete user successfully", func() {
			userID := "123"

			mock.ExpectExec(regexp.QuoteMeta("DELETE FROM customer WHERE id = $1")).
				WithArgs(userID).
				WillReturnResult(sqlmock.NewResult(0, 1))

			err := repo.DeleteUser(userID)
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should return error on database failure", func() {
			userID := "123"

			mock.ExpectExec(regexp.QuoteMeta("DELETE FROM customer WHERE id = $1")).
				WithArgs(userID).
				WillReturnError(fmt.Errorf("database error"))

			err := repo.DeleteUser(userID)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("failed to delete user"))
		})
	})
})
