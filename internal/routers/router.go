package routers

import (
	"awesomeProject/internal/handlers"
	"awesomeProject/internal/middlewares/authentication"
	"awesomeProject/internal/middlewares/logging"
	"awesomeProject/internal/middlewares/middleware"
	"github.com/gorilla/mux"
)

func NewRouter(users handlers.Userer, categories handlers.Categorer, products handlers.Producter, auth handlers.Auther) *mux.Router {
	router := mux.NewRouter()

	middlewares := []middleware.Middleware{
		logging.LoggingMiddleware,
		authentication.IsAuthenticated,
	}

	r := router.PathPrefix("/api/v1").Subrouter()

	r.HandleFunc("/signup", auth.Register).Methods("POST")
	r.HandleFunc("/login", auth.Login).Methods("POST")
	r.HandleFunc("/logout", auth.Logout).Methods("POST")

	r.HandleFunc("/users/{username}", middleware.ChainMiddleware(
		users.GetUserByUsername,
		middlewares...,
	)).Methods("GET")
	r.HandleFunc("/users", middleware.ChainMiddleware(
		users.GetAllUsers,
		middlewares...,
	)).Methods("GET")
	r.HandleFunc("/users/{user_id}", middleware.ChainMiddleware(
		users.UpdateUser,
		middlewares...,
	)).Methods("PUT")
	r.HandleFunc("/users", middleware.ChainMiddleware(
		users.CreateUser,
		middlewares...,
	)).Methods("POST")
	r.HandleFunc("/users/{user_id}", middleware.ChainMiddleware(
		users.DeleteUser,
		middlewares...,
	)).Methods("DELETE")

	r.HandleFunc("/categories", middleware.ChainMiddleware(
		categories.CreateCategoryHandler,
		middlewares...,
	)).Methods("POST")
	r.HandleFunc("/categories/{category_id}", middleware.ChainMiddleware(
		categories.GetCategoryHandler,
		middlewares...)).Methods("GET")
	r.HandleFunc("/categories/{category_id}", middleware.ChainMiddleware(
		categories.UpdateCategoryHandler,
		middlewares...)).Methods("PUT")
	r.HandleFunc("/categories/{category_id}", middleware.ChainMiddleware(
		categories.DeleteCategoryHandler,
		middlewares...)).Methods("DELETE")

	r.HandleFunc("/products", middleware.ChainMiddleware(
		products.CreateProductHandler,
		middlewares...)).Methods("POST")
	r.HandleFunc("/products/{product_id}", middleware.ChainMiddleware(
		products.GetProductHandler,
		middlewares...)).Methods("GET")
	r.HandleFunc("/products/{product_id}", middleware.ChainMiddleware(
		products.UpdateProductHandler,
		middlewares...)).Methods("PUT")
	r.HandleFunc("/products/{product_id}", middleware.ChainMiddleware(
		products.DeleteProductHandler,
		middlewares...)).Methods("DELETE")

	return r
}
