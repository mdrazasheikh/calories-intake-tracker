package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	SupabaseURL     string
	SupabaseKey     string
	OpenFoodFacts   string
	USDAAPIEndpoint string
}

func MustLoad() Config {
	_ = godotenv.Load()

	cfg := Config{
		Port:            getEnv("PORT", "8080"),
		SupabaseURL:     getEnv("SUPABASE_URL", ""),
		SupabaseKey:     getEnv("SUPABASE_SERVICE_ROLE_KEY", ""),
		OpenFoodFacts:   getEnv("OPEN_FOOD_FACTS_BASE_URL", "https://world.openfoodfacts.org/api/v2"),
		USDAAPIEndpoint: getEnv("USDA_API_ENDPOINT", "https://api.nal.usda.gov/fdc/v1"),
	}

	if cfg.SupabaseURL == "" || cfg.SupabaseKey == "" {
		log.Println("warning: Supabase credentials missing. data persistence will fail until configured")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
