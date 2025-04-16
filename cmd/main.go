package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"todo/internal/db"

	"todo/internal/handler"
	"todo/internal/repo"
	"todo/internal/usecases"
)

func main() {
	database, err := db.Connection()
	if err != nil {
		panic(err)
	}
	repository := repo.NewTaskRepositoryImpl(database)
	service := usecases.NewTaskServiceImpl(repository)
	handlerService := handler.NewTaskHandler(service)

	router := mux.NewRouter()

	router.StrictSlash(true)

	router.HandleFunc("/tasks", handlerService.GetAllTasks).Methods(http.MethodGet)
	router.HandleFunc("/tasks", handlerService.CreateTask).Methods(http.MethodPost)
	router.HandleFunc("/tasks", handlerService.DeleteAllTasks).Methods(http.MethodDelete)
	router.HandleFunc("/tasks/{id}", handlerService.GetTask).Methods(http.MethodGet)
	router.HandleFunc("/tasks/{id}", handlerService.UpdateTask).Methods(http.MethodPatch)
	router.HandleFunc("/tasks/{id}", handlerService.DeleteTask).Methods(http.MethodDelete)
	router.HandleFunc("/hello", handlerService.Hello).Methods(http.MethodGet)

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		return
	}
}
