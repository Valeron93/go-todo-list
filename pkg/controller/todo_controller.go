package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Valeron93/go-todo-list/pkg/repository"
	"github.com/Valeron93/go-todo-list/pkg/schemas"
)

type TodoController struct {
	repo   *repository.TodoRepository
	logger *slog.Logger
}

func NewTodoController(repo *repository.TodoRepository, logger *slog.Logger) (*TodoController, error) {
	return &TodoController{
		repo:   repo,
		logger: logger,
	}, nil
}

func (c *TodoController) RegisterRoutes() http.Handler {

	router := http.NewServeMux()

	router.HandleFunc("POST /todos", c.handleAddTodo)
	router.HandleFunc("GET /todos", c.handleGetAllTodos)
	router.HandleFunc("GET /todos/{id}", c.handleGetTodo)
	router.HandleFunc("DELETE /todos/{id}", c.handleDeleteTodo)
	router.HandleFunc("PATCH /todos/{id}", c.handlePatchUpdate)

	return router

}
func (c *TodoController) handleAddTodo(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	var todoCreate schemas.TodoCreate

	if !readJson(w, r, &todoCreate) {
		return
	}

	todoItem, err := c.repo.AddTodo(r.Context(), todoCreate)
	if err != nil {
		c.dbError(w, err)
		return
	}

	writeJson(w, todoItem)
	w.WriteHeader(http.StatusCreated)
}

func (c *TodoController) handleGetAllTodos(w http.ResponseWriter, r *http.Request) {

	todos, err := c.repo.GetAllTodos(r.Context())
	if err != nil {
		c.dbError(w, err)
		return
	}

	writeJson(w, todos)
}

func (c *TodoController) handleGetTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	todoItem, err := c.repo.GetTodo(r.Context(), id)
	if err != nil {
		c.dbError(w, err)
		return
	}

	writeJson(w, todoItem)
}

func (c *TodoController) handleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := c.repo.DeleteTodo(r.Context(), id); err != nil {
		c.dbError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *TodoController) handlePatchUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var todoItemUpdate schemas.TodoPatchUpdate

	defer r.Body.Close()
	if !readJson(w, r, &todoItemUpdate) {
		return
	}

	todoItem, err := c.repo.PatchTodo(r.Context(), id, todoItemUpdate)
	if err != nil {
		c.dbError(w, err)
		return
	}

	writeJson(w, todoItem)
}

func (c *TodoController) dbError(w http.ResponseWriter, err error) {
	if err == repository.ErrTodoNotFound {
		http.Error(w, "not found", http.StatusNotFound)
	} else {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		c.logger.Error("db error", slog.Any("err", err))
	}
}

func writeJson(w http.ResponseWriter, object any) error {
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(object)
}

func readJson(w http.ResponseWriter, r *http.Request, object any) bool {
	if err := json.NewDecoder(r.Body).Decode(object); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return false
	}

	return true
}
