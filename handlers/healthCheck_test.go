package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"errors"
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

var pingMock func() error

type databaseMock struct{}

func (db databaseMock) Ping() error {
	return pingMock()
}

func Test_testDatabaseConnection(t *testing.T) {
	type args struct {
		pingResult error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Connection to database is valid",
			args: args{
				pingResult: nil,
			},
			want: true,
		},
		{
			name: "Connection to database is invalid",
			args: args{
				pingResult: errors.New(""),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var db databaseMock
			pingMock = func() error {
				return tt.args.pingResult
			}

			if got := testDatabaseConnection(db); got != tt.want {
				t.Errorf("testDatabaseConnection() = %v, want %v", got, tt.want)
			}
		})
	}
}
