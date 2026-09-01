package routes

import (
	"fmt"
	"net/http"
	"taskmanager/controllers"
	"taskmanager/middlewares"
)

func RegisterRoutes(mux *http.ServeMux, taskController *controllers.TaskController, authController *controllers.AuthController, categoryController *controllers.CategoryController, tagController *controllers.TagController, subtaskController *controllers.SubtaskController, activityLogController *controllers.ActivityLogController, swaggerController *controllers.SwaggerController, authMiddleware *middlewares.AuthMiddleware) {
	// health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API is runningggg 🚀")
	})

	// swagger documentation endpoints (public)
	mux.HandleFunc("GET /swagger/index.html", swaggerController.ServeSwaggerUI)
	mux.HandleFunc("GET /swagger/openapi.yaml", swaggerController.ServeOpenAPISpec)
	mux.HandleFunc("GET /swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
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
	mux.HandleFunc("GET /tasks/analytics/summary", authMiddleware.Authenticate(taskController.GetAnalyticsSummary))

	// task action endpoints (protected)
	mux.HandleFunc("PATCH /tasks/{id}/complete", authMiddleware.Authenticate(taskController.MarkTaskAsCompleted))
	mux.HandleFunc("PATCH /tasks/{id}/uncomplete", authMiddleware.Authenticate(taskController.MarkTaskAsUncompleted))
	mux.HandleFunc("PATCH /tasks/{id}/restore", authMiddleware.Authenticate(taskController.RestoreDeletedTask))

	// task-tag relationship endpoints (protected)
	mux.HandleFunc("POST /tasks/{id}/tags", authMiddleware.Authenticate(taskController.AddTagToTask))
	mux.HandleFunc("DELETE /tasks/{id}/tags/{tagId}", authMiddleware.Authenticate(taskController.RemoveTagFromTask))
	mux.HandleFunc("GET /tasks/{id}/tags", authMiddleware.Authenticate(tagController.GetTagsByTaskID))

	// subtask endpoints (protected)
	mux.HandleFunc("POST /tasks/{id}/subtasks", authMiddleware.Authenticate(subtaskController.CreateSubtask))
	mux.HandleFunc("GET /tasks/{id}/subtasks", authMiddleware.Authenticate(subtaskController.GetSubtasksByTaskID))
	mux.HandleFunc("GET /subtasks/{id}", authMiddleware.Authenticate(subtaskController.GetSubtaskByID))
	mux.HandleFunc("PUT /subtasks/{id}", authMiddleware.Authenticate(subtaskController.UpdateSubtask))
	mux.HandleFunc("DELETE /subtasks/{id}", authMiddleware.Authenticate(subtaskController.DeleteSubtask))
	mux.HandleFunc("PATCH /subtasks/{id}/toggle", authMiddleware.Authenticate(subtaskController.ToggleSubtask))

	// activity log endpoints (protected)
	mux.HandleFunc("GET /activity-logs", authMiddleware.Authenticate(activityLogController.GetUserActivityLogs))
	mux.HandleFunc("GET /activity-logs/{id}", authMiddleware.Authenticate(activityLogController.GetActivityLogByID))
	mux.HandleFunc("GET /tasks/{id}/activity-logs", authMiddleware.Authenticate(activityLogController.GetTaskActivityLogs))
}
