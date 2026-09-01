package routes

import (
	"fmt"
	"net/http"
	"taskmanager/controllers"
	"taskmanager/middlewares"
)

func RegisterRoutes(mux *http.ServeMux, taskController *controllers.TaskController, authController *controllers.AuthController, categoryController *controllers.CategoryController, authMiddleware *middlewares.AuthMiddleware) {
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
}
