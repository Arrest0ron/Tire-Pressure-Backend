package repository

import (
	"errors"
	"fmt"
	"time"
	"web_backend/internal/app/ds"
	"web_backend/internal/app/utils"

	"gorm.io/gorm"
)

func (r *Repository) GetTirePressureCount(creatorID uint) int64 {
	var appID uint
	err := r.db.Model(&ds.TirePressure{}).
		Where("creator_id = ? AND status = ?", creatorID, "draft").
		Select("id").First(&appID).Error
	if err != nil {
		return 0
	}

	var count int64
	err = r.db.Model(&ds.TirePressureEntry{}).
		Where("tire_pressure_id = ?", appID).Count(&count).Error
	if err != nil {
		return 0
	}
	return count
}

func (r *Repository) GetActiveTirePressureID(creatorID uint) uint {
	var appID uint
	err := r.db.Model(&ds.TirePressure{}).
		Where("creator_id = ? AND status = ?", creatorID, "draft").
		Select("id").First(&appID).Error
	if err != nil {
		return 0
	}
	return appID
}

func (r *Repository) GetTirePressure(id int, creatorID uint) ([]ds.TirePressureEntry, error) {
	var app ds.TirePressure
	err := r.db.Where("id = ? AND creator_id = ? AND status != ?",
		id, creatorID, "deleted").First(&app).Error
	if err != nil {
		return nil, err
	}

	// If request is deleted, return empty items
	if app.Status == "deleted" {
		return []ds.TirePressureEntry{}, nil
	}

	var items []ds.TirePressureEntry
	err = r.db.Where("tire_pressure_id = ?", id).
		Preload("Tire").Preload("TirePressure").
		Order("id ASC").Find(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *Repository) AddTireToPressure(tireID int, creatorID uint) error {
	// Check if user already has a draft
	var existingApp ds.TirePressure
	err := r.db.Where("creator_id = ? AND status = ?", creatorID, "draft").
		First(&existingApp).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new draft if none exists
		app := ds.TirePressure{
			Status:     "draft",
			DateCreate: time.Now(),
			CreatorID:  creatorID,
		}
		if err := r.db.Create(&app).Error; err != nil {
			return err
		}

		// Add tire to the new draft with 0 pressure
		item := ds.TirePressureEntry{
			TirePressureID: app.ID,
			TireID:         tireID,
			CoatingCoeff:   1.0,
			Pressure:       0.0,
		}
		if err := r.db.Create(&item).Error; err != nil {
			return err
		}
		return nil
	} else if err != nil {
		return err
	}

	// If draft exists, add tire to existing draft
	var count int64
	r.db.Model(&ds.TirePressureEntry{}).
		Where("tire_pressure_id = ? AND tire_id = ?", existingApp.ID, tireID).
		Count(&count)

	if count == 0 {
		item := ds.TirePressureEntry{
			TirePressureID: existingApp.ID,
			TireID:         tireID,
			CoatingCoeff:   1.0,
			Pressure:       0.0,
		}
		if err := r.db.Create(&item).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) DeleteTirePressure(appID uint) error {
	query := `
		UPDATE tire_pressures
		SET status = 'deleted'
		WHERE id = $1;
	`
	result := r.db.Exec(query, appID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("tire_pressure with id %d not found", appID)
	}
	return nil
}

func (r *Repository) UpdateCoatingCoeff(appID, tireID uint, coeff float64) error {
	if err := r.db.Model(&ds.TirePressureEntry{}).
		Where("tire_pressure_id = ? AND tire_id = ?", appID, tireID).
		Update("coating_coeff", coeff).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdatePressure(appID, tireID uint, pressure float64) error {
	if err := r.db.Model(&ds.TirePressureEntry{}).
		Where("tire_pressure_id = ? AND tire_id = ?", appID, tireID).
		Update("pressure", pressure).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) CalculateTotalPressure(appID uint) (float64, error) {
	var total float64
	err := r.db.Model(&ds.TirePressureEntry{}).
		Where("tire_pressure_id = ?", appID).
		Select("sum(pressure)").Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *Repository) CompleteTirePressure(appID uint) error {
	// Calculate pressure for all tires using their coating coefficients
	var entries []ds.TirePressureEntry
	err := r.db.Where("tire_pressure_id = ?", appID).Find(&entries).Error
	if err != nil {
		return err
	}

	// Calculate pressure for each tire
	for _, entry := range entries {
		pressure := utils.CalculatePressure(entry.CoatingCoeff)
		if err := r.db.Model(&ds.TirePressureEntry{}).
			Where("id = ?", entry.ID).
			Update("pressure", pressure).Error; err != nil {
			return err
		}
	}

	// Calculate total pressure
	total, err := r.CalculateTotalPressure(appID)
	if err != nil {
		return err
	}

	// Update status and total pressure
	query := `
		UPDATE tire_pressures
		SET status = 'completed',
			total_pressure = $2,
			date_update = CURRENT_TIMESTAMP
		WHERE id = $1;
	`
	result := r.db.Exec(query, appID, total)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("tire_pressure with id %d not found", appID)
	}
	return nil
}

func (r *Repository) DeleteTireEntry(appID, tireID uint) error {
	if err := r.db.Where("tire_pressure_id = ? AND tire_id = ?", appID, tireID).
		Delete(&ds.TirePressureEntry{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) IsDraftTirePressure(appID int, creatorID uint) (bool, error) {
	var app ds.TirePressure
	err := r.db.Select("status").Where("id = ? AND creator_id = ?",
		appID, creatorID).First(&app).Error
	if err != nil {
		return false, err
	}
	return app.Status == "draft", nil
}

func (r *Repository) GetTirePressures(creatorID uint) ([]ds.TirePressure, error) {
	var applications []ds.TirePressure
	err := r.db.Where("creator_id = ? AND status != ?", creatorID, "deleted").
		Order("date_create DESC").Find(&applications).Error
	if err != nil {
		return nil, err
	}
	return applications, nil
}

func (r *Repository) GetTirePressureByID(id int, creatorID uint) (ds.TirePressure, error) {
	var app ds.TirePressure
	err := r.db.Where("id = ? AND creator_id = ? AND status != ?",
		id, creatorID, "deleted").First(&app).Error
	if err != nil {
		return ds.TirePressure{}, err
	}
	return app, nil
}