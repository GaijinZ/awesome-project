package repositories

import (
	"database/sql"
	"errors"
	"regexp"
	"time"

	"awesomeProject/internal/models"

	"github.com/DATA-DOG/go-sqlmock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Product Repository", func() {
	var (
		db              *sql.DB
		mock            sqlmock.Sqlmock
		repo            ProductRepository
		product         *models.Product
		productResponse *models.ProductResponse
		err             error
	)

	BeforeEach(func() {
		db, mock, err = sqlmock.New()
		Expect(err).Should(BeNil())

		repo = NewProduct(db)
		product = &models.Product{
			ID:         "1",
			Name:       "test product",
			CategoryID: 1,
			CreatedAt:  time.Time{},
			UpdatedAt:  time.Time{},
		}

		productResponse = &models.ProductResponse{}
	})

	AfterEach(func() {
		db.Close()
	})

	Describe("Get Product", func() {
		It("should return product, nil", func() {
			rows := sqlmock.NewRows([]string{"name", "category_id", "created_at", "updated_at"}).
				AddRow("test product", 1, time.Time{}, time.Time{})

			mock.ExpectQuery(regexp.QuoteMeta(
				"SELECT name, category_id, created_at, updated_at FROM products WHERE id = $1")).
				WithArgs(product.ID).WillReturnRows(rows)

			productResponse, err = repo.GetProduct(product.ID)
			Expect(err).Should(BeNil())
			Expect(productResponse).Should(Equal(&models.ProductResponse{
				Name:       "test product",
				CategoryID: 1,
				CreatedAt:  time.Time{},
				UpdatedAt:  time.Time{},
			}))
		})
		It("should return error when product not found", func() {
			mock.ExpectQuery(regexp.QuoteMeta("SELECT name, category_id, created_at, updated_at FROM products WHERE id = $1")).
				WithArgs(product.ID).
				WillReturnError(sql.ErrNoRows)

			productResponse, err := repo.GetProduct(product.ID)
			Expect(err).Should(HaveOccurred())
			Expect(productResponse).Should(BeNil())
			Expect(errors.Is(err, sql.ErrNoRows)).Should(BeTrue())
		})
		It("should return error on query failure", func() {
			mock.ExpectQuery(regexp.QuoteMeta("SELECT name, category_id, created_at, updated_at FROM products WHERE id = $1")).
				WithArgs(product.ID).WillReturnError(errors.New("query error"))

			productResponse, err = repo.GetProduct(product.ID)
			Expect(err).Should(HaveOccurred())
			Expect(productResponse).Should(BeNil())
			Expect(err.Error()).Should(ContainSubstring("query error"))
		})
		It("should return error on scan failure", func() {
			rows := sqlmock.NewRows([]string{"name", "category_id", "created_at", "updated_at"}).
				AddRow("test product", 1, "invalid time", time.Now())

			mock.ExpectQuery(regexp.QuoteMeta("SELECT name, category_id, created_at, updated_at FROM products WHERE id = $1")).
				WithArgs(product.ID).WillReturnRows(rows)

			productResponse, err = repo.GetProduct(product.ID)
			Expect(err).Should(HaveOccurred())
			Expect(productResponse).Should(BeNil())
			Expect(err.Error()).Should(ContainSubstring("sql: Scan error"))
		})
	})

	Describe("UpdateProduct", func() {
		It("should update product successfully", func() {
			mock.ExpectExec(regexp.QuoteMeta("UPDATE products SET name = $2, category_id = $3, updated_at = $4 WHERE id = $1")).
				WithArgs(product.ID, product.Name, product.CategoryID, sqlmock.AnyArg()).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := repo.UpdateProduct(product)
			Expect(err).Should(BeNil())
		})

		It("should return error when update fails", func() {
			mock.ExpectExec(regexp.QuoteMeta("UPDATE products SET name = $2, category_id = $3, updated_at = $4 WHERE id = $1")).
				WithArgs(product.ID, product.Name, product.CategoryID, sqlmock.AnyArg()).
				WillReturnError(errors.New("update error"))

			err := repo.UpdateProduct(product)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("failed to update product: update error"))
		})

		It("should return error when product ID is missing", func() {
			product.ID = ""

			err := repo.UpdateProduct(product)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("failed to update product:"))
		})
	})

	Describe("CreateProduct", func() {
		It("should create product successfully", func() {
			mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS (SELECT 1 FROM products WHERE name = $1)")).
				WithArgs(product.Name).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO products (name, category_id) VALUES ($1, $2)")).
				WithArgs(product.Name, product.CategoryID).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := repo.CreateProduct(product)
			Expect(err).Should(BeNil())
		})

		It("should return error when product already exists", func() {
			mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS (SELECT 1 FROM products WHERE name = $1)")).
				WithArgs(product.Name).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

			err := repo.CreateProduct(product)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("product already exists"))
		})

		It("should return error when checking product existence fails", func() {
			mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS (SELECT 1 FROM products WHERE name = $1)")).
				WithArgs(product.Name).WillReturnError(errors.New("db error"))

			err := repo.CreateProduct(product)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("failed to check product exist: db error"))
		})

		It("should return error when creating product fails", func() {
			mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS (SELECT 1 FROM products WHERE name = $1)")).
				WithArgs(product.Name).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

			mock.ExpectExec(regexp.QuoteMeta(
				"INSERT INTO products (name, category_id) VALUES ($1, $2)")).
				WithArgs(product.Name, product.CategoryID).
				WillReturnError(errors.New("insert error"))

			err := repo.CreateProduct(product)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("failed to create product: insert error"))
		})
	})

	Describe("DeleteProduct", func() {
		It("should delete product successfully", func() {
			mock.ExpectExec(regexp.QuoteMeta("DELETE FROM products WHERE id = $1")).
				WithArgs(product.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := repo.DeleteProduct(product.ID)
			Expect(err).Should(BeNil())
		})

		It("should return error when deletion fails", func() {
			mock.ExpectExec(regexp.QuoteMeta("DELETE FROM products WHERE id = $1")).
				WithArgs(product.ID).
				WillReturnError(errors.New("delete error"))

			err := repo.DeleteProduct(product.ID)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("delete error"))
		})
	})
})
