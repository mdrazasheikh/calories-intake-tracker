package services

import (
	"encoding/json"
	"time"

	"calories-intake-tracker/backend/internal/config"
	"calories-intake-tracker/backend/internal/models"
	supabase "github.com/supabase-community/supabase-go"
)

type GoalService struct {
	client *supabase.Client
}

func NewGoalService(cfg config.Config) *GoalService {
	if cfg.SupabaseURL == "" || cfg.SupabaseKey == "" {
		return &GoalService{}
	}
	c, err := supabase.NewClient(cfg.SupabaseURL, cfg.SupabaseKey, nil)
	if err != nil {
		return &GoalService{}
	}
	return &GoalService{client: c}
}

func (s *GoalService) SaveEntry(entry models.IntakeEntry) error {
	if s.client == nil {
		return nil
	}
	_, _, err := s.client.From("intake_entries").Insert(entry, false, "", "", "").Execute()
	return err
}

func (s *GoalService) ListEntries(userID string, day time.Time) ([]models.IntakeEntry, error) {
	if s.client == nil {
		return []models.IntakeEntry{}, nil
	}
	start := day.UTC().Format(time.RFC3339)
	end := day.UTC().Add(24 * time.Hour).Format(time.RFC3339)
	res, _, err := s.client.From("intake_entries").Select("*", "", false).Eq("user_id", userID).Gte("consumed_at", start).Lt("consumed_at", end).Execute()
	if err != nil {
		return nil, err
	}
	var out []models.IntakeEntry
	if err := json.Unmarshal(res, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *GoalService) SaveDailyGoal(goal models.DailyGoal) error {
	if s.client == nil {
		return nil
	}
	_, _, err := s.client.From("daily_goals").Upsert(goal, "", "", "").Execute()
	return err
}

func (s *GoalService) SaveWeightGoal(goal models.WeightGoal) error {
	if s.client == nil {
		return nil
	}
	_, _, err := s.client.From("weight_goals").Upsert(goal, "", "", "").Execute()
	return err
}
