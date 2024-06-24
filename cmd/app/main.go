package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"awesomeProject/internal/handlers"
	"awesomeProject/internal/repositories"
	"awesomeProject/internal/routers"
	"awesomeProject/pkg/database"
)

func main() {
	configDB := database.ConnectionConfig{
		DriverName:      os.Getenv("AWP_DB_DRIVER"),
		DataSourceName:  os.Getenv("AWP_DB_DATASOURCE"),
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	}

	db, err := database.NewConnection(configDB)
	if err != nil {
		fmt.Errorf("failed to configure db connection: %v", err)
	}
	log.Printf("Pinging db %v", db.Ping())

	userRepository := repositories.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepository)
	authRepository := repositories.NewAuthRepositoryImpl(db)
	authHandler := handlers.NewAuth(authRepository)
	categoryRepository := repositories.NewCategory(db)
	categoryHandler := handlers.NewCategoryHandler(categoryRepository)
	productRepository := repositories.NewProduct(db)
	productHandler := handlers.NewProductHandler(productRepository)

	router := routers.NewRouter(userHandler, categoryHandler, productHandler, authHandler)

	httpServer := http.Server{
		Addr:    ":" + os.Getenv("AWP_PORT"),
		Handler: router,
	}

	fmt.Printf("Server starting at :%v", os.Getenv("AWP_PORT")+"\n")
	go func() {
		if err = httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP routers error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		db.Close()
		shutdownRelease()
	}()

	if err = httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}

	log.Println("Graceful shutdown complete.")
}
