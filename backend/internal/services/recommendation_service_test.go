package services

import (
	"math"
	"testing"

	"calories-intake-tracker/backend/internal/models"
)

func TestCalculateRecommendation(t *testing.T) {
	tests := []struct {
		name         string
		req          models.RecommendationRequest
		wantBMR      float64
		wantMaint    float64
		wantCalories float64
	}{
		{
			name: "female fat loss uses deficit with activity multiplier",
			req: models.RecommendationRequest{
				Sex:           "female",
				Age:           30,
				HeightCm:      165,
				WeightKg:      65,
				ActivityLevel: "moderate",
				GoalType:      "lose",
			},
			wantBMR:      1370.25,
			wantMaint:    2123.8875,
			wantCalories: 1623.8875,
		},
		{
			name: "male gain adds surplus",
			req: models.RecommendationRequest{
				Sex:           "male",
				Age:           40,
				HeightCm:      180,
				WeightKg:      80,
				ActivityLevel: "light",
				GoalType:      "gain",
			},
			wantBMR:      1730,
			wantMaint:    2378.75,
			wantCalories: 2678.75,
		},
		{
			name: "unknown activity defaults to sedentary maintenance",
			req: models.RecommendationRequest{
				Sex:           "female",
				Age:           25,
				HeightCm:      160,
				WeightKg:      55,
				ActivityLevel: "unknown",
				GoalType:      "maintain",
			},
			wantBMR:      1264,
			wantMaint:    1516.8,
			wantCalories: 1516.8,
		},
	}

	service := NewRecommendationService()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.Calculate(tt.req)
			assertFloat(t, got.BMR, tt.wantBMR)
			assertFloat(t, got.MaintenanceCalories, tt.wantMaint)
			assertFloat(t, got.RecommendedCalories, tt.wantCalories)
			if got.GoalType != tt.req.GoalType {
				t.Fatalf("GoalType = %q, want %q", got.GoalType, tt.req.GoalType)
			}
			if got.Guidance == "" {
				t.Fatal("Guidance should not be empty")
			}
		})
	}
}

func TestCalculateRecommendationLossFloor(t *testing.T) {
	got := NewRecommendationService().Calculate(models.RecommendationRequest{
		Sex:           "female",
		Age:           70,
		HeightCm:      150,
		WeightKg:      40,
		ActivityLevel: "sedentary",
		GoalType:      "lose",
	})

	if got.RecommendedCalories != 1200 {
		t.Fatalf("RecommendedCalories = %v, want floor of 1200", got.RecommendedCalories)
	}
}

func TestEstimateWeightGoal(t *testing.T) {
	tests := []struct {
		name         string
		currentKg    float64
		targetKg     float64
		weeklyChange float64
		wantType     string
		wantDays     int
		wantKcal     int
		wantErr      bool
	}{
		{name: "loss goal", currentKg: 80, targetKg: 75, weeklyChange: 0.5, wantType: "lose", wantDays: 70, wantKcal: -550},
		{name: "gain goal", currentKg: 70, targetKg: 73, weeklyChange: 0.25, wantType: "gain", wantDays: 84, wantKcal: 275},
		{name: "maintain goal", currentKg: 70, targetKg: 70, weeklyChange: 0.5, wantType: "maintain", wantDays: 0, wantKcal: 0},
		{name: "invalid weekly change", currentKg: 70, targetKg: 65, weeklyChange: 0, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, gotDays, gotKcal, err := EstimateWeightGoal(tt.currentKg, tt.targetKg, tt.weeklyChange)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotType != tt.wantType || gotDays != tt.wantDays || gotKcal != tt.wantKcal {
				t.Fatalf("got (%q, %d, %d), want (%q, %d, %d)", gotType, gotDays, gotKcal, tt.wantType, tt.wantDays, tt.wantKcal)
			}
		})
	}
}

func assertFloat(t *testing.T, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 0.0001 {
		t.Fatalf("got %v, want %v", got, want)
	}
}
