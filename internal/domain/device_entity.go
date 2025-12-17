package domain

import (
	"time"
)

type DeviceState string

const (
	DeviceStateAvailable DeviceState = "available"
	DeviceStateInUse     DeviceState = "in-use"
	DeviceStateInactive  DeviceState = "inactive"
)

type Device struct {
	ID        string      `json:"id" gorm:"primaryKey"`
	Name      string      `json:"name"`
	Brand     string      `json:"brand"`
	State     DeviceState `json:"state"`
	CreatedAt time.Time   `json:"created_at"`
}

func NewDevice(id, name, brand string) *Device {
	return &Device{
		ID:        id,
		Name:      name,
		Brand:     brand,
		State:     DeviceStateAvailable,
		CreatedAt: time.Now(),
	}
}

func (d *Device) CanUpdateDetails(newState DeviceState) error {
	if d.State == DeviceStateInUse && (d.State != newState) {
		return nil 
	}
	return nil
}

func (d *Device) UpdateDetails(name, brand string) error {
    if d.State == DeviceStateInUse {
        if d.Name != name || d.Brand != brand {
            return ErrDeviceInUse
        }
    }
    d.Name = name
    d.Brand = brand
    return nil
}

func (d *Device) UpdateState(state DeviceState) {
    d.State = state
}

func (d *Device) CanBeDeleted() error {
	if d.State == DeviceStateInUse {
		return ErrDeviceInUse
	}
	return nil
}
