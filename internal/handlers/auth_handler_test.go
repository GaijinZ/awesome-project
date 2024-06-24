package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"time"

	"awesomeProject/internal/handlers/mocks"
	"awesomeProject/internal/models"
	"awesomeProject/pkg/utils"

	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth Handler", func() {
	var (
		mockCtrl         *gomock.Controller
		mockRepo         *mocks.MockAuthRepository
		authHandler      *AuthHandler
		responseRecorder *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockAuthRepository(mockCtrl)
		authHandler = &AuthHandler{authRepository: mockRepo}
		responseRecorder = httptest.NewRecorder()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Register", func() {
		It("should return 201 for successful registration", func() {
			user := &models.Auth{
				Username: "testuser",
				Email:    "testuser@example.com",
				Password: "password",
				Role:     "user",
			}

			requestBody, err := json.Marshal(user)
			Expect(err).NotTo(HaveOccurred())
			request, err := http.NewRequest("POST", "/api/v1/signup", bytes.NewBuffer(requestBody))
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			mockRepo.EXPECT().Register(gomock.Any()).Return(nil).Times(1)

			authHandler.Register(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusCreated))
		})
		It("should return 400 for invalid input", func() {
			request, err := http.NewRequest("POST", "/api/v1/signup", nil)
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			authHandler.Register(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Describe("Login", func() {
		It("should return 202 for successful login", func() {
			auth := &models.Auth{
				Email:    "testuser@example.com",
				Password: "password",
			}

			hashedPassword, err := utils.GenerateHashPassword("password")
			Expect(err).NotTo(HaveOccurred())
			userResponse := &models.UserResponse{
				ID:       "1",
				Username: "test",
				Email:    "testuser@example.com",
				Password: hashedPassword,
				Role:     "user",
			}

			mockRepo.EXPECT().
				Login(gomock.Eq(auth)).
				Return(userResponse, nil).
				Times(1)

			requestBody, err := json.Marshal(auth)
			Expect(err).NotTo(HaveOccurred())
			request, err := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(requestBody))
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			authHandler.Login(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusAccepted))
		})

		It("should return 401 for invalid credentials", func() {
			auth := &models.Auth{
				Email:    "invaliduser@example.com",
				Password: "wrongpassword",
			}
			userResponse := &models.UserResponse{}

			requestBody, err := json.Marshal(auth)
			Expect(err).NotTo(HaveOccurred())
			request, err := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(requestBody))
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			mockRepo.EXPECT().Login(auth).Return(userResponse, errors.New("invalid login credentials")).Times(1)

			authHandler.Login(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusUnauthorized))
		})
	})

	Describe("Logout", func() {
		It("should clear the token cookie and redirect", func() {
			request, err := http.NewRequest("POST", "/api/v1/logout", nil)
			Expect(err).NotTo(HaveOccurred())

			authHandler.Logout(responseRecorder, request)

			cookie := responseRecorder.Result().Cookies()[0]
			Expect(cookie.Name).To(Equal("token"))
			Expect(cookie.Value).To(BeEmpty())
			Expect(cookie.Expires.Before(time.Now())).To(BeTrue())
			Expect(responseRecorder.Code).To(Equal(http.StatusSeeOther))
		})
	})
})
