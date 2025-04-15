package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"todo/internal/apperr"
	"todo/internal/entity"
)

type TaskService interface {
	CreateTask(title, description string) (entity.Task, error)
	GetAllTasks() ([]entity.Task, error)
	GetTaskById(taskId string) (entity.Task, error)
	DeleteTaskById(taskId string) error
	DeleteAllTasks() error
	UpdateTask(id string, task entity.Task) (entity.Task, error)
}

type TaskHandler struct {
	service TaskService
}

func NewTaskHandler(service TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {

	// Не должно так быть. Давай делать кодогенерацию и описовать контракты
	// TODO https://github.com/oapi-codegen/oapi-codegen
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateTitle(input.Title, 50); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := validateDescription(input.Description, 500); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	respTask, err := h.service.CreateTask(input.Title, input.Description)
	if err != nil {
		// TODO изучить врапинг и использовать apperr.Is apperr.As
		if errors.Is(err, apperr.DatabaseError) {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		} else {

			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(respTask); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]
	w.Header().Set("Content-Type", "application/json")
	task, err := h.service.GetTaskById(id)
	if err != nil {
		if errors.Is(err, apperr.NotFoundError) {
			http.Error(w, "task wasn't found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]
	w.Header().Set("Content-Type", "application/json")
	if err := h.service.DeleteTaskById(id); err != nil {
		if errors.Is(err, apperr.NotFoundError) {
			http.Error(w, "task wasn't found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]
	var input struct {
		Title       string `json:"title,omitempty"`
		Description string `json:"description,omitempty"`
		Status      string `json:"status,omitempty"`
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	currentTask, err := h.service.GetTaskById(id)
	if err != nil {
		if errors.Is(err, apperr.NotFoundError) {
			http.Error(w, "task wasn't found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}
	if input.Title != "" {
		currentTask.Title = input.Title
	}
	if input.Description != "" {
		currentTask.Description = input.Description
	}
	if input.Status != "" {
		currentTask.Status = input.Status
	}

	if err := validateTitle(input.Title, 50); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := validateDescription(input.Description, 500); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	updatedTask, err := h.service.UpdateTask(id, currentTask)
	if err != nil {
		if errors.Is(err, apperr.DatabaseError) {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedTask); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	tasks, err := h.service.GetAllTasks()
	if err != nil {
		if errors.Is(err, apperr.DatabaseError) {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) DeleteAllTasks(w http.ResponseWriter, r *http.Request) {

	if err := h.service.DeleteAllTasks(); err != nil {
		if errors.Is(err, apperr.DatabaseError) {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)

}

func validateTitle(title string, maxLen int) error {
	if maxLen == 0 {
		maxLen = 100
	}
	if len(title) == 0 {
		return errors.New("bad title")
	}
	if len(title) > maxLen {
		return errors.New("bad title")
	}
	return nil
}

func validateDescription(description string, maxLen int) error {
	if maxLen == 0 {
		maxLen = 100
	}

	if len(description) == 0 {
		return errors.New("bad description")
	}

	if len(description) > maxLen {
		return errors.New("bad description")

	}
	return nil
}
