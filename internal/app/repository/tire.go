package repository

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"metoda/internal/app/ds"
	minio "metoda/internal/app/minioClient"
	"metoda/internal/app/serializer"

	"gorm.io/gorm"
)

// ─── HTML Pages (2 страницы: главная + шина) ────────────────────────────────
// ... остальной код

func (r *Repository) GetAllTires() ([]ds.Tire, error) {
	var tires []ds.Tire
	err := r.db.Where("is_delete = false").Find(&tires).Error
	if err != nil {
		return nil, err
	}
	return tires, nil
}

func (r *Repository) GetTireByID(id int) (*ds.Tire, error) {
	var tire ds.Tire
	err := r.db.Where("tire_id = ? AND is_delete = false", id).First(&tire).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: шина с id %d", ErrNotFound, id)
		}
		return nil, err
	}
	return &tire, nil
}

func (r *Repository) SearchTiresByTitle(title string) ([]ds.Tire, error) {
	var tires []ds.Tire
	err := r.db.Where("tire_title ILIKE ? AND is_delete = ?", "%"+title+"%", false).Find(&tires).Error
	if err != nil {
		return nil, err
	}
	return tires, nil
}

func (r *Repository) CreateTire(j serializer.TireJSON) (ds.Tire, error) {
	if j.TireTitle == "" {
		return ds.Tire{}, fmt.Errorf("поле tire_title обязательно")
	}
	if j.TireMaterialCoefficient == 0 {
		return ds.Tire{}, fmt.Errorf("поле tire_material_coefficient обязательно")
	}
	if j.TireThicknessCoefficient == 0 {
		return ds.Tire{}, fmt.Errorf("поле tire_thickness_coefficient обязательно")
	}
	if j.Description == "" {
		return ds.Tire{}, fmt.Errorf("поле description обязательно")
	}
	t := serializer.TireFromJSON(j)
	if err := r.db.Create(&t).Error; err != nil {
		return ds.Tire{}, err
	}
	return t, nil
}

func (r *Repository) UploadTirePhoto(ctx context.Context, tireID int, file *multipart.FileHeader) (ds.Tire, error) {
	t, err := r.GetTireByID(tireID)
	if err != nil {
		return ds.Tire{}, err
	}
	if err := minio.EnsureBucket(ctx, r.mc, minio.TiresBucket); err != nil {
		return ds.Tire{}, err
	}
	objectName, err := minio.UploadFile(ctx, r.mc, minio.TiresBucket, file)
	if err != nil {
		return ds.Tire{}, err
	}
	t.Photo = objectName
	if err := r.db.Save(t).Error; err != nil {
		return ds.Tire{}, err
	}
	return *t, nil
}

func (r *Repository) UploadTireVideo(ctx context.Context, tireID int, file *multipart.FileHeader) (ds.Tire, error) {
	t, err := r.GetTireByID(tireID)
	if err != nil {
		return ds.Tire{}, err
	}
	if err := minio.EnsureBucket(ctx, r.mc, minio.TiresBucket); err != nil {
		return ds.Tire{}, err
	}
	objectName, err := minio.UploadFile(ctx, r.mc, minio.TiresBucket, file)
	if err != nil {
		return ds.Tire{}, err
	}
	t.Video = objectName
	if err := r.db.Save(t).Error; err != nil {
		return ds.Tire{}, err
	}
	return *t, nil
}
