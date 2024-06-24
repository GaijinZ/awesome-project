package handlers

import (
	"awesomeProject/internal/handlers/mocks"
	"awesomeProject/internal/models"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Category Handler", func() {
	var (
		mockCtrl         *gomock.Controller
		mockRepo         *mocks.MockCategorer
		categoryHandler  Categorer
		responseRecorder *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockCategorer(mockCtrl)
		categoryHandler = NewCategoryHandler(mockRepo)
		responseRecorder = httptest.NewRecorder()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("GetCategoryHandler", func() {
		It("should return 200", func() {
			request, err := http.NewRequest("GET", "/api/v1/categories/1", nil)
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			categoryID := mux.Vars(request)["category_id"]

			mockRepo.EXPECT().
				GetCategory(categoryID).
				Return(&models.CategoryResponse{
					Name:      "Books",
					ProductID: 1,
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				}, nil).
				Times(1)

			categoryHandler.GetCategoryHandler(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
		})
		It("should return 404 when no category id provided", func() {
			request, err := http.NewRequest("GET", "/api/v1/categories", nil)
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			categoryID := mux.Vars(request)["category_id"]

			mockRepo.EXPECT().
				GetCategory(categoryID).
				Return(nil, errors.New("not found")).
				Times(1)

			categoryHandler.GetCategoryHandler(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
		})
		It("should return 404, wrong category id provided", func() {
			request, err := http.NewRequest("GET", "/api/v1/categories/2", nil)
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			categoryID := mux.Vars(request)["category_id"]

			mockRepo.EXPECT().
				GetCategory(categoryID).
				Return(nil, errors.New("not found")).
				Times(1)

			categoryHandler.GetCategoryHandler(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
		})
	})

	Describe("UpdateCategoryHandler", func() {
		It("should return 200", func() {
			category := models.Category{
				Name:      "test",
				ProductID: 1,
				UpdatedAt: time.Time{},
			}

			requestBody, err := json.Marshal(category)
			Expect(err).NotTo(HaveOccurred())
			request, err := http.NewRequest("PUT", "/api/v1/categories/1", bytes.NewBuffer(requestBody))
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			mockRepo.EXPECT().
				UpdateCategory(gomock.Eq(category)).
				Return(nil).
				Times(1)

			categoryHandler.UpdateCategoryHandler(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
		})
		It("should return 200, empty request", func() {
			category := models.Category{}

			requestBody, err := json.Marshal(category)
			Expect(err).NotTo(HaveOccurred())
			request, err := http.NewRequest("PUT", "/api/v1/categories/1", bytes.NewBuffer(requestBody))
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			mockRepo.EXPECT().
				UpdateCategory(gomock.Eq(category)).
				Return(nil).
				Times(1)

			categoryHandler.UpdateCategoryHandler(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
		})
		It("should return 400, nil body", func() {
			request, err := http.NewRequest("PUT", "/api/v1/categories/1", nil)
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			mockRepo.EXPECT().UpdateCategory(gomock.Any()).Times(0)

			categoryHandler.UpdateCategoryHandler(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Describe("CreateCategoryHandler", func() {
		It("should return 200", func() {
			category := models.Category{
				Name:      "testBook",
				ProductID: 1,
			}

			requestBody, err := json.Marshal(category)
			Expect(err).NotTo(HaveOccurred())
			request, err := http.NewRequest("POST", "/api/v1/categories", bytes.NewBuffer(requestBody))
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			mockRepo.EXPECT().CreateCategory(gomock.Any()).Return(nil).Times(1)

			categoryHandler.CreateCategoryHandler(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusCreated))
		})
		It("should return 400, nil body", func() {
			request, err := http.NewRequest("POST", "/api/v1/categories", nil)
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			mockRepo.EXPECT().CreateCategory(gomock.Any()).Times(0)

			categoryHandler.CreateCategoryHandler(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
		})
		It("should return 400, empty body", func() {
			category := models.Category{}

			requestBody, err := json.Marshal(category)
			Expect(err).NotTo(HaveOccurred())
			request, err := http.NewRequest("POST", "/api/v1/categories", bytes.NewBuffer(requestBody))
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			mockRepo.EXPECT().CreateCategory(gomock.Any()).Times(0)

			categoryHandler.CreateCategoryHandler(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Describe("DeleteCategoryHandler", func() {
		It("should return 200", func() {
			request, err := http.NewRequest("DELETE", "/api/v1/categories/1", nil)
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			categoryID := mux.Vars(request)["category_id"]

			mockRepo.EXPECT().DeleteCategory(categoryID).Return(nil).Times(1)
			categoryHandler.DeleteCategoryHandler(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
		})
		It("should return 404 when no category id provided", func() {
			request, err := http.NewRequest("DELETE", "/api/v1/categories", nil)
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			mockRepo.EXPECT().DeleteCategory(gomock.Any()).Return(errors.New("not found")).Times(1)
			categoryHandler.DeleteCategoryHandler(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
		})
	})
})
