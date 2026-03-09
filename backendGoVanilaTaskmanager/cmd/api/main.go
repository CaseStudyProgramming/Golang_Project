package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"taskmanager/internal/handler"
	"taskmanager/internal/repository"
	"taskmanager/pkg/database"
)

func main() {
	// sementara hardcode dulu (nanti kita load dari config.yaml)
	db := database.NewPostgresDB("localhost", 5432, "postgres", "berjuang02", "taskmanager", "disable")
	defer db.Close()

	TaskRepository := repository.NewTaskRepository(db)
	TaskHandler := handler.NewTaskHandler(TaskRepository)

	// endpoint
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			TaskHandler.CreateTask(w, r)
			return
		}
		if r.Method == "GET" {
			TaskHandler.GetAllTasks(w, r)
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	// endpoint with id
	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			TaskHandler.GetTaskByID(w, r)
			return
		}
		if r.Method == "PUT" {
			TaskHandler.UpdateTask(w, r)
			return
		}
		if r.Method == "DELETE" {
			TaskHandler.DeleteTask(w, r)
			return
		}
		if r.Method == "PATCH" {
			if strings.Contains(r.URL.Path, "/complete") {
				TaskHandler.MarkTaskAsCompleted(w, r)
				return
			}
			if strings.Contains(r.URL.Path, "/uncomplete") {
				TaskHandler.MarkTaskAsUncompleted(w, r)
				return
			}
			if strings.Contains(r.URL.Path, "/restore") {
				TaskHandler.RestoreDeletedTask(w, r)
				return
			}
			http.Error(w, "Invalid endpoint", http.StatusBadRequest)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	// health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_ = r // use in future
		fmt.Fprintln(w, "API is runningggg 🚀")
	})

	log.Println("Server running at :8080") // start server123446
	log.Fatal(http.ListenAndServe(":8080", nil))
}
