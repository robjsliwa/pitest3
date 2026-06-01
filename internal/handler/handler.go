package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/robjsliwa/pitest3/internal/model"
	"github.com/robjsliwa/pitest3/internal/storage"
)

// maxBodySize limits request body to 64 KB.
const maxBodySize = 64 * 1024

// Handler processes HTTP requests for Todo operations.
type Handler struct {
	store *storage.Storage
}

// NewHandler creates a new Handler with the given Storage.
func NewHandler(store *storage.Storage) *Handler {
	return &Handler{store: store}
}

// CreateTodo handles POST /todos — creates a new Todo from the JSON body.
func (h *Handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		DueDate     string `json:"due_date"`
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: %v", err)
		return
	}

	todo := model.Todo{
		Title:       input.Title,
		Description: input.Description,
		DueDate:     input.DueDate,
	}

	created, err := h.store.Create(todo)
	if err != nil {
		writeError(w, http.StatusBadRequest, "validation failed: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createResponse{Todo: created}); err != nil {
		writeResponseError(err)
	}
}

type createResponse struct {
	Todo model.Todo `json:"todo"`
}

func writeError(w http.ResponseWriter, status int, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(errorResponse{
		Error:   http.StatusText(status),
		Message: msg,
	}); err != nil {
		writeResponseError(err)
	}
}

type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func writeResponseError(err error) {
	// Client disconnect or similar — nothing actionable, but log for visibility.
	// In production this would use a proper logger.
	_ = err
}
