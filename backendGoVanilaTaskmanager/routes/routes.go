package routes

import (
	"fmt"
	"net/http"
	"taskmanager/controllers"
	"taskmanager/middlewares"
)

func RegisterRoutes(mux *http.ServeMux, taskController *controllers.TaskController, authController *controllers.AuthController, categoryController *controllers.CategoryController, tagController *controllers.TagController, authMiddleware *middlewares.AuthMiddleware) {
	// health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API is runningggg 🚀")
	})

	// auth endpoints (public)
	mux.HandleFunc("POST /auth/register", authController.Register)
	mux.HandleFunc("POST /auth/login", authController.Login)
	mux.HandleFunc("POST /auth/logout", authController.Logout)
	mux.HandleFunc("GET /auth/me", authMiddleware.Authenticate(authController.GetCurrentUser))

	// categories endpoints (protected)
	mux.HandleFunc("GET /categories", authMiddleware.Authenticate(categoryController.GetAllCategories))
	mux.HandleFunc("POST /categories", authMiddleware.Authenticate(categoryController.CreateCategory))
	mux.HandleFunc("GET /categories/{id}", authMiddleware.Authenticate(categoryController.GetCategoryByID))
	mux.HandleFunc("PUT /categories/{id}", authMiddleware.Authenticate(categoryController.UpdateCategory))
	mux.HandleFunc("DELETE /categories/{id}", authMiddleware.Authenticate(categoryController.DeleteCategory))

	// tags endpoints (protected)
	mux.HandleFunc("GET /tags", authMiddleware.Authenticate(tagController.GetAllTags))
	mux.HandleFunc("POST /tags", authMiddleware.Authenticate(tagController.CreateTag))
	mux.HandleFunc("GET /tags/{id}", authMiddleware.Authenticate(tagController.GetTagByID))
	mux.HandleFunc("PUT /tags/{id}", authMiddleware.Authenticate(tagController.UpdateTag))
	mux.HandleFunc("DELETE /tags/{id}", authMiddleware.Authenticate(tagController.DeleteTag))
	mux.HandleFunc("GET /tags/{id}/tasks", authMiddleware.Authenticate(tagController.GetTasksByTagID))

	// tasks endpoints (protected)
	mux.HandleFunc("GET /tasks", authMiddleware.Authenticate(taskController.GetAllTasks))
	mux.HandleFunc("POST /tasks", authMiddleware.Authenticate(taskController.CreateTask))
	mux.HandleFunc("GET /tasks/{id}", authMiddleware.Authenticate(taskController.GetTaskByID))
	mux.HandleFunc("PUT /tasks/{id}", authMiddleware.Authenticate(taskController.UpdateTask))
	mux.HandleFunc("DELETE /tasks/{id}", authMiddleware.Authenticate(taskController.DeleteTask))

	// task action endpoints (protected)
	mux.HandleFunc("PATCH /tasks/{id}/complete", authMiddleware.Authenticate(taskController.MarkTaskAsCompleted))
	mux.HandleFunc("PATCH /tasks/{id}/uncomplete", authMiddleware.Authenticate(taskController.MarkTaskAsUncompleted))
	mux.HandleFunc("PATCH /tasks/{id}/restore", authMiddleware.Authenticate(taskController.RestoreDeletedTask))

	// task-tag relationship endpoints (protected)
	mux.HandleFunc("POST /tasks/{id}/tags", authMiddleware.Authenticate(taskController.AddTagToTask))
	mux.HandleFunc("DELETE /tasks/{id}/tags/{tagId}", authMiddleware.Authenticate(taskController.RemoveTagFromTask))
	mux.HandleFunc("GET /tasks/{id}/tags", authMiddleware.Authenticate(tagController.GetTagsByTaskID))
}
