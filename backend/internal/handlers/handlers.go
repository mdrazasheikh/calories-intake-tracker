package handlers

import (
	"net/http"
	"time"

	"calories-intake-tracker/backend/internal/models"
	"calories-intake-tracker/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	food *services.FoodService
	goal *services.GoalService
	rec  *services.RecommendationService
}

func NewHandler(food *services.FoodService, goal *services.GoalService, rec *services.RecommendationService) *Handler {
	return &Handler{food: food, goal: goal, rec: rec}
}

func (h *Handler) Health(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) }

func (h *Handler) SearchFood(c *gin.Context) {
	q := c.Query("q")
	foods, err := h.food.SearchFoods(q)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, foods)
}

func (h *Handler) LookupBarcode(c *gin.Context) {
	item, err := h.food.LookupBarcode(c.Param("barcode"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *Handler) CreateEntry(c *gin.Context) {
	var e models.IntakeEntry
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if e.ConsumedAt.IsZero() {
		e.ConsumedAt = time.Now().UTC()
	}
	if err := h.goal.SaveEntry(e); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, e)
}

func (h *Handler) ListEntries(c *gin.Context) {
	userID := c.Query("user_id")
	day := time.Now().UTC()
	if d := c.Query("day"); d != "" {
		if parsed, err := time.Parse("2006-01-02", d); err == nil {
			day = parsed
		}
	}
	rows, err := h.goal.ListEntries(userID, day)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *Handler) SetDailyGoal(c *gin.Context) {
	var g models.DailyGoal
	if err := c.ShouldBindJSON(&g); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.goal.SaveDailyGoal(g); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, g)
}

func (h *Handler) SetWeightGoal(c *gin.Context) {
	var g models.WeightGoal
	if err := c.ShouldBindJSON(&g); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	goalType, days, delta, err := services.EstimateWeightGoal(g.CurrentWeightKg, g.TargetWeightKg, g.WeeklyChangeKg)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	g.GoalType, g.EstimatedDaysGoal, g.TargetDailyKCal = goalType, days, delta
	if err := h.goal.SaveWeightGoal(g); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, g)
}

func (h *Handler) CalculateRecommendation(c *gin.Context) {
	var req models.RecommendationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, h.rec.Calculate(req))
}
