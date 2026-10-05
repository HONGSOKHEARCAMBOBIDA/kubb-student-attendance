package repository

import (
	"context"
	"mysql/model"

	"gorm.io/gorm"
)

type LocationRepository interface {
	GetProvince(ctx context.Context) ([]model.Province, error)
	GetDistrict(ctx context.Context, id int) ([]model.District, error)
	GetCommune(ctx context.Context, id int) ([]model.Communce, error)
	GetVillage(ctx context.Context, id int) ([]model.Village, error)
}

type locationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) LocationRepository {
	return &locationRepository{
		db: db,
	}
}

func (r *locationRepository) GetProvince(ctx context.Context) ([]model.Province, error) {
	var data []model.Province
	err := r.db.WithContext(ctx).Find(&data).Error
	return data, err
}

func (r *locationRepository) GetDistrict(ctx context.Context, id int) ([]model.District, error) {
	var data []model.District
	err := r.db.WithContext(ctx).
		Where("province_id = ?", id).
		Find(&data).Error

	return data, err
}

func (r *locationRepository) GetCommune(ctx context.Context, id int) ([]model.Communce, error) {
	var data []model.Communce
	err := r.db.WithContext(ctx).
		Where("district_id = ?", id).
		Find(&data).Error

	return data, err
}

func (r *locationRepository) GetVillage(ctx context.Context, id int) ([]model.Village, error) {
	var data []model.Village
	err := r.db.WithContext(ctx).
		Where("commune_id = ?", id).
		Find(&data).Error

	return data, err
}
