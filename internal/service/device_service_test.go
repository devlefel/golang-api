package service_test

import (
	"device-api/internal/domain"
	"device-api/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of domain.IDeviceRepository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Save(device *domain.Device) error {
	args := m.Called(device)
	return args.Error(0)
}

func (m *MockRepository) FindByID(id string) (*domain.Device, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Device), args.Error(1)
}

func (m *MockRepository) FindAll() ([]*domain.Device, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Device), args.Error(1)
}

func (m *MockRepository) FindByBrand(brand string) ([]*domain.Device, error) {
	args := m.Called(brand)
	return args.Get(0).([]*domain.Device), args.Error(1)
}

func (m *MockRepository) FindByState(state domain.DeviceState) ([]*domain.Device, error) {
	args := m.Called(state)
	return args.Get(0).([]*domain.Device), args.Error(1)
}

func (m *MockRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) Update(device *domain.Device) error {
	args := m.Called(device)
	return args.Error(0)
}

func TestCreateDevice(t *testing.T) {
	t.Run("success", func(t *testing.T) {
        mockRepo := new(MockRepository)
        svc := service.NewDeviceService(mockRepo)
		mockRepo.On("FindByID", "123").Return(nil, domain.ErrDeviceNotFound)
		mockRepo.On("Save", mock.AnythingOfType("*domain.Device")).Return(nil)

		device, err := svc.CreateDevice("123", "Pixel", "Google")
		assert.NoError(t, err)
		assert.Equal(t, "123", device.ID)
		assert.Equal(t, "Pixel", device.Name)
        assert.Equal(t, domain.DeviceStateAvailable, device.State)
		mockRepo.AssertExpectations(t)
	})

	t.Run("already exists", func(t *testing.T) {
        mockRepo := new(MockRepository)
        svc := service.NewDeviceService(mockRepo)
		existing := &domain.Device{ID: "123"}
		mockRepo.On("FindByID", "123").Return(existing, nil)

		_, err := svc.CreateDevice("123", "Pixel", "Google")
		assert.ErrorIs(t, err, domain.ErrDeviceAlreadyExists)
	})
}

func TestUpdateDevice(t *testing.T) {
	t.Run("success_update_details", func(t *testing.T) {
        mockRepo := new(MockRepository)
        svc := service.NewDeviceService(mockRepo)
        existing := domain.NewDevice("123", "Old", "OldBrand")
		mockRepo.On("FindByID", "123").Return(existing, nil)
		mockRepo.On("Update", mock.MatchedBy(func(d *domain.Device) bool {
			return d.Name == "New" && d.Brand == "NewBrand"
		})).Return(nil)

		updated, err := svc.UpdateDevice("123", "New", "NewBrand")
		assert.NoError(t, err)
		assert.Equal(t, "New", updated.Name)
	})

    t.Run("fail_in_use_update_details", func(t *testing.T) {
        mockRepo := new(MockRepository)
        svc := service.NewDeviceService(mockRepo)
        existing := domain.NewDevice("123", "Old", "OldBrand")
        existing.State = domain.DeviceStateInUse
		mockRepo.On("FindByID", "123").Return(existing, nil)
        
        // No update call expected
		_, err := svc.UpdateDevice("123", "New", "NewBrand")
		assert.ErrorIs(t, err, domain.ErrDeviceInUse)
	})
}

func TestDeleteDevice(t *testing.T) {
    t.Run("success_delete", func(t *testing.T) {
        mockRepo := new(MockRepository)
        svc := service.NewDeviceService(mockRepo)
        existing := domain.NewDevice("123", "Old", "OldBrand")
        mockRepo.On("FindByID", "123").Return(existing, nil)
        mockRepo.On("Delete", "123").Return(nil)
        
        err := svc.DeleteDevice("123")
        assert.NoError(t, err)
    })
    
     t.Run("fail_delete_in_use", func(t *testing.T) {
        mockRepo := new(MockRepository)
        svc := service.NewDeviceService(mockRepo)
        existing := domain.NewDevice("123", "Old", "OldBrand")
        existing.State = domain.DeviceStateInUse
        mockRepo.On("FindByID", "123").Return(existing, nil)
        
        err := svc.DeleteDevice("123")
        assert.ErrorIs(t, err, domain.ErrDeviceInUse)
    })
}
