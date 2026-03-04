package repository

import (
	"web_backend/internal/app/ds"
)

func (r *Repository) GetTires() ([]ds.Tire, error) {
	var tires []ds.Tire
	err := r.db.Find(&tires).Error
	return tires, err
}

func (r *Repository) GetTiresByTitle(title string) ([]ds.Tire, error) {
	var tires []ds.Tire
	err := r.db.Where("title ILIKE ?", "%"+title+"%").Find(&tires).Error
	return tires, err
}

func (r *Repository) GetTire(id int) (*ds.Tire, error) {
	var tire ds.Tire
	err := r.db.First(&tire, id).Error
	if err != nil {
		return nil, err
	}
	return &tire, nil
}