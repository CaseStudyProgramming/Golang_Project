package main

import (
	"fmt"
	"log"
	"net/http"

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

	http.HandleFunc("/tasks", TaskHandler.CreateTask)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API is runningggg 🚀")
	})

	log.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
