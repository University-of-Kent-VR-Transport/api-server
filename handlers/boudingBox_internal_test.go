package handlers

import (
	"reflect"
	"server/types"
	"testing"
)

func Test_parseCoordinate(t *testing.T) {
	type args struct {
		coordinate string
	}
	tests := []struct {
		name    string
		args    args
		want    types.Coordinate
		wantErr bool
	}{
		{
			name: "Parses the coordinate \"1.234,1.632\"",
			args: args{
				coordinate: "1.234,1.632",
			},
			want: types.Coordinate{
				Longitude: 1.234,
				Latitude: 1.632,
			},
			wantErr: false,
		},
		{
			name: "Fails to parses when too many arguments are passed \"1.234,1.632,2.353\"",
			args: args{
				coordinate: "1.234,1.632,2.353",
			},
			want: types.Coordinate{},
			wantErr: true,
		},
		{
			name: "Fails to parses a incorrectly formated longitude \"dang,1.632\"",
			args: args{
				coordinate: "dang,1.632",
			},
			want: types.Coordinate{},
			wantErr: true,
		},
		{
			name: "Fails to parses a incorrectly formated latitude \"1.234,fails\"",
			args: args{
				coordinate: "1.234,fails",
			},
			want: types.Coordinate{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCoordinate(tt.args.coordinate)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCoordinate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCoordinate() = %v, want %v", got, tt.want)
			}
		})
	}
}
