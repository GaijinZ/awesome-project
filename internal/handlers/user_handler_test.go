package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"awesomeProject/internal/handlers/mocks"
	"awesomeProject/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("User Handler", func() {
	var (
		mockCtrl         *gomock.Controller
		mockRepo         *mocks.MockUserRepository
		userHandler      Userer
		responseRecorder *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockUserRepository(mockCtrl)
		userHandler = NewUserHandler(mockRepo)
		responseRecorder = httptest.NewRecorder()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Get User by name Handler", func() {
		It("should return 200", func() {
			request, err := http.NewRequest("GET", "/api/v1/users/1", nil)
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			userID := mux.Vars(request)["user_id"]

			mockRepo.EXPECT().
				GetUserByUsername(userID).
				Return(&models.UserResponse{
					Email: "user@example.com",
					Role:  "user",
				}, nil).Times(1)

			userHandler.GetUserByUsername(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
		})
		It("should return 404 when no category id provided", func() {
			request, err := http.NewRequest("GET", "/api/v1/users", nil)
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			categoryID := mux.Vars(request)["user_id"]

			mockRepo.EXPECT().
				GetUserByUsername(categoryID).
				Return(nil, errors.New("not found")).
				Times(1)

			userHandler.GetUserByUsername(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
		})
		It("should return 404, wrong category id provided", func() {
			request, err := http.NewRequest("GET", "/api/v1/users/2", nil)
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			categoryID := mux.Vars(request)["user_id"]

			mockRepo.EXPECT().
				GetUserByUsername(categoryID).
				Return(nil, errors.New("not found")).
				Times(1)

			userHandler.GetUserByUsername(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
		})
	})

	Describe("Get All Users Handler", func() {
		It("should return 200", func() {
			request, _ := http.NewRequest("GET", "/api/v1/users", nil)
			request.Header.Set("Content-Type", "application/json")

			expectedUsers := []models.UserResponse{
				{Email: "user1@example.com", Role: "user"},
				{Email: "user2@example.com", Role: "admin"},
			}

			mockRepo.EXPECT().
				GetAllUsers().
				Return(expectedUsers, nil).
				Times(1)

			userHandler.GetAllUsers(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
		})
		It("should return 404 not found", func() {
			request, err := http.NewRequest("GET", "/api/v1/users", nil)
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			mockRepo.EXPECT().
				GetAllUsers().
				Return([]models.UserResponse{}, errors.New("no users found"))

			userHandler.GetAllUsers(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
		})
		It("should return 404 not found", func() {
			request, err := http.NewRequest("GET", "/api/v1/users", nil)
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			mockRepo.EXPECT().GetAllUsers().Return(nil, errors.New("database error"))

			userHandler.GetAllUsers(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
		})
	})

	Describe("UpdateUser", func() {
		It("should update user successfully", func() {
			user := &models.User{
				Username: "updated_username",
				Email:    "updated_email@example.com",
				Role:     "user",
			}

			userResponse := &models.UserResponse{
				Username: "testuser",
				Email:    "testuser@example.com",
				Role:     "user",
			}

			userJSON, err := json.Marshal(user)
			Expect(err).NotTo(HaveOccurred())

			request, err := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer(userJSON))
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")
			request.AddCookie(&http.Cookie{Name: "userID", Value: "1"})
			request.AddCookie(&http.Cookie{Name: "username", Value: "testuser"})

			request = mux.SetURLVars(request, map[string]string{"user_id": "1"})
			user.ID = "1"

			mockRepo.EXPECT().
				GetUserByUsername("testuser").
				Return(userResponse, nil).
				Times(1)

			mockRepo.EXPECT().
				UpdateUser(user).
				Return(nil).
				Times(1)

			userHandler.UpdateUser(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
		})
		It("should return 401 Unauthorized", func() {
			user := models.User{Username: "testuser"}

			requestBody, err := json.Marshal(user)
			Expect(err).NotTo(HaveOccurred())
			request, err := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer(requestBody))
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")
			request.AddCookie(&http.Cookie{Name: "userID", Value: "2"})

			userHandler.UpdateUser(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusUnauthorized))
			Expect(responseRecorder.Body.String()).To(Equal("Unauthorized\n"))
		})
		It("should return 404 Not Found, missing cookie", func() {
			user := models.User{Username: "testuser"}

			requestBody, err := json.Marshal(user)
			Expect(err).NotTo(HaveOccurred())
			request, err := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer(requestBody))
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")
			request.AddCookie(&http.Cookie{Name: "userID", Value: "1"})

			request = mux.SetURLVars(request, map[string]string{"user_id": "1"})
			user.ID = "1"

			userHandler.UpdateUser(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
			Expect(responseRecorder.Body.String()).To(Equal("Cookie not found\n"))
		})
		It("should return 404 Not Found in repository", func() {
			user := models.User{Username: "testuser"}

			requestBody, err := json.Marshal(user)
			Expect(err).NotTo(HaveOccurred())
			request, err := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer(requestBody))
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")
			request.AddCookie(&http.Cookie{Name: "userID", Value: "1"})
			request.AddCookie(&http.Cookie{Name: "username", Value: "testuser"})

			request = mux.SetURLVars(request, map[string]string{"user_id": "1"})
			user.ID = "1"

			mockRepo.EXPECT().
				GetUserByUsername(user.Username).
				Return(nil, errors.New("user not found")).
				Times(1)

			userHandler.UpdateUser(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
			Expect(responseRecorder.Body.String()).To(Equal("User not found\n"))
		})
	})

	Describe("CreateUser", func() {
		It("should return 400 Bad Request if JSON is invalid", func() {
			request, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBufferString("invalid json"))
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			userHandler.CreateUser(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
			Expect(responseRecorder.Body.String()).To(ContainSubstring("Failed to encode request"))
		})
		It("should return 404 Not Found if user already exists", func() {
			user := &models.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password",
				Role:     "user",
			}

			userJSON, err := json.Marshal(user)
			Expect(err).NotTo(HaveOccurred())

			request, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(userJSON))
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			mockRepo.EXPECT().
				CreateUser(gomock.Any()).
				Return(errors.New("user already exists")).
				Times(1)

			userHandler.CreateUser(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
			Expect(responseRecorder.Body.String()).To(ContainSubstring("Failed to create user"))
		})
		It("should return 200 OK if user is created successfully", func() {
			user := &models.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password",
				Role:     "user",
			}

			userJSON, err := json.Marshal(user)
			Expect(err).NotTo(HaveOccurred())

			request, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(userJSON))
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			mockRepo.EXPECT().
				CreateUser(gomock.Any()).
				Return(nil).
				Times(1)

			userHandler.CreateUser(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
		})
	})

	Describe("DeleteUser", func() {
		It("should return 404 Not Found if user not found", func() {
			request, err := http.NewRequest("DELETE", "/api/v1/users/{user_id}", nil)
			Expect(err).NotTo(HaveOccurred())

			userID := mux.Vars(request)["user_id"]

			mockRepo.EXPECT().
				DeleteUser(userID).
				Return(errors.New("user not found")).
				Times(1)

			userHandler.DeleteUser(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
			Expect(responseRecorder.Body.String()).To(ContainSubstring("Failed to delete user"))
		})
		It("should return 200 OK if user is deleted successfully", func() {
			request, err := http.NewRequest("DELETE", "/api/v1/users/{user_id}", nil)
			Expect(err).NotTo(HaveOccurred())

			userID := mux.Vars(request)["user_id"]

			mockRepo.EXPECT().
				DeleteUser(userID).
				Return(nil).
				Times(1)

			userHandler.DeleteUser(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
		})
	})
})
