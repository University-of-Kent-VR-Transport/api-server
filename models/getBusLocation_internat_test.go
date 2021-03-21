package models

import (
	"server/types"
	"testing"
)

func Test_buildRequestURL(t *testing.T) {
	type args struct {
		coordinateTopLeft     types.Coordinate
		coordinateBottomRight types.Coordinate
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "some test",
			args: args{
				coordinateTopLeft: types.Coordinate{
					Longitude: 1.234,
					Latitude: 2.345,
				},
				coordinateBottomRight: types.Coordinate{
					Longitude: 3.456,
					Latitude: 4.567,
				},
			},
			want: "https://data.bus-data.dft.gov.uk/api/v1/datafeed?status=published&api_key=&boundingBox=1.234,4.567,3.456,2.345",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildRequestURL(tt.args.coordinateTopLeft, tt.args.coordinateBottomRight)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildRequestURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("buildRequestURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
