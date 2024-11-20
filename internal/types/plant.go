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

type Species string

const (
	SpeciesAnnuum        Species = "Capsicum annuum"
	SpeciesChinense      Species = "Capsicum chinense"
	SpeciesBaccatum      Species = "Capsicum baccatum"
	SpeciesFruitescens   Species = "Capsicum frutescens"
	SpeciesPubescens     Species = "Capsicum pubescens"
	SpeciesRhomboideum   Species = "Capsicum rhomboideum"
	SpeciesPraetermissum Species = "Capsicum praetermissum"
	SpeciesCardenasii    Species = "Capsicum cardenasii"
	SpeciesEximium       Species = "Capsicum eximium"
	SpeciesGalapagoense  Species = "Capsicum galapagoense"
	SpeciesTovarii       Species = "Capsicum tovarii"
	SpeciesFlexuosum     Species = "Capsicum flexuosum"
	SpeciesExile         Species = "Capsicum exile"
)

type Plant struct {
	ID             int         `db:"id"`
	Name           string      `db:"name"`
	Species        Species     `db:"species"`
	Health         PlantHealth `db:"health"`
	GrowthStage    GrowthStage `db:"growth_stage"`
	PlantingDate   time.Time   `db:"planting_date"`
	LastWatered    *time.Time  `db:"last_watered_at"`
	LastFertilized *time.Time  `db:"last_fertilized_at"`
	ImagePath      string      `db:"image_path"`
	Notes          string      `db:"notes"`
	DeletedAt      *time.Time  `db:"deleted_at"`
	CreatedAt      time.Time   `db:"created_at"`
	UpdatedAt      time.Time   `db:"updated_at"`
}

type PlantWithDates struct {
	ID              int         `db:"id"`
	Name            string      `db:"name"`
	Species         Species     `db:"species"`
	Health          PlantHealth `db:"health"`
	GrowthStage     GrowthStage `db:"growth_stage"`
	PlantingDate    time.Time   `db:"planting_date"`
	ImagePath       string      `db:"image_path"`
	Notes           string      `db:"notes"`
	CreatedAt       time.Time   `db:"created_at"`
	UpdatedAt       time.Time   `db:"updated_at"`
	DeletedAt       *time.Time  `db:"deleted_at"`
	LastWatering    *time.Time  `db:"last_watered_at"`
	LastFertilizing *time.Time  `db:"last_fertilized_at"`
}

type JournalEntry struct {
	ID          int       `db:"id"`
	PlantID     int       `db:"plant_id"`
	Title       string    `db:"title"`
	EntryType   string    `db:"entry_type"`
	Description string    `db:"description"`
	ImagePath   string    `db:"image_path"`
	EntryDate   time.Time `db:"entry_date"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
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

func ParseSpecies(s string) (Species, error) {
	switch s {
	case "Capsicum annuum":
		return SpeciesAnnuum, nil
	case "Capsicum chinense":
		return SpeciesChinense, nil
	case "Capsicum baccatum":
		return SpeciesBaccatum, nil
	case "Capsicum frutescens":
		return SpeciesFruitescens, nil
	case "Capsicum pubescens":
		return SpeciesPubescens, nil
	case "Capsicum rhomboideum":
		return SpeciesRhomboideum, nil
	case "Capsicum praetermissum":
		return SpeciesPraetermissum, nil
	case "Capsicum cardenasii":
		return SpeciesCardenasii, nil
	case "Capsicum eximium":
		return SpeciesEximium, nil
	case "Capsicum galapagoense":
		return SpeciesGalapagoense, nil
	case "Capsicum tovarii":
		return SpeciesTovarii, nil
	case "Capsicum flexuosum":
		return SpeciesFlexuosum, nil
	case "Capsicum exile":
		return SpeciesExile, nil
	default:
		return "", fmt.Errorf("invalid species value: %s", s)
	}
}
