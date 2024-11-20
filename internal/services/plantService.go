package services

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"pepper-analytics-ai/internal/types"
	"time"
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

func (s *PlantService) GetJournalEntries(plantID int) ([]types.JournalEntry, error) {
	var entries []types.JournalEntry
	query := `
        SELECT id, plant_id, title, entry_type, description, image_path, 
               entry_date, created_at, updated_at 
        FROM journal_entries 
        WHERE plant_id = $1 
        ORDER BY entry_date DESC
    `
	err := s.db.Select(&entries, query, plantID)
	if err != nil {
		log.Printf("Error fetching journal entries: %v", err)
		return nil, fmt.Errorf("failed to fetch journal entries: %w", err)
	}
	return entries, nil
}
func (s *PlantService) CreateJournalEntry(entry *types.JournalEntry) error {
	query := `
        INSERT INTO journal_entries (
            plant_id, title, entry_type, description, 
            image_path, entry_date, created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
        )
        RETURNING id, created_at, updated_at
    `

	err := s.db.QueryRow(
		query,
		entry.PlantID,
		entry.Title,
		entry.EntryType,
		entry.Description,
		entry.ImagePath,
		entry.EntryDate,
	).Scan(&entry.ID, &entry.CreatedAt, &entry.UpdatedAt)

	if err != nil {
		log.Printf("Error creating journal entry: %v", err)
		return fmt.Errorf("failed to create journal entry: %w", err)
	}

	return nil
}

func (s *PlantService) GetLastWateringDate(plantID int) (*time.Time, error) {
	var entryDate time.Time
	query := `
        SELECT entry_date 
        FROM journal_entries 
        WHERE plant_id = $1 AND entry_type = 'Watering'
        ORDER BY entry_date DESC 
        LIMIT 1
    `
	err := s.db.Get(&entryDate, query, plantID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entryDate, nil
}

func (s *PlantService) GetLastFertilizingDate(plantID int) (*time.Time, error) {
	var entryDate time.Time
	query := `
        SELECT entry_date 
        FROM journal_entries 
        WHERE plant_id = $1 AND entry_type = 'Fertilizing'
        ORDER BY entry_date DESC 
        LIMIT 1
    `
	err := s.db.Get(&entryDate, query, plantID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entryDate, nil
}
