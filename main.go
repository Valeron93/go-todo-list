package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/Valeron93/go-todo-list/pkg/controller"
	"github.com/Valeron93/go-todo-list/pkg/middleware"
	"github.com/Valeron93/go-todo-list/pkg/models"
	"github.com/Valeron93/go-todo-list/pkg/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {

	gormLogger := logger.New(nil, logger.Config{
		LogLevel: logger.Silent,
	})

	db, err := gorm.Open(sqlite.Open("db.sqlite"), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.TodoItem{})
	logger := slog.Default()

	todoRepo := repository.NewTodoRepository(db, logger)

	todoController, err := controller.NewTodoController(todoRepo, logger)
	if err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()
	router.Handle("/api/", http.StripPrefix("/api", todoController.RegisterRoutes()))

	const addr = ":3000"

	loggerMiddleware := middleware.NewSlogMiddleware(logger)
	logger.Info("listening", slog.Any("addr", addr))
	if err := http.ListenAndServe(addr, loggerMiddleware(router)); err != nil {
		log.Fatal(err)
	}

}
