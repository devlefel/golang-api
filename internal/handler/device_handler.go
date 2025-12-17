package handler

import (
	"device-api/internal/domain"
	"device-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	service *service.DeviceService
}

func NewDeviceHandler(s *service.DeviceService) *DeviceHandler {
	return &DeviceHandler{service: s}
}

// CreateDevice godoc
// @Summary Create a new device
// @Description Create a new device with the input payload
// @Tags devices
// @Accept  json
// @Produce  json
// @Param device body CreateDeviceRequest true "Create Device"
// @Success 201 {object} domain.Device
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /devices [post]
func (h *DeviceHandler) CreateDevice(c *gin.Context) {
	var req CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

    // ID is required by domain.NewDevice? The prompt didn't specify if ID is generic or user provided.
    // The previous implementation of Service.CreateDevice assumes ID is passed.
    // If not in request, maybe generate it?
    // Let's assume user provides it or we generate UUID if missing.
    // Spec says: "Id" in Device Domain. 
    // I will add ID to request or generate it.
    // "Fetch a single device" -> by ID.
    // Let's assume for now request has it or we generate it. 
    // Ideally we'd use UUID.
    if req.ID == "" {
         c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
         return
    }

	device, err := h.service.CreateDevice(req.ID, req.Name, req.Brand)
	if err != nil {
		if err == domain.ErrDeviceAlreadyExists {
			c.JSON(http.StatusConflict, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, device)
}

// GetDevice godoc
// @Summary Get a device by ID
// @Description Get details of a single device
// @Tags devices
// @Produce  json
// @Param id path string true "Device ID"
// @Success 200 {object} domain.Device
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /devices/{id} [get]
func (h *DeviceHandler) GetDevice(c *gin.Context) {
	id := c.Param("id")
	device, err := h.service.GetDevice(id)
	if err != nil {
		if err == domain.ErrDeviceNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, device)
}

// ListDevices godoc
// @Summary List all devices
// @Description Get a list of devices, optionally filtered by brand or state
// @Tags devices
// @Produce  json
// @Param brand query string false "Brand filter"
// @Param state query string false "State filter (available, in-use, inactive)"
// @Success 200 {array} domain.Device
// @Failure 500 {object} ErrorResponse
// @Router /devices [get]
func (h *DeviceHandler) ListDevices(c *gin.Context) {
	brand := c.Query("brand")
	state := c.Query("state")

	var devices []*domain.Device
	var err error

	if brand != "" {
		devices, err = h.service.ListDevicesByBrand(brand)
	} else if state != "" {
		devices, err = h.service.ListDevicesByState(domain.DeviceState(state))
	} else {
		devices, err = h.service.ListAllDevices()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, devices)
}

// UpdateDevice godoc
// @Summary Update a device
// @Description Fully or partially update a device (details or state)
// @Tags devices
// @Accept  json
// @Produce  json
// @Param id path string true "Device ID"
// @Param device body UpdateDeviceRequest true "Update Device"
// @Success 200 {object} domain.Device
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /devices/{id} [put]
// @Router /devices/{id} [patch]
func (h *DeviceHandler) UpdateDevice(c *gin.Context) {
    id := c.Param("id")
    var req UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

    var device *domain.Device
    var err error

    if req.State != "" {
        device, err = h.service.UpdateDeviceState(id, domain.DeviceState(req.State))
        if err != nil {
            if err == domain.ErrDeviceNotFound {
                c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
                return
            }
             c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
             return
        }
    }

    if req.Name != "" || req.Brand != "" {
        device, err = h.service.UpdateDevice(id, req.Name, req.Brand)
         if err != nil {
            if err == domain.ErrDeviceNotFound {
                c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
                return
            }
            if err == domain.ErrDeviceInUse {
                 c.JSON(http.StatusUnprocessableEntity, ErrorResponse{Error: err.Error()}) 
                 return
            }
            c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
            return
        }
    }
    
    if req.Name == "" && req.Brand == "" && req.State == "" {
        device, err = h.service.GetDevice(id)
        if err != nil {
             if err == domain.ErrDeviceNotFound {
                c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
                return
            }
            c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
            return
        }
    }

    c.JSON(http.StatusOK, device)
}

// DeleteDevice godoc
// @Summary Delete a device
// @Description Delete a device by ID
// @Tags devices
// @Param id path string true "Device ID"
// @Success 204 "No Content"
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /devices/{id} [delete]
func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	id := c.Param("id")
	err := h.service.DeleteDevice(id)
	if err != nil {
		if err == domain.ErrDeviceNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
        if err == domain.ErrDeviceInUse {
             c.JSON(http.StatusUnprocessableEntity, ErrorResponse{Error: err.Error()}) 
             return
        }
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}


type CreateDeviceRequest struct {
    ID    string `json:"id"`
	Name  string `json:"name" binding:"required"`
	Brand string `json:"brand" binding:"required"`
}

type UpdateDeviceRequest struct {
	Name  string `json:"name"`
	Brand string `json:"brand"`
    State string `json:"state"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
