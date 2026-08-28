package main

import (
	"log"
	"net/http"

	"taskmanager/config"
	"taskmanager/controllers"
	"taskmanager/models"
	"taskmanager/routes"
	"taskmanager/services"
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

	taskModel := models.NewTaskModel(db)
	taskService := services.NewTaskService(taskModel)
	taskController := controllers.NewTaskController(taskService)

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, taskController)

	log.Printf("Server running at :%s", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, mux))
}
