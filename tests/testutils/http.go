package testutils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// MakeRequest creates and executes an HTTP GET request with optional user ID header
func MakeRequest(router *gin.Engine, userID string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", "/test", nil)
	if userID != "" {
		req.Header.Set("X-User-ID", userID)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// MakeRequestToPath creates and executes an HTTP GET request to a specific path with optional user ID header
func MakeRequestToPath(router *gin.Engine, path, userID string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", path, nil)
	if userID != "" {
		req.Header.Set("X-User-ID", userID)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// AssertStatusCode checks if the HTTP response has the expected status code
func AssertStatusCode(t *testing.T, response *httptest.ResponseRecorder, expectedCode int, context string) {
	t.Helper()
	if response.Code != expectedCode {
		t.Errorf("%s: expected status %d, got %d", context, expectedCode, response.Code)
	}
}

// MakeMultipleRequests executes multiple HTTP requests and returns all responses
func MakeMultipleRequests(router *gin.Engine, userID string, count int) []*httptest.ResponseRecorder {
	responses := make([]*httptest.ResponseRecorder, count)
	for i := 0; i < count; i++ {
		responses[i] = MakeRequest(router, userID)
	}
	return responses
}

func SetupTestRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	for _, middleware := range middlewares {
		router.Use(middleware)
	}

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	return router
}
