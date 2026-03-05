package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockSourceService struct {
	sources   []string
	listErr   error
	addErr    error
	deleteErr error
	addedID   string
	deletedID string
}

func (m *mockSourceService) ListSources(_ context.Context) ([]string, error) {
	return m.sources, m.listErr
}

func (m *mockSourceService) AddSource(_ context.Context, id string) error {
	m.addedID = id
	return m.addErr
}

func (m *mockSourceService) DeleteSource(_ context.Context, id string) error {
	m.deletedID = id
	return m.deleteErr
}

func setupSourceRouter(svc *mockSourceService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := NewSourceAPI(svc)
	r.GET("/api/sources", api.listSources)
	r.POST("/api/sources", api.createSource)
	r.DELETE("/api/sources", api.deleteSource)
	return r
}

// --- ListSources ---

func TestListSources_ReturnsSources(t *testing.T) {
	svc := &mockSourceService{sources: []string{"/photos", "/backup"}}
	r := setupSourceRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/sources", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	body := w.Body.String()
	if !strings.Contains(body, "/photos") || !strings.Contains(body, "/backup") {
		t.Errorf("body = %s, want sources listed", body)
	}
}

func TestListSources_Empty(t *testing.T) {
	svc := &mockSourceService{sources: []string{}}
	r := setupSourceRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/sources", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if body := w.Body.String(); body != "[]" {
		t.Errorf("body = %s, want []", body)
	}
}

func TestListSources_ServiceError(t *testing.T) {
	svc := &mockSourceService{listErr: errors.New("fail")}
	r := setupSourceRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/sources", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// --- CreateSource ---

func TestCreateSource_Success(t *testing.T) {
	svc := &mockSourceService{}
	r := setupSourceRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/sources", strings.NewReader(`{"id":"/new/photos"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
	}
	if svc.addedID != "/new/photos" {
		t.Errorf("addedID = %q, want %q", svc.addedID, "/new/photos")
	}
}

func TestCreateSource_MissingID(t *testing.T) {
	svc := &mockSourceService{}
	r := setupSourceRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/sources", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCreateSource_InvalidJSON(t *testing.T) {
	svc := &mockSourceService{}
	r := setupSourceRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/sources", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCreateSource_ServiceError(t *testing.T) {
	svc := &mockSourceService{addErr: errors.New("disk full")}
	r := setupSourceRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/sources", strings.NewReader(`{"id":"/fail"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// --- DeleteSource ---

func TestDeleteSource_Success(t *testing.T) {
	svc := &mockSourceService{}
	r := setupSourceRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/sources", strings.NewReader(`{"id":"/home/photos"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
	if svc.deletedID != "/home/photos" {
		t.Errorf("deletedID = %q, want %q", svc.deletedID, "/home/photos")
	}
}

func TestDeleteSource_MissingID(t *testing.T) {
	svc := &mockSourceService{}
	r := setupSourceRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/sources", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestDeleteSource_ServiceError(t *testing.T) {
	svc := &mockSourceService{deleteErr: errors.New("not found")}
	r := setupSourceRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/sources", strings.NewReader(`{"id":"bad"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}
