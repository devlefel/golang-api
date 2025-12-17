package repository

import (
	"device-api/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Save(device *domain.Device) error {
	result := r.db.Create(device)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return domain.ErrDeviceAlreadyExists
		}
		return result.Error
	}
	return nil
}

func (r *PostgresRepository) FindByID(id string) (*domain.Device, error) {
	var device domain.Device
	result := r.db.First(&device, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrDeviceNotFound
		}
		return nil, result.Error
	}
	return &device, nil
}

func (r *PostgresRepository) FindAll() ([]*domain.Device, error) {
	var devices []*domain.Device
	result := r.db.Find(&devices)
	return devices, result.Error
}

func (r *PostgresRepository) FindByBrand(brand string) ([]*domain.Device, error) {
	var devices []*domain.Device
	result := r.db.Where("brand = ?", brand).Find(&devices)
	return devices, result.Error
}

func (r *PostgresRepository) FindByState(state domain.DeviceState) ([]*domain.Device, error) {
	var devices []*domain.Device
	result := r.db.Where("state = ?", state).Find(&devices)
	return devices, result.Error
}

func (r *PostgresRepository) Delete(id string) error {
	result := r.db.Delete(&domain.Device{}, "id = ?", id)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return domain.ErrDeviceNotFound
    }
	return nil
}

func (r *PostgresRepository) Update(device *domain.Device) error {
	result := r.db.Save(device)
	return result.Error
}
