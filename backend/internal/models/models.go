package models

import "time"

type FoodItem struct {
	Name     string  `json:"name"`
	Brand    string  `json:"brand,omitempty"`
	Barcode  string  `json:"barcode,omitempty"`
	Calories float64 `json:"calories"`
	Serving  string  `json:"serving,omitempty"`
	Source   string  `json:"source"`
}

type IntakeEntry struct {
	UserID      string    `json:"user_id"`
	FoodName    string    `json:"food_name"`
	Calories    float64   `json:"calories"`
	ConsumedAt  time.Time `json:"consumed_at"`
	Quantity    float64   `json:"quantity"`
	ServingSize string    `json:"serving_size"`
}

type DailyGoal struct {
	UserID          string `json:"user_id"`
	TargetCalories  int    `json:"target_calories"`
	ProteinGrams    int    `json:"protein_grams,omitempty"`
	CarbsGrams      int    `json:"carbs_grams,omitempty"`
	FatGrams        int    `json:"fat_grams,omitempty"`
	ActivityLevel   string `json:"activity_level"`
	GoalDescription string `json:"goal_description,omitempty"`
}

type WeightGoal struct {
	UserID            string  `json:"user_id"`
	CurrentWeightKg   float64 `json:"current_weight_kg"`
	TargetWeightKg    float64 `json:"target_weight_kg"`
	WeeklyChangeKg    float64 `json:"weekly_change_kg"`
	GoalType          string  `json:"goal_type"`
	TargetDailyKCal   int     `json:"target_daily_kcal"`
	EstimatedDaysGoal int     `json:"estimated_days_goal"`
}

type RecommendationRequest struct {
	Sex           string  `json:"sex"`
	Age           int     `json:"age"`
	HeightCm      float64 `json:"height_cm"`
	WeightKg      float64 `json:"weight_kg"`
	ActivityLevel string  `json:"activity_level"`
	GoalType      string  `json:"goal_type"`
}

type RecommendationResponse struct {
	BMR                 float64 `json:"bmr"`
	MaintenanceCalories float64 `json:"maintenance_calories"`
	RecommendedCalories float64 `json:"recommended_calories"`
	GoalType            string  `json:"goal_type"`
	Guidance            string  `json:"guidance"`
}
