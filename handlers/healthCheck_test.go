package handlers

import (
	"strings"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	getRequest(t)
	postRequest(t)
}

func getRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept", contentTypeJson)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheck)

	handler.ServeHTTP(responseRecorder, req)

	if status := responseRecorder.Code; status != http.StatusInternalServerError {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusInternalServerError,
		)
	}

	var expectedContentType string = "application/json"
	receivedContentType := responseRecorder.Header().Get("Content-Type")
	if receivedContentType != expectedContentType {
		t.Errorf(
			"handler returned wrong content type header: got %v want %v",
			receivedContentType,
			expectedContentType,
		)
	}

	expected := `{"database": false}`
	if responseRecorder.Body.String() != expected {
		t.Errorf(
			"handler returned unexpected body: got %v want %v",
			responseRecorder.Body.String(),
			expected,
		)
	}
}

func postRequest(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept", contentTypeJson)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheck)

	handler.ServeHTTP(responseRecorder, req)

	if status := responseRecorder.Code; status != http.StatusMethodNotAllowed {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusMethodNotAllowed,
		)
	}

	expectedContentType := []string{http.MethodGet, http.MethodOptions}
	receivedContentType := responseRecorder.Header().Get("Allow")
	if receivedContentType != strings.Join(expectedContentType, ", ") {
		t.Errorf(
			"handler returned wrong allow header: got %v want %v",
			receivedContentType,
			expectedContentType,
		)
	}

	expected := "Method Not Allowed"
	if responseRecorder.Body.String() != expected {
		t.Errorf(
			"handler returned unexpected body: got %v want %v",
			responseRecorder.Body.String(),
			expected,
		)
	}
}