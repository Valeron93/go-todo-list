package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Valeron93/go-todo-list/pkg/models"
	"github.com/Valeron93/go-todo-list/pkg/schemas"
	"gorm.io/gorm"
)

var ErrTodoNotFound = errors.New("not found")

type TodoRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewTodoRepository(db *gorm.DB, logger *slog.Logger) *TodoRepository {
	return &TodoRepository{
		db:     db,
		logger: logger,
	}
}

func (r *TodoRepository) GetAllTodos(ctx context.Context) ([]models.TodoItem, error) {
	return gorm.G[models.TodoItem](r.db).Find(ctx)
}

func (r *TodoRepository) GetTodo(ctx context.Context, id string) (models.TodoItem, error) {

	todoItem, err := gorm.G[models.TodoItem](r.db).Where("id = ?", id).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.TodoItem{}, ErrTodoNotFound
		}

		return models.TodoItem{}, err
	}

	return todoItem, nil
}

func (r *TodoRepository) AddTodo(ctx context.Context, item schemas.TodoCreate) (models.TodoItem, error) {
	todoItem := models.TodoItem{
		Title:   item.Title,
		Content: item.Content,
	}

	if err := gorm.G[models.TodoItem](r.db).Create(ctx, &todoItem); err != nil {
		return models.TodoItem{}, nil
	}

	return todoItem, nil

}

func (r *TodoRepository) DeleteTodo(ctx context.Context, id string) error {
	rowsAffected, err := gorm.G[models.TodoItem](r.db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTodoNotFound
	}

	return nil
}

func (r *TodoRepository) PatchTodo(ctx context.Context, id string, item schemas.TodoPatchUpdate) (models.TodoItem, error) {

	var todoItem models.TodoItem

	if item.Content != nil {
		todoItem.Content = *item.Content
	}

	if item.Title != nil {
		todoItem.Title = *item.Title
	}

	rowsAffected, err := gorm.G[models.TodoItem](r.db).Where("id = ?", id).Updates(ctx, todoItem)
	if err != nil {
		return models.TodoItem{}, err
	}

	if rowsAffected == 0 {
		return models.TodoItem{}, ErrTodoNotFound
	}

	return r.GetTodo(ctx, id)
}
