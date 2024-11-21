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

func (s *PlantService) GetPlants() ([]types.PlantWithDates, error) {
	query := `
        WITH LastWatering AS (
            SELECT plant_id, MAX(entry_date) as last_watered_at   -- Changed column alias
            FROM journal_entries
            WHERE entry_type = 'Watering'
            GROUP BY plant_id
        ),
        LastFertilizing AS (
            SELECT plant_id, MAX(entry_date) as last_fertilized_at  -- Changed column alias
            FROM journal_entries
            WHERE entry_type = 'Fertilizing'
            GROUP BY plant_id
        )
        SELECT p.id,
               p.name,
               p.species,
               p.health,
               p.growth_stage,
               p.planting_date,
               p.image_path,
               p.notes,
               p.created_at,
               p.updated_at,
               p.deleted_at,
               lw.last_watered_at,      -- Matches struct tag
               lf.last_fertilized_at    -- Matches struct tag
        FROM plants p
        LEFT JOIN LastWatering lw ON p.id = lw.plant_id
        LEFT JOIN LastFertilizing lf ON p.id = lf.plant_id
        WHERE p.deleted_at IS NULL
        ORDER BY p.created_at DESC
    `
	var plants []types.PlantWithDates
	err := s.db.Select(&plants, query)
	if err != nil {
		return nil, fmt.Errorf("error fetching plants: %w", err)
	}
	return plants, nil
}

func (s *PlantService) GetPlant(id int) (*types.PlantWithDates, error) {
	query := `
        WITH LastWatering AS (
            SELECT plant_id, entry_date as last_watered_at
            FROM journal_entries je1
            WHERE entry_type = 'Watering'
            AND entry_date = (
                SELECT MAX(entry_date)
                FROM journal_entries je2
                WHERE je2.plant_id = je1.plant_id
                AND entry_type = 'Watering'
            )
        ),
        LastFertilizing AS (
            SELECT plant_id, entry_date as last_fertilized_at
            FROM journal_entries je1
            WHERE entry_type = 'Fertilizing'
            AND entry_date = (
                SELECT MAX(entry_date)
                FROM journal_entries je2
                WHERE je2.plant_id = je1.plant_id
                AND entry_type = 'Fertilizing'
            )
        )
        SELECT p.*, 
               lw.last_watered_at,
               lf.last_fertilized_at
        FROM plants p
        LEFT JOIN LastWatering lw ON p.id = lw.plant_id
        LEFT JOIN LastFertilizing lf ON p.id = lf.plant_id
        WHERE p.id = $1 AND p.deleted_at IS NULL
    `
	var plant types.PlantWithDates
	err := s.db.Get(&plant, query, id)
	if err != nil {
		return nil, err
	}
	return &plant, nil
}

func (s *PlantService) CreatePlant(plant *types.PlantWithDates) error {
	query := `
        INSERT INTO plants (
            name, species, health, growth_stage, planting_date, 
            image_path, notes, is_cross, generation
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, created_at, updated_at
    `

	// Set generation to NULL if not a cross
	if !plant.IsCross {
		plant.Generation = sql.NullString{}
	}

	return s.db.QueryRow(
		query,
		plant.Name,
		plant.Species,
		plant.Health,
		plant.GrowthStage,
		plant.PlantingDate,
		plant.ImagePath,
		plant.Notes,
		plant.IsCross,
		plant.Generation,
	).Scan(&plant.ID, &plant.CreatedAt, &plant.UpdatedAt)
}

func (s *PlantService) UpdatePlant(plant *types.PlantWithDates) error {
	query := `
        UPDATE plants 
        SET name = $1, species = $2, health = $3, growth_stage = $4,
            image_path = $5, notes = $6, planting_date = $7, is_cross = $8,
            generation = $9, updated_at = CURRENT_TIMESTAMP
        WHERE id = $10 AND deleted_at IS NULL`

	// Set generation to NULL if not a cross
	if !plant.IsCross {
		plant.Generation = sql.NullString{}
	}

	result, err := s.db.Exec(
		query,
		plant.Name,
		plant.Species,
		plant.Health,
		plant.GrowthStage,
		plant.ImagePath,
		plant.Notes,
		plant.PlantingDate,
		plant.IsCross,
		plant.Generation,
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

func (s *PlantService) GetPlantsWithLastDates() ([]types.PlantWithDates, error) {
	query := `
        WITH LastWatering AS (
            SELECT plant_id, entry_date as last_watered_at
            FROM journal_entries je1
            WHERE entry_type = 'Watering'
            AND entry_date = (
                SELECT MAX(entry_date)
                FROM journal_entries je2
                WHERE je2.plant_id = je1.plant_id
                AND entry_type = 'Watering'
            )
        ),
        LastFertilizing AS (
            SELECT plant_id, entry_date as last_fertilized_at
            FROM journal_entries je1
            WHERE entry_type = 'Fertilizing'
            AND entry_date = (
                SELECT MAX(entry_date)
                FROM journal_entries je2
                WHERE je2.plant_id = je1.plant_id
                AND entry_type = 'Fertilizing'
            )
        )
        SELECT p.*, 
               lw.last_watered_at,
               lf.last_fertilized_at,
               p.is_cross,
               p.generation
        FROM plants p
        LEFT JOIN LastWatering lw ON p.id = lw.plant_id
        LEFT JOIN LastFertilizing lf ON p.id = lf.plant_id
        WHERE p.deleted_at IS NULL
        ORDER BY p.created_at DESC
    `
	var plants []types.PlantWithDates
	err := s.db.Select(&plants, query)
	if err != nil {
		return nil, fmt.Errorf("error fetching plants with dates: %w", err)
	}
	return plants, nil
}

func (s *PlantService) DeleteJournalEntry(plantID, entryID int) error {
	query := `
        DELETE FROM journal_entries 
        WHERE id = $1 AND plant_id = $2
    `
	result, err := s.db.Exec(query, entryID, plantID)
	if err != nil {
		return fmt.Errorf("error deleting journal entry: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("journal entry not found")
	}
	return nil
}

func (s *PlantService) GetJournalEntry(entryID, plantID int) (*types.JournalEntry, error) {
	var entry types.JournalEntry
	query := `
        SELECT * FROM journal_entries 
        WHERE id = $1 
        AND plant_id = $2 
        AND deleted_at IS NULL
    `
	err := s.db.Get(&entry, query, entryID, plantID)
	if err != nil {
		return nil, fmt.Errorf("error getting journal entry: %w", err)
	}
	return &entry, nil
}

func (s *PlantService) UpdateJournalEntry(entry *types.JournalEntry) error {
	query := `
        UPDATE journal_entries 
        SET title = $1, 
            entry_type = $2, 
            description = $3, 
            entry_date = $4,
            image_path = COALESCE($5, image_path),
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $6 AND plant_id = $7
        RETURNING created_at, updated_at
    `

	return s.db.QueryRow(
		query,
		entry.Title,
		entry.EntryType,
		entry.Description,
		entry.EntryDate,
		entry.ImagePath,
		entry.ID,
		entry.PlantID,
	).Scan(&entry.CreatedAt, &entry.UpdatedAt)
}
