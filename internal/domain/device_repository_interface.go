package domain

type IDeviceRepository interface {
	Save(device *Device) error
	FindByID(id string) (*Device, error)
	FindAll() ([]*Device, error)
	FindByBrand(brand string) ([]*Device, error)
	FindByState(state DeviceState) ([]*Device, error)
	Delete(id string) error
	Update(device *Device) error
}
