package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	getRequest(t)
	postRequest(t)
}

func getRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheckHandler)

	handler.ServeHTTP(responseRecorder, req)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusOK,
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

	expected := `{"alive": true}`
	if responseRecorder.Body.String() != expected {
		t.Errorf(
			"handler returned unexpected body: got %v want %v",
			responseRecorder.Body.String(),
			expected,
		)
	}
}

func postRequest(t *testing.T) {
	req, err := http.NewRequest("POST", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheckHandler)

	handler.ServeHTTP(responseRecorder, req)

	if status := responseRecorder.Code; status != http.StatusMethodNotAllowed {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusMethodNotAllowed,
		)
	}

	var expectedContentType string = http.MethodGet
	receivedContentType := responseRecorder.Header().Get("Allow")
	if receivedContentType != expectedContentType {
		t.Errorf(
			"handler returned wrong allow header: got %v want %v",
			receivedContentType,
			expectedContentType,
		)
	}

	expected := ""
	if responseRecorder.Body.String() != expected {
		t.Errorf(
			"handler returned unexpected body: got %v want %v",
			responseRecorder.Body.String(),
			expected,
		)
	}
}
