package controllers

import (
	"net/http"
	"os"
)

type SwaggerController struct{}

func NewSwaggerController() *SwaggerController {
	return &SwaggerController{}
}

// ServeSwaggerUI serves the Swagger UI HTML page
func (c *SwaggerController) ServeSwaggerUI(w http.ResponseWriter, r *http.Request) {
	// Read the index.html file
	content, err := os.ReadFile("swagger/index.html")
	if err != nil {
		http.Error(w, "Swagger UI not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(content)
}

// ServeOpenAPISpec serves the OpenAPI YAML specification
func (c *SwaggerController) ServeOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	// Read the openapi.yaml file
	content, err := os.ReadFile("swagger/openapi.yaml")
	if err != nil {
		http.Error(w, "OpenAPI spec not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/x-yaml")
	w.Write(content)
}