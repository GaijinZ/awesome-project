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

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Category Handler", func() {
	var (
		mockCtrl         *gomock.Controller
		mockRepo         *mocks.MockProductRepository
		productHandler   Producter
		responseRecorder *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockProductRepository(mockCtrl)
		productHandler = NewProductHandler(mockRepo)
		responseRecorder = httptest.NewRecorder()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("GetProductHandler", func() {
		It("should return 200", func() {
			request, err := http.NewRequest("GET", "/api/v1/products/1", nil)
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			productID := mux.Vars(request)["product_id"]

			mockRepo.EXPECT().
				GetProduct(productID).
				Return(&models.ProductResponse{
					Name:       "Hobbit",
					CategoryID: 1,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}, nil).
				Times(1)

			productHandler.GetProductHandler(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
		})
		It("should return 404 when no category id provided", func() {
			request, err := http.NewRequest("GET", "/api/v1/products", nil)
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			productID := mux.Vars(request)["product_id"]

			mockRepo.EXPECT().
				GetProduct(productID).
				Return(nil, errors.New("no category id provided")).
				Times(1)

			productHandler.GetProductHandler(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
		})
		It("should return 404, wrong category id provided", func() {
			request, err := http.NewRequest("GET", "/api/v1/products/2", nil)
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			productID := mux.Vars(request)["product_id"]

			mockRepo.EXPECT().
				GetProduct(productID).
				Return(nil, errors.New("no category id provided")).
				Times(1)

			productHandler.GetProductHandler(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
		})
	})

	Describe("UpdateProductHandler", func() {
		It("should return 200", func() {
			product := &models.Product{
				Name:       "Hobbit",
				CategoryID: 1,
				UpdatedAt:  time.Now().Round(time.Microsecond),
			}

			requestBody, err := json.Marshal(product)
			Expect(err).NotTo(HaveOccurred())
			request, err := http.NewRequest("PUT", "/api/v1/products/1", bytes.NewBuffer(requestBody))
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			mockRepo.EXPECT().
				UpdateProduct(gomock.Eq(product)).
				Return(nil).
				Times(1)

			productHandler.UpdateProductHandler(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
		})
		It("should return 200, empty request", func() {
			product := &models.Product{}

			requestBody, err := json.Marshal(product)
			Expect(err).NotTo(HaveOccurred())
			request, err := http.NewRequest("PUT", "/api/v1/products/1", bytes.NewBuffer(requestBody))
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			mockRepo.EXPECT().
				UpdateProduct(gomock.Eq(product)).
				Return(nil).
				Times(1)

			productHandler.UpdateProductHandler(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
		})
		It("should return 400, nil body", func() {
			request, err := http.NewRequest("PUT", "/api/v1/products/1", nil)
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			mockRepo.EXPECT().UpdateProduct(gomock.Any()).Times(0)

			productHandler.UpdateProductHandler(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Describe("CreateProductHandler", func() {
		It("should return 200", func() {
			product := &models.Product{Name: "testBook"}

			requestBody, err := json.Marshal(product)
			Expect(err).NotTo(HaveOccurred())
			request, err := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(requestBody))
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			mockRepo.EXPECT().CreateProduct(product).Return(nil).Times(1)

			productHandler.CreateProductHandler(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusCreated))
		})
		It("should return 400, nil body", func() {
			request, err := http.NewRequest("POST", "/api/v1/products", nil)
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			mockRepo.EXPECT().CreateProduct(gomock.Any()).Times(0)

			productHandler.CreateProductHandler(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
		})
		It("should return 400, empty body", func() {
			category := &models.Category{}

			requestBody, _ := json.Marshal(category)
			request, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(requestBody))
			request.Header.Set("Content-Type", "application/json")

			mockRepo.EXPECT().CreateProduct(gomock.Any()).Times(0)

			productHandler.CreateProductHandler(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Describe("DeleteProductHandler", func() {
		It("should return 200", func() {
			request, err := http.NewRequest("DELETE", "/api/v1/products/1", nil)
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			categoryID := mux.Vars(request)["category_id"]

			mockRepo.EXPECT().DeleteProduct(categoryID).Return(nil).Times(1)
			productHandler.DeleteProductHandler(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
		})
		It("should return 404 when no category id provided", func() {
			request, err := http.NewRequest("DELETE", "/api/v1/categories", nil)
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			mockRepo.EXPECT().DeleteProduct(gomock.Any()).Return(errors.New("not found")).Times(1)
			productHandler.DeleteProductHandler(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
		})
	})
})
