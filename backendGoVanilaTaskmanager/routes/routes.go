package routes

import (
	"fmt"
	"net/http"
	"strings"
	"taskmanager/controllers"
)

func RegisterRoutes(mux *http.ServeMux, taskController *controllers.TaskController) {
	// health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API is runningggg 🚀")
	})

	// /tasks endpoint
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			taskController.CreateTask(w, r)
			return
		}
		if r.Method == http.MethodGet {
			taskController.GetAllTasks(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	// /tasks/ endpoint
	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			taskController.GetTaskByID(w, r)
			return
		}
		if r.Method == http.MethodPut {
			taskController.UpdateTask(w, r)
			return
		}
		if r.Method == http.MethodDelete {
			taskController.DeleteTask(w, r)
			return
		}
		if r.Method == http.MethodPatch {
			if strings.Contains(r.URL.Path, "/complete") {
				taskController.MarkTaskAsCompleted(w, r)
				return
			}
			if strings.Contains(r.URL.Path, "/uncomplete") {
				taskController.MarkTaskAsUncompleted(w, r)
				return
			}
			if strings.Contains(r.URL.Path, "/restore") {
				taskController.RestoreDeletedTask(w, r)
				return
			}
			http.Error(w, "Invalid endpoint", http.StatusBadRequest)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})
}
