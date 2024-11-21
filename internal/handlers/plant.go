package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"pepper-analytics-ai/internal/services"
	"pepper-analytics-ai/internal/types"
	"pepper-analytics-ai/templates/pages"
	"strconv"
	"time"
)

type PlantHandler struct {
	plantService *services.PlantService
	fileService  *services.FileService
	uploadDir    string
}

func NewPlantHandler(plantService *services.PlantService, fileService *services.FileService) *PlantHandler {
	return &PlantHandler{
		plantService: plantService,
		fileService:  fileService,
		uploadDir:    "uploads",
	}
}

func (h *PlantHandler) HandlePlantList(c *gin.Context) {
	plants, err := h.plantService.GetPlantsWithLastDates()
	if err != nil {
		log.Printf("Error fetching plants: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := pages.Plant(plants).Render(context.Background(), c.Writer); err != nil {
		log.Printf("Error rendering template: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
}

func (h *PlantHandler) HandleCreatePlant(c *gin.Context) {
	log.Printf("Received form data: %+v", c.Request.Form)

	plantingDate, err := time.Parse("2006-01-02", c.PostForm("planting_date"))
	if err != nil {
		log.Printf("Error parsing planting date: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting date"})
		return
	}

	health, err := types.ParsePlantHealth(c.PostForm("health"))
	if err != nil {
		log.Printf("Error parsing health: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	growthStage, err := types.ParseGrowthStage(c.PostForm("growth_stage"))
	if err != nil {
		log.Printf("Error parsing growth stage: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	species, err := types.ParseSpecies(c.PostForm("species"))
	if err != nil {
		log.Printf("Error parsing species: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse the cross status and generation
	isCross := c.PostForm("cross") == "Yes"
	generation := sql.NullString{}
	if isCross {
		generation = sql.NullString{
			String: c.PostForm("generation"),
			Valid:  true,
		}
	}

	plant := &types.PlantWithDates{
		Name:         c.PostForm("name"),
		Species:      species,
		PlantingDate: plantingDate,
		Health:       health,
		GrowthStage:  growthStage,
		Notes:        c.PostForm("notes"),
		IsCross:      isCross,
		Generation:   generation,
	}

	// Handle image upload
	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()

		filePath := fmt.Sprintf("%s/%s", h.uploadDir, header.Filename)
		if err := h.fileService.SaveFile(file, filePath); err != nil {
			log.Printf("Error saving file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
		plant.ImagePath = filePath
	}

	if err := h.plantService.CreatePlant(plant); err != nil {
		log.Printf("Error creating plant: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create plant"})
		return
	}

	// Set header to trigger modal close
	c.Writer.Header().Set("HX-Trigger", "closeModal")

	// After successful creation, fetch all plants with dates and cross information
	plants, err := h.plantService.GetPlantsWithLastDates() // Make sure this fetches cross and generation info
	if err != nil {
		log.Printf("Error fetching plants: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Writer.Header().Set("Content-Type", "text/html")
	templ.Handler(pages.PlantsGrid(plants)).ServeHTTP(c.Writer, c.Request)
}

func (h *PlantHandler) HandleNewPlantForm(c *gin.Context) {
	component := pages.NewPlantForm()
	_ = component.Render(context.Background(), c.Writer)
}

func (h *PlantHandler) HandleEditPlantForm(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	plant, err := h.plantService.GetPlant(id)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	component := pages.EditPlantForm(*plant)
	_ = component.Render(context.Background(), c.Writer)
}

func (h *PlantHandler) HandleUpdatePlant(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	plant, err := h.plantService.GetPlant(id)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	plant.Name = c.PostForm("name")
	plantingDate, err := time.Parse("2006-01-02", c.PostForm("planting_date"))
	if err != nil {
		log.Printf("Error parsing planting date: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting date"})
		return
	}
	plant.PlantingDate = plantingDate

	if species, err := types.ParseSpecies(c.PostForm("species")); err == nil {
		plant.Species = species
	}

	if health, err := types.ParsePlantHealth(c.PostForm("health")); err == nil {
		plant.Health = health
	}

	if growthStage, err := types.ParseGrowthStage(c.PostForm("growth_stage")); err == nil {
		plant.GrowthStage = growthStage
	}

	plant.Notes = c.PostForm("notes")

	// Update cross status and generation
	plant.IsCross = c.PostForm("cross") == "Yes"
	if plant.IsCross {
		plant.Generation = sql.NullString{
			String: c.PostForm("generation"),
			Valid:  true,
		}
	} else {
		plant.Generation = sql.NullString{}
	}

	// Handle new image upload
	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()

		if plant.ImagePath != "" {
			h.fileService.DeleteFile(plant.ImagePath)
		}

		filePath := fmt.Sprintf("%s/%s", h.uploadDir, header.Filename)
		if err := h.fileService.SaveFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
		plant.ImagePath = filePath
	}

	if err := h.plantService.UpdatePlant(plant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update plant"})
		return
	}

	// Set header to trigger modal close
	c.Writer.Header().Set("HX-Trigger", "closeModal")

	// After successful update, fetch all plants with dates
	plants, err := h.plantService.GetPlants()
	if err != nil {
		log.Printf("Error fetching plants: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Writer.Header().Set("Content-Type", "text/html")
	templ.Handler(pages.PlantsGrid(plants)).ServeHTTP(c.Writer, c.Request)
}

func (h *PlantHandler) HandleDeletePlant(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	plant, err := h.plantService.GetPlant(id)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if plant.ImagePath != "" {
		h.fileService.DeleteFile(plant.ImagePath)
	}

	if err := h.plantService.DeletePlant(id); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Writer.Header().Set("Content-Type", "text/html")
	c.String(http.StatusOK, "")
}

func (h *PlantHandler) HandleJournal(c *gin.Context) {
	plantID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	plant, err := h.plantService.GetPlant(plantID)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	entries, err := h.plantService.GetJournalEntries(plantID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	templ.Handler(pages.Journal(*plant, entries)).ServeHTTP(c.Writer, c.Request)
}

func (h *PlantHandler) HandleCreateJournalEntry(c *gin.Context) {
	plantID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	entryDate, err := time.Parse("2006-01-02", c.PostForm("entry_date"))
	if err != nil {
		log.Printf("Error parsing entry date: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entry date"})
		return
	}

	entry := &types.JournalEntry{
		PlantID:     plantID,
		Title:       c.PostForm("title"),
		EntryType:   c.PostForm("entry_type"),
		Description: c.PostForm("description"),
		EntryDate:   entryDate,
	}

	// Handle image upload if present
	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()

		filePath := fmt.Sprintf("/%s/journal/%s", h.uploadDir, header.Filename)
		if err := h.fileService.SaveFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
		entry.ImagePath = filePath
	}

	if err := h.plantService.CreateJournalEntry(entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create journal entry"})
		return
	}

	c.Writer.Header().Set("Content-Type", "text/html")
	templ.Handler(pages.JournalEntry(*entry)).ServeHTTP(c.Writer, c.Request)
}

func (h *PlantHandler) HandleDeleteJournalEntry(c *gin.Context) {
	plantID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	entryID, err := strconv.Atoi(c.Param("entryId"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.plantService.DeleteJournalEntry(plantID, entryID); err != nil {
		log.Printf("Error deleting journal entry: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *PlantHandler) HandleEditJournalEntry(c *gin.Context) {
	plantID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	entryID, err := strconv.Atoi(c.Param("entryId"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	entry, err := h.plantService.GetJournalEntry(entryID, plantID)
	if err != nil {
		log.Printf("Error getting journal entry: %v", err)
		c.Status(http.StatusNotFound)
		return
	}

	if err := pages.EditJournalEntry(*entry).Render(context.Background(), c.Writer); err != nil {
		log.Printf("Error rendering edit form: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
}
func (h *PlantHandler) HandleUpdateJournalEntry(c *gin.Context) {
	plantID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	entryID, err := strconv.Atoi(c.Param("entryId"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	entryDate, err := time.Parse("2006-01-02", c.PostForm("entry_date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entry date"})
		return
	}

	entry := &types.JournalEntry{
		ID:          entryID,
		PlantID:     plantID,
		Title:       c.PostForm("title"),
		EntryType:   c.PostForm("entry_type"),
		Description: c.PostForm("description"),
		EntryDate:   entryDate,
	}

	// Handle image upload if present
	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()

		filePath := fmt.Sprintf("%s/journal/%s", h.uploadDir, header.Filename)
		if err := h.fileService.SaveFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
		entry.ImagePath = filePath
	}

	if err := h.plantService.UpdateJournalEntry(entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update entry"})
		return
	}

	// Return the updated entry
	c.Writer.Header().Set("Content-Type", "text/html")
	templ.Handler(pages.JournalEntry(*entry)).ServeHTTP(c.Writer, c.Request)
}
