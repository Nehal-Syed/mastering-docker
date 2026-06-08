package routes

import (
	"database/sql"
	"net/http"

	"mastering-docker/internal/config"
	"mastering-docker/internal/handlers"
	"mastering-docker/internal/middleware"
	"mastering-docker/internal/services"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRoutes(db *sql.DB, cfg *config.Config) *mux.Router {
	// Initialize service and handler
	productService := services.NewProductService(db)
	productHandler := handlers.NewProductHandler(productService)

	// Create router with middleware
	router := mux.NewRouter()
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.RecoveryMiddleware)
	router.Use(middleware.InstanceMiddleware)
	router.Use(middleware.MetricsMiddleware)

	// API routes
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/products", productHandler.GetAllProducts).Methods("GET")
	api.HandleFunc("/products", productHandler.CreateProduct).Methods("POST")
	api.HandleFunc("/products/{id:[0-9]+}", productHandler.GetProduct).Methods("GET")
	api.HandleFunc("/products/{id:[0-9]+}", productHandler.UpdateProduct).Methods("PUT")
	api.HandleFunc("/products/{id:[0-9]+}", productHandler.DeleteProduct).Methods("DELETE")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Metrics endpoint for Prometheus
	router.Handle("/metrics", promhttp.Handler())

	return router
}
