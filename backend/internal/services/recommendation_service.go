package services

import (
	"fmt"
	"math"

	"calories-intake-tracker/backend/internal/models"
)

type RecommendationService struct{}

func NewRecommendationService() *RecommendationService { return &RecommendationService{} }

func (s *RecommendationService) Calculate(req models.RecommendationRequest) models.RecommendationResponse {
	bmr := 10*req.WeightKg + 6.25*req.HeightCm - 5*float64(req.Age)
	if req.Sex == "male" {
		bmr += 5
	} else {
		bmr -= 161
	}

	multiplier := map[string]float64{"sedentary": 1.2, "light": 1.375, "moderate": 1.55, "active": 1.725, "very_active": 1.9}[req.ActivityLevel]
	if multiplier == 0 {
		multiplier = 1.2
	}
	maintenance := bmr * multiplier
	recommended := maintenance
	guidance := "Maintain current weight by staying near maintenance calories."

	switch req.GoalType {
	case "lose":
		recommended = math.Max(1200, maintenance-500)
		guidance = "For fat loss, target about 500 kcal/day deficit and monitor weekly trends."
	case "gain":
		recommended = maintenance + 300
		guidance = "For weight gain, target ~300 kcal/day surplus with resistance training."
	}

	return models.RecommendationResponse{BMR: bmr, MaintenanceCalories: maintenance, RecommendedCalories: recommended, GoalType: req.GoalType, Guidance: guidance}
}

func EstimateWeightGoal(currentKg, targetKg, weeklyChange float64) (goalType string, days int, kcal int, err error) {
	if weeklyChange <= 0 {
		return "", 0, 0, fmt.Errorf("weekly_change_kg must be > 0")
	}
	diff := targetKg - currentKg
	goalType = "maintain"
	if diff < 0 {
		goalType = "lose"
	} else if diff > 0 {
		goalType = "gain"
	}
	weeks := math.Abs(diff) / weeklyChange
	days = int(math.Ceil(weeks * 7))
	kcalDelta := int(math.Round((weeklyChange * 7700) / 7))
	if goalType == "lose" {
		kcal = -kcalDelta
	} else if goalType == "gain" {
		kcal = kcalDelta
	}
	return
}
