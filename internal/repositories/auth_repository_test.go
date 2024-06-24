package repositories

import (
	"awesomeProject/internal/models"
	_ "awesomeProject/pkg/database"
	"database/sql"
	"errors"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AuthRepository", func() {
	var (
		db   *sql.DB
		mock sqlmock.Sqlmock
		repo *AuthRepositoryImpl
		user *models.UserResponse
		auth *models.Auth
		err  error
	)

	BeforeEach(func() {
		db, mock, err = sqlmock.New()
		Expect(err).Should(BeNil())

		repo = &AuthRepositoryImpl{db: db}
		user = &models.UserResponse{}
		auth = &models.Auth{
			Username: "testuser",
			Email:    "testuser@example.com",
			Password: "password",
			Role:     "user",
		}
	})

	AfterEach(func() {
		db.Close()
	})

	Context("Register user", func() {
		It("should register user", func() {
			mock.ExpectExec("INSERT INTO customer").
				WithArgs(auth.Username, auth.Email, auth.Password, auth.Role).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err = repo.Register(auth)
			Expect(err).Should(BeNil())
		})

		It("should return an error if there's a database error", func() {
			mock.ExpectExec("INSERT INTO customer").
				WithArgs(auth.Username, auth.Email, auth.Password, auth.Role).
				WillReturnError(errors.New("database error"))

			err = repo.Register(auth)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("database error"))
		})
	})

	Describe("Login", func() {
		It("should login a valid user", func() {
			rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "role"}).
				AddRow(1, "testuser", "testuser@example.com", "hashedpassword", "user")

			mock.ExpectQuery(regexp.QuoteMeta("SELECT id, username, email, password, role FROM customer WHERE email = $1")).
				WithArgs(auth.Email).
				WillReturnRows(rows)

			user, err = repo.Login(auth)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(user.Email).To(Equal("testuser@example.com"))
		})

		It("should return an error if the user is not found", func() {
			mock.ExpectQuery(regexp.QuoteMeta("SELECT id, username, email, password, role FROM customer WHERE email = $1")).
				WithArgs(auth.Email).
				WillReturnError(sql.ErrNoRows)

			_, err = repo.Login(auth)
			Expect(err).Should(HaveOccurred())
			Expect(errors.Is(err, sql.ErrNoRows)).To(BeTrue())
		})

		It("should return an error if there's a database error", func() {
			mock.ExpectQuery(regexp.QuoteMeta("SELECT id, username, email, password, role FROM customer WHERE email = $1")).
				WithArgs(auth.Email).
				WillReturnError(errors.New("database error"))

			_, err = repo.Login(auth)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("database error"))
		})
	})
})
