package routes

import (
	"fmt"
	"net/http"
	"taskmanager/controllers"
)

func RegisterRoutes(mux *http.ServeMux, taskController *controllers.TaskController) {
	// health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API is runningggg 🚀")
	})

	// tasks endpoints
	mux.HandleFunc("GET /tasks", taskController.GetAllTasks)
	mux.HandleFunc("POST /tasks", taskController.CreateTask)
	mux.HandleFunc("GET /tasks/{id}", taskController.GetTaskByID)
	mux.HandleFunc("PUT /tasks/{id}", taskController.UpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", taskController.DeleteTask)

	// task action endpoints
	mux.HandleFunc("PATCH /tasks/{id}/complete", taskController.MarkTaskAsCompleted)
	mux.HandleFunc("PATCH /tasks/{id}/uncomplete", taskController.MarkTaskAsUncompleted)
	mux.HandleFunc("PATCH /tasks/{id}/restore", taskController.RestoreDeletedTask)
}
