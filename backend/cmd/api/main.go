package main

import (
	"log"

	"calories-intake-tracker/backend/internal/config"
	"calories-intake-tracker/backend/internal/handlers"
	"calories-intake-tracker/backend/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.MustLoad()

	foodService := services.NewFoodService(cfg)
	goalService := services.NewGoalService(cfg)
	recommendationService := services.NewRecommendationService()

	h := handlers.NewHandler(foodService, goalService, recommendationService)

	r := gin.Default()
	r.Use(cors.Default())

	api := r.Group("/api/v1")
	{
		api.GET("/health", h.Health)
		api.GET("/foods/search", h.SearchFood)
		api.GET("/foods/barcode/:barcode", h.LookupBarcode)
		api.POST("/entries", h.CreateEntry)
		api.GET("/entries", h.ListEntries)
		api.POST("/goals/daily", h.SetDailyGoal)
		api.POST("/goals/weight", h.SetWeightGoal)
		api.POST("/recommendations/calories", h.CalculateRecommendation)
	}

	log.Printf("server listening on %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
