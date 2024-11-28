package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/iximiuz/todolist/models"
)

type TodoHandler struct {
	rdb *redis.Client
}

func NewTodoHandler(rdb *redis.Client) *TodoHandler {
	return &TodoHandler{rdb: rdb}
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo.ID = uuid.New().String()

	if err := h.rdb.HSet(r.Context(), "todos", todo.ID, todo.Task).Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *TodoHandler) ListTodos(w http.ResponseWriter, r *http.Request) {
	todosMap, err := h.rdb.HGetAll(r.Context(), "todos").Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	todos := []models.Todo{}
	for id, task := range todosMap {
		todos = append(todos, models.Todo{ID: id, Task: task})
	}

	if err := json.NewEncoder(w).Encode(todos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	id := r.URL.Path[len("/todos/"):]
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	if err := h.rdb.HDel(r.Context(), "todos", id).Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
