package integration_test

import (
	"bytes"
	"device-api/internal/domain"
	"device-api/internal/handler"
	"device-api/internal/repository"
	"device-api/internal/service"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite" // Use sqlite for in-memory integration test
	"gorm.io/gorm"
)

func setupTestRouter() (*gin.Engine, *gorm.DB) {
	// Use in-memory SQLite for integration testing
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&domain.Device{})

	repo := repository.NewPostgresRepository(db) // It works with any Gorm DB
	svc := service.NewDeviceService(repo)
	h := handler.NewDeviceHandler(svc)

	r := gin.Default()
	handler.RegisterRoutes(r, h)
	return r, db
}

func TestCreateAndGetDevice(t *testing.T) {
	r, _ := setupTestRouter()

	// Create
	deviceReq := handler.CreateDeviceRequest{
        ID:    "integration-1",
		Name:  "Integration Device",
		Brand: "Test Brand",
	}
    body, _ := json.Marshal(deviceReq)
	req, _ := http.NewRequest("POST", "/api/v1/devices", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
    
    // Get
    req, _ = http.NewRequest("GET", "/api/v1/devices/integration-1", nil)
    w = httptest.NewRecorder()
    r.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    var device domain.Device
    json.Unmarshal(w.Body.Bytes(), &device)
    assert.Equal(t, "Integration Device", device.Name)
}

func TestDeviceLifecycle(t *testing.T) {
    r, _ := setupTestRouter()
    
    // 1. Create
    id := "lifecycle-1"
    reqBody, _ := json.Marshal(handler.CreateDeviceRequest{ID: id, Name: "Phone", Brand: "BrandA"})
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/api/v1/devices", bytes.NewBuffer(reqBody))
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusCreated, w.Code)
    
    // 2. Update to In-Use
    // Assuming PUT/PATCH to update state? Need to check handlers. 
    // Wait, handlers: UpdateDevice only binds Name/Brand. 
    // We don't have an endpoint exposed to update STATE specifically in the handler implementation!
    // The service has `UpdateDeviceState`, but `DeviceHandler` has `UpdateDevice` which takes `UpdateDeviceRequest` (Name, Brand).
    // I missed exposing state update in the handler!
    // The requirement says "Fully and/or partially update an existing device".
    // I should probably allow state update in `UpdateDevice` or add a specific endpoint. 
    // Given the domain constraints ("Name.. cannot be updated if in use"), being able to set "in-use" is crucial.
    // I will add State to UpdateDeviceRequest and handle it. Or add a separate endpoint.
    // Let's modify the handler and integration test to support state update.
    // But first, let's finish writing this test file assuming I fix the handler.
    
    // Let's assume I fix the handler to accept State.
}
