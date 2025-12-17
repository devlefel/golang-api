package service

import (
	"device-api/internal/domain"
)

type DeviceService struct {
	repo domain.IDeviceRepository
}

func NewDeviceService(repo domain.IDeviceRepository) *DeviceService {
	return &DeviceService{repo: repo}
}

func (s *DeviceService) CreateDevice(id, name, brand string) (*domain.Device, error) {
	// Requirements check: "Create a new device."
    // Note: ID handling. User didn't specify who provides ID. 
    // Usually client or generated. I'll assume client provides it for this task, or I can generate it.
    // Given the domain object takes ID, I'll pass it through.
    
    // Check if exists
    existing, _ := s.repo.FindByID(id)
    if existing != nil {
        return nil, domain.ErrDeviceAlreadyExists
    }

	device := domain.NewDevice(id, name, brand)
	err := s.repo.Save(device)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (s *DeviceService) GetDevice(id string) (*domain.Device, error) {
	return s.repo.FindByID(id)
}

func (s *DeviceService) ListAllDevices() ([]*domain.Device, error) {
	return s.repo.FindAll()
}

func (s *DeviceService) ListDevicesByBrand(brand string) ([]*domain.Device, error) {
	return s.repo.FindByBrand(brand)
}

func (s *DeviceService) ListDevicesByState(state domain.DeviceState) ([]*domain.Device, error) {
	return s.repo.FindByState(state)
}

func (s *DeviceService) UpdateDevice(id string, name, brand string) (*domain.Device, error) {
	device, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if err := device.UpdateDetails(name, brand); err != nil {
		return nil, err
	}

	if err := s.repo.Update(device); err != nil {
		return nil, err
	}
	return device, nil
}

func (s *DeviceService) UpdateDeviceState(id string, state domain.DeviceState) (*domain.Device, error) {
    device, err := s.repo.FindByID(id)
    if err != nil {
        return nil, err
    }

    // No specific domain rules for simpler state transitions were given other than enum types, 
    // but we should technically have a method on validation?
    // "Creation time cannot be updated" - Handled by not touching it.
    // "Name and brand... cannot be updated if... in use" - Handled in UpdateDevice.
    
    device.UpdateState(state)
    if err := s.repo.Update(device); err != nil {
        return nil, err // Should likely handle concurrent modifications handling if needed, but not asked.
    }
    return device, nil
}


func (s *DeviceService) DeleteDevice(id string) error {
	device, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if err := device.CanBeDeleted(); err != nil {
		return err
	}

	return s.repo.Delete(id)
}
