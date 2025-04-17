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

	if err := validateLen(input.Title, 50); err != nil {
		http.Error(w, fmt.Sprintf("title %s", err.Error()), http.StatusBadRequest)
	}
	if err := validateLen(input.Description, 500); err != nil {
		http.Error(w, fmt.Sprintf("description %s", err.Error()), http.StatusBadRequest)
	}

	respTask, err := h.service.CreateTask(input.Title, input.Description)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
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
			http.Error(w, "task not found", http.StatusNotFound)
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
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]
	var input struct {
		Title       string            `json:"title,omitempty"`
		Description string            `json:"description,omitempty"`
		Status      entity.TaskStatus `json:"status,omitempty"`
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	currentTask, err := h.service.GetTaskById(id)
	if err != nil {
		if errors.Is(err, apperr.NotFoundError) {
			http.Error(w, "task not found", http.StatusNotFound)
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

	if err := validateLen(input.Title, 50); err != nil {
		http.Error(w, fmt.Sprintf("title %s", err.Error()), http.StatusBadRequest)
	}
	if err := validateLen(input.Description, 500); err != nil {
		http.Error(w, fmt.Sprintf("description %s", err.Error()), http.StatusBadRequest)
	}

	updatedTask, err := h.service.UpdateTask(id, currentTask)
	if err != nil {

		http.Error(w, "Internal server error", http.StatusInternalServerError)
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
		if errors.Is(err, apperr.NotFoundError) {
			http.Error(w, "task wasn't found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)

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

		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}
	w.WriteHeader(http.StatusNoContent)

}

func validateLen(data string, maxLen int) error {
	if len(data) <= 0 {
		return errors.New("can't be equal 0")
	}
	if len(data) > maxLen {
		return errors.New("too long")
	}
	return nil
}
