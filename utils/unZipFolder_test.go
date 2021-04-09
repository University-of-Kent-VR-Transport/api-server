package utils

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func Test_UnZipFile(t *testing.T) {
	type args struct {
		zippedFolderPath string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Reads \"Hello World!\" from \"hellWorld.zip\"",
			args: args{
				zippedFolderPath: "testdata/helloWorld.zip",
			},
			want:    []byte("Hello World!\n"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zippedFolder, err := ioutil.ReadFile(tt.args.zippedFolderPath)
			if err != nil {
				t.Fatal("Failed to read zip folder", err)
			}

			got, err := UnZipFile(zippedFolder)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnZipFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			file, err := got[0].Open()
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			fileContent, err := ioutil.ReadAll(file)
			if err != nil {
				t.Fatal("Failed to read UnZipFile byte stream", err)
			}

			if !reflect.DeepEqual(fileContent, tt.want) {
				t.Errorf("UnZipFile() = %v, want %v", got, tt.want)
			}
		})
	}
}