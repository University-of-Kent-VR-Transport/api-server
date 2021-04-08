package controllers

import (
	"testing"
	"errors"
)

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

			if got := HealthCheck(db); got != tt.want {
				t.Errorf("testDatabaseConnection() = %v, want %v", got, tt.want)
			}
		})
	}
}
