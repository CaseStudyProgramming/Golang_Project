package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"taskmanager/config"
	"taskmanager/controllers"
	"taskmanager/middlewares"
	"taskmanager/models"
	"taskmanager/routes"
	"taskmanager/services"
	"taskmanager/utils"
)

func main() {
	cfg, err := config.LoadConfig("env/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db := config.NewPostgresDB(
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)
	defer db.Close()

	// Initialize JWT manager
	jwtManager := utils.NewJWTManager(cfg.JWT.Secret)

	// Initialize user components
	userModel := models.NewUserModel(db)
	userService := services.NewUserService(userModel, jwtManager)
	authController := controllers.NewAuthController(userService)

	// Initialize auth middleware
	authMiddleware := middlewares.NewAuthMiddleware(jwtManager)

	// Initialize task components
	taskModel := models.NewTaskModel(db)
	tagModel := models.NewTagModel(db)

	// Initialize activity log components (needed for task service)
	activityLogModel := models.NewActivityLogModel(db)
	activityLogService := services.NewActivityLogService(activityLogModel)

	taskService := services.NewTaskService(taskModel, tagModel, activityLogService)
	taskController := controllers.NewTaskController(taskService)

	// Initialize category components
	categoryModel := models.NewCategoryModel(db)
	categoryService := services.NewCategoryService(categoryModel)
	categoryController := controllers.NewCategoryController(categoryService)

	// Initialize tag components
	tagService := services.NewTagService(tagModel)
	tagController := controllers.NewTagController(tagService)

	// Initialize subtask components
	subtaskModel := models.NewSubtaskModel(db)
	subtaskService := services.NewSubtaskService(subtaskModel, taskModel)
	subtaskController := controllers.NewSubtaskController(subtaskService)

	// Initialize activity log controller
	activityLogController := controllers.NewActivityLogController(activityLogService)

	// Initialize swagger controller
	swaggerController := controllers.NewSwaggerController()

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, taskController, authController, categoryController, tagController, subtaskController, activityLogController, swaggerController, authMiddleware)

	// Apply middlewares in order: Recovery -> Logger -> CORS -> Routes
	handler := middlewares.Recovery(middlewares.Logger(middlewares.CORS(mux)))

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: handler,
	}

	// Graceful shutdown setup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Server running at :%s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")12
}
