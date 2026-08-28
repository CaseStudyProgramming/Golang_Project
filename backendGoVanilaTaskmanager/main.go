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
	// Sementara hardcode dulu (nanti kita load dari config.yaml)
	db := config.NewPostgresDB("localhost", 5432, "postgres", "berjuang02", "taskmanager", "disable")
	defer db.Close()

	taskModel := models.NewTaskModel(db)
	taskService := services.NewTaskService(taskModel)
	taskController := controllers.NewTaskController(taskService)

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, taskController)

	log.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
