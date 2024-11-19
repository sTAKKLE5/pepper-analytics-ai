package types

import (
	"fmt"
	"time"
)

type PlantHealth string

const (
	PlantHealthExcellent PlantHealth = "Excellent"
	PlantHealthGood      PlantHealth = "Good"
	PlantHealthFair      PlantHealth = "Fair"
	PlantHealthPoor      PlantHealth = "Poor"
)

type GrowthStage string

const (
	GrowthStageSeedling   GrowthStage = "Seedling"
	GrowthStageVegetative GrowthStage = "Vegetative"
	GrowthStageFlowering  GrowthStage = "Flowering"
	GrowthStageFruiting   GrowthStage = "Fruiting"
)

type Plant struct {
	ID           int         `db:"id"`
	Name         string      `db:"name"`
	Species      string      `db:"species"`
	Health       PlantHealth `db:"health"`
	GrowthStage  GrowthStage `db:"growth_stage"`
	PlantingDate time.Time   `db:"planting_date"`
	LastWatered  *time.Time  `db:"last_watered_at"`
	ImagePath    string      `db:"image_path"`
	Notes        string      `db:"notes"`
	DeletedAt    *time.Time  `db:"deleted_at"`
	CreatedAt    time.Time   `db:"created_at"`
	UpdatedAt    time.Time   `db:"updated_at"`
}

type JournalEntry struct {
	ID          int        `db:"id"`
	PlantID     int        `db:"plant_id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	EntryType   string     `db:"entry_type"`
	ImagePath   string     `db:"image_path"`
	DeletedAt   *time.Time `db:"deleted_at"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

func ParsePlantHealth(s string) (PlantHealth, error) {
	switch s {
	case "Excellent":
		return PlantHealthExcellent, nil
	case "Good":
		return PlantHealthGood, nil
	case "Fair":
		return PlantHealthFair, nil
	case "Poor":
		return PlantHealthPoor, nil
	default:
		return "", fmt.Errorf("invalid plant health value: %s", s)
	}
}

func ParseGrowthStage(s string) (GrowthStage, error) {
	switch s {
	case "Seedling":
		return GrowthStageSeedling, nil
	case "Vegetative":
		return GrowthStageVegetative, nil
	case "Flowering":
		return GrowthStageFlowering, nil
	case "Fruiting":
		return GrowthStageFruiting, nil
	default:
		return "", fmt.Errorf("invalid growth stage value: %s", s)
	}
}
