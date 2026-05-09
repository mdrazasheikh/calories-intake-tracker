package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"calories-intake-tracker/backend/internal/config"
	"calories-intake-tracker/backend/internal/models"
)

type FoodService struct{ cfg config.Config }

func NewFoodService(cfg config.Config) *FoodService { return &FoodService{cfg: cfg} }

func (s *FoodService) SearchFoods(query string) ([]models.FoodItem, error) {
	u := fmt.Sprintf("%s/foods/search?search_terms=%s&fields=product_name,brands,code,nutriments&json=true", s.cfg.OpenFoodFacts, url.QueryEscape(query))
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		Products []struct {
			Name       string `json:"product_name"`
			Brand      string `json:"brands"`
			Code       string `json:"code"`
			Nutriments struct {
				EnergyKcal100g float64 `json:"energy-kcal_100g"`
			} `json:"nutriments"`
		} `json:"products"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	items := make([]models.FoodItem, 0, len(data.Products))
	for _, p := range data.Products {
		if p.Name == "" {
			continue
		}
		items = append(items, models.FoodItem{Name: p.Name, Brand: p.Brand, Barcode: p.Code, Calories: p.Nutriments.EnergyKcal100g, Serving: "100g", Source: "openfoodfacts"})
		if len(items) >= 20 {
			break
		}
	}
	return items, nil
}

func (s *FoodService) LookupBarcode(code string) (models.FoodItem, error) {
	u := fmt.Sprintf("%s/product/%s?fields=product_name,brands,code,nutriments", s.cfg.OpenFoodFacts, url.PathEscape(code))
	resp, err := http.Get(u)
	if err != nil {
		return models.FoodItem{}, err
	}
	defer resp.Body.Close()

	var data struct {
		Product struct {
			Name       string `json:"product_name"`
			Brand      string `json:"brands"`
			Code       string `json:"code"`
			Nutriments struct {
				EnergyKcal100g float64 `json:"energy-kcal_100g"`
			} `json:"nutriments"`
		} `json:"product"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return models.FoodItem{}, err
	}
	if data.Product.Name == "" {
		return models.FoodItem{}, fmt.Errorf("barcode not found")
	}

	return models.FoodItem{
		Name:     data.Product.Name,
		Brand:    data.Product.Brand,
		Barcode:  data.Product.Code,
		Calories: data.Product.Nutriments.EnergyKcal100g,
		Serving:  "100g",
		Source:   "openfoodfacts",
	}, nil
}
