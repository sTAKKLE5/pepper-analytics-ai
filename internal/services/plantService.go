package services

import (
	"github.com/jmoiron/sqlx"
	"pepper-analytics-ai/internal/types"
)

type PlantService struct {
	db *sqlx.DB
}

func NewPlantService(db *sqlx.DB) *PlantService {
	return &PlantService{db: db}
}

func (s *PlantService) GetPlants() ([]types.Plant, error) {
	var plants []types.Plant
	query := `SELECT * FROM plants WHERE deleted_at IS NULL ORDER BY created_at DESC`
	err := s.db.Select(&plants, query)
	return plants, err
}

func (s *PlantService) GetPlant(id int) (*types.Plant, error) {
	var plant types.Plant
	query := `SELECT * FROM plants WHERE id = $1 AND deleted_at IS NULL`
	err := s.db.Get(&plant, query, id)
	return &plant, err
}

func (s *PlantService) CreatePlant(plant *types.Plant) error {
	query := `
        INSERT INTO plants (name, species, health, growth_stage, planting_date, image_path, notes)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at`

	return s.db.QueryRow(
		query,
		plant.Name,
		plant.Species,
		plant.Health,
		plant.GrowthStage,
		plant.PlantingDate,
		plant.ImagePath,
		plant.Notes,
	).Scan(&plant.ID, &plant.CreatedAt, &plant.UpdatedAt)
}

func (s *PlantService) UpdatePlant(plant *types.Plant) error {
	query := `
        UPDATE plants 
        SET name = $1, species = $2, health = $3, growth_stage = $4,
            image_path = $5, notes = $6, updated_at = CURRENT_TIMESTAMP
        WHERE id = $7 AND deleted_at IS NULL`

	result, err := s.db.Exec(
		query,
		plant.Name,
		plant.Species,
		plant.Health,
		plant.GrowthStage,
		plant.ImagePath,
		plant.Notes,
		plant.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrPlantNotFound
	}
	return nil
}

func (s *PlantService) DeletePlant(id int) error {
	query := `UPDATE plants SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL`
	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrPlantNotFound
	}
	return nil
}
