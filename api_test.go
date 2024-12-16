package useragentvaluestore

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	api := NewApi()
	return api
}

func TestGet(t *testing.T) {
	router := setupRouter()

	// Test getting a non-existent user-agent
	req, _ := http.NewRequest("GET", VERSION+"/value", nil)
	req.Header.Set("User-Agent", "test-agent")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "userAgent 'test-agent' has nothing stored")

	// Test getting an existing user-agent
	req, _ = http.NewRequest("POST", VERSION+"/value", bytes.NewBufferString("test-value"))
	req.Header.Set("User-Agent", "test-agent")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	req, _ = http.NewRequest("GET", VERSION+"/value", nil)
	req.Header.Set("User-Agent", "test-agent")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "test-value")
}

func TestUpsert(t *testing.T) {
	router := setupRouter()

	// Test upserting a value
	req, _ := http.NewRequest("POST", VERSION+"/value", bytes.NewBufferString("test-value"))
	req.Header.Set("User-Agent", "test-agent")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "test-value")

	// Test updating the value
	req, _ = http.NewRequest("POST", VERSION+"/value", bytes.NewBufferString("new-value"))
	req.Header.Set("User-Agent", "test-agent")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "new-value")
}

func TestDelete(t *testing.T) {
	router := setupRouter()

	// Test deleting a non-existent user-agent
	req, _ := http.NewRequest("DELETE", VERSION+"/value", nil)
	req.Header.Set("User-Agent", "test-agent")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "userAgent 'test-agent' has nothing stored")

	// Test deleting an existing user-agent
	req, _ = http.NewRequest("POST", VERSION+"/value", bytes.NewBufferString("test-value"))
	req.Header.Set("User-Agent", "test-agent")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	req, _ = http.NewRequest("DELETE", VERSION+"/value", nil)
	req.Header.Set("User-Agent", "test-agent")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHomePage(t *testing.T) {
	router := setupRouter()

	// Test home page with no user-agents
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", "test-agent")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"numUserAgents:":0`)

	// Test home page with one user-agent
	req, _ = http.NewRequest("POST", VERSION+"/value", bytes.NewBufferString("test-value"))
	req.Header.Set("User-Agent", "test-agent")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", "test-agent")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"numUserAgents:":1`)
}
