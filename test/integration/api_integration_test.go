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
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestRouter() (*gin.Engine, *gorm.DB) {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&domain.Device{})

	repo := repository.NewPostgresRepository(db)
	svc := service.NewDeviceService(repo)
	h := handler.NewDeviceHandler(svc)

	r := gin.Default()
	handler.RegisterRoutes(r, h)
	return r, db
}

func TestCreateAndGetDevice(t *testing.T) {
	r, _ := setupTestRouter()

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
    
    id := "lifecycle-1"
    reqBody, _ := json.Marshal(handler.CreateDeviceRequest{ID: id, Name: "Phone", Brand: "BrandA"})
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/api/v1/devices", bytes.NewBuffer(reqBody))
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusCreated, w.Code)
}
