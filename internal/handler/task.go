package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"todo/internal/errors"
	"todo/internal/usecases"
)

type TaskHandler struct {
	service usecases.TaskService
}

func NewTaskHandler(service usecases.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	respTask, err := h.service.CreateTask(input.Title, input.Description)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			http.Error(w, appErr.Message, appErr.Code)
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
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := mux.Vars(r)["id"]
	w.Header().Set("Content-Type", "application/json")
	task, err := h.service.GetTaskById(id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			http.Error(w, appErr.Message, appErr.Code)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := mux.Vars(r)["id"]
	w.Header().Set("Content-Type", "application/json")
	if err := h.service.DeleteTaskById(id); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			http.Error(w, appErr.Message, appErr.Code)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := mux.Vars(r)["id"]
	var input struct {
		Title       string `json:"title,omitempty"`
		Description string `json:"description,omitempty"`
		Status      *bool  `json:"status,omitempty"`
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	currentTask, err := h.service.GetTaskById(id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			http.Error(w, appErr.Message, appErr.Code)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	if input.Title != "" {
		currentTask.Title = input.Title
	}
	if input.Description != "" {
		currentTask.Description = input.Description
	}
	if input.Status != nil {
		currentTask.Status = *input.Status
	}
	updatedTask, err := h.service.UpdateTask(id, currentTask)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			http.Error(w, appErr.Message, appErr.Code)
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
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	tasks, err := h.service.GetAllTasks()
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			http.Error(w, appErr.Message, appErr.Code)
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
	if r.Method != http.MethodDelete {

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := h.service.DeleteAllTasks(); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {

			http.Error(w, appErr.Message, appErr.Code)
		} else {

			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
