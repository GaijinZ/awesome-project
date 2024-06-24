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

var _ = Describe("CategoryRepository", func() {
	var (
		mock             sqlmock.Sqlmock
		db               *sql.DB
		repo             Categorer
		category         *models.Category
		categoryResponse *models.CategoryResponse
		err              error
	)

	BeforeEach(func() {
		db, mock, err = sqlmock.New()
		Expect(err).ShouldNot(HaveOccurred())

		repo = NewCategory(db)
		category = &models.Category{
			ID:        "1",
			Name:      "test product",
			ProductID: 1,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		}

		categoryResponse = &models.CategoryResponse{}
	})

	AfterEach(func() {
		err := mock.ExpectationsWereMet()
		Expect(err).ShouldNot(HaveOccurred())
		db.Close()
	})

	Describe("GetCategory", func() {
		It("should return category when found", func() {
			now := time.Now()
			rows := sqlmock.NewRows([]string{"name", "product_id", "created_at", "updated_at"}).
				AddRow("Books", 1, now, now)

			mock.ExpectQuery(regexp.QuoteMeta("SELECT name, product_id, created_at, updated_at FROM category WHERE id = $1")).
				WithArgs(category.ID).
				WillReturnRows(rows)

			categoryResponse, err = repo.GetCategory(category.ID)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(categoryResponse).ShouldNot(BeNil())
			Expect(categoryResponse.Name).Should(Equal("Books"))
			Expect(categoryResponse.ProductID).Should(Equal(1))
			Expect(categoryResponse.CreatedAt).Should(Equal(now))
			Expect(categoryResponse.UpdatedAt).Should(Equal(now))
		})

		It("should return error when category not found", func() {
			mock.ExpectQuery(regexp.QuoteMeta("SELECT name, product_id, created_at, updated_at FROM category WHERE id = $1")).
				WithArgs(category.ID).
				WillReturnError(sql.ErrNoRows)

			categoryResponse, err = repo.GetCategory(category.ID)
			Expect(err).Should(HaveOccurred())
			Expect(categoryResponse).Should(BeNil())
			Expect(errors.Is(err, sql.ErrNoRows)).Should(BeTrue())
		})

		It("should return error on query failure", func() {
			mock.ExpectQuery(regexp.QuoteMeta("SELECT name, product_id, created_at, updated_at FROM category WHERE id = $1")).
				WithArgs(category.ID).WillReturnError(errors.New("query error"))

			categoryResponse, err = repo.GetCategory(category.ID)
			Expect(err).Should(HaveOccurred())
			Expect(categoryResponse).Should(BeNil())
			Expect(err.Error()).Should(ContainSubstring("failed to get category: query error"))
		})

		It("should return error on scan failure", func() {
			rows := sqlmock.NewRows([]string{"name", "product_id", "created_at", "updated_at"}).
				AddRow("test category", 1, "invalid time", time.Now())

			mock.ExpectQuery(regexp.QuoteMeta("SELECT name, product_id, created_at, updated_at FROM category WHERE id = $1")).
				WithArgs(category.ID).
				WillReturnRows(rows)

			categoryResponse, err = repo.GetCategory(category.ID)
			Expect(err).Should(HaveOccurred())
			Expect(categoryResponse).Should(BeNil())
			Expect(err.Error()).Should(ContainSubstring("failed to get category: sql: Scan error"))
		})
	})

	Describe("Update Category", func() {
		It("should update category successfully", func() {
			mock.ExpectExec(regexp.QuoteMeta("UPDATE category SET name = $2, product_id = $3, updated_at = $4 WHERE id = $1")).
				WithArgs(category.ID, category.Name, category.ProductID, sqlmock.AnyArg()).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := repo.UpdateCategory(*category)
			Expect(err).Should(BeNil())
		})

		It("should return error when update fails", func() {
			mock.ExpectExec(regexp.QuoteMeta("UPDATE category SET name = $2, product_id = $3, updated_at = $4 WHERE id = $1")).
				WithArgs(category.ID, category.Name, category.ProductID, sqlmock.AnyArg()).
				WillReturnError(errors.New("update error"))

			err := repo.UpdateCategory(*category)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("failed to update category: update error"))
		})

		It("should return error when category ID is missing", func() {
			category.ID = ""

			err := repo.UpdateCategory(*category)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("failed to update category:"))
		})
	})

	Describe("Create Category", func() {
		It("should create category successfully", func() {
			mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS (SELECT 1 FROM category WHERE name = $1)")).
				WithArgs(category.Name).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO category (name, product_id) VALUES ($1, $2)")).
				WithArgs(category.Name, category.ProductID).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := repo.CreateCategory(*category)
			Expect(err).Should(BeNil())
		})

		It("should return error when category already exists", func() {
			mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS (SELECT 1 FROM category WHERE name = $1)")).
				WithArgs(category.Name).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

			err := repo.CreateCategory(*category)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("category already exists"))
		})

		It("should return error when checking category existence fails", func() {
			mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS (SELECT 1 FROM category WHERE name = $1)")).
				WithArgs(category.Name).WillReturnError(errors.New("db error"))

			err := repo.CreateCategory(*category)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("failed to check category exist: db error"))
		})

		It("should return error when creating category fails", func() {
			mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS (SELECT 1 FROM category WHERE name = $1)")).
				WithArgs(category.Name).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

			mock.ExpectExec(regexp.QuoteMeta(
				"INSERT INTO category (name, product_id) VALUES ($1, $2)")).
				WithArgs(category.Name, category.ProductID).
				WillReturnError(errors.New("insert error"))

			err := repo.CreateCategory(*category)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("failed to create category: insert error"))
		})
	})

	Describe("Delete Category", func() {
		It("should delete category successfully", func() {
			mock.ExpectExec(regexp.QuoteMeta("DELETE FROM category WHERE id = $1")).
				WithArgs(category.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := repo.DeleteCategory(category.ID)
			Expect(err).Should(BeNil())
		})

		It("should return error when deletion fails", func() {
			mock.ExpectExec(regexp.QuoteMeta("DELETE FROM category WHERE id = $1")).
				WithArgs(category.ID).
				WillReturnError(errors.New("delete error"))

			err := repo.DeleteCategory(category.ID)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("delete error"))
		})
	})
})
