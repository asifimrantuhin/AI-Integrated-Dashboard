package models

import "time"

type PrescriptiveRecommendation struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Module         string    `json:"module"`
	EntityType     string    `json:"entity_type"`
	EntityID       string    `json:"entity_id"`
	RiskLevel      string    `json:"risk_level"`
	Recommendation string    `json:"recommendation"`
	ImpactScore    float64   `json:"impact_score"`
	Metadata       string    `json:"metadata"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ScenarioSimulation struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Module      string    `json:"module"`
	ScenarioKey string    `json:"scenario_key"`
	BaseMetrics string    `json:"base_metrics"`
	Adjustments string    `json:"adjustments"`
	Results     string    `json:"results"`
	CreatedBy   *uint     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
