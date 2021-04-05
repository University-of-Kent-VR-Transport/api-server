package handlers

import (
	"strings"
	"io/ioutil"
	"os"
	"reflect"
	"server/models"
	"testing"
	"net/http"
	"errors"

	_ "github.com/lib/pq"
)

func Test_unZipFile(t *testing.T) {
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

			got, err := unZipFile(zippedFolder)
			if (err != nil) != tt.wantErr {
				t.Errorf("unZipFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer got.Close()

			fileContent, err := ioutil.ReadAll(got)
			if err != nil {
				t.Fatal("Failed to read unZipFile byte stream", err)
			}

			if !reflect.DeepEqual(fileContent, tt.want) {
				t.Errorf("unZipFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseXML(t *testing.T) {
	type args struct {
		xmlFilePath string
	}
	tests := []struct {
		name    string
		args    args
		want    []models.BusStop
		wantErr bool
	}{
		{
			name: "Tests a simple data set containing two stop points",
			args: args{
				xmlFilePath: "testdata/simpleNaPTAN.xml",
			},
			want: []models.BusStop{
				models.BusStop{
					ID:        "010000001",
					Name:      "Cassell Road",
					Longitude: -2.51701423067,
					Latitude:  51.4843326109,
					Bearing:   225,
				},
				models.BusStop{
					ID:        "010000002",
					Name:      "The Centre",
					Longitude: -2.59725334008,
					Latitude:  51.45306504329,
					Bearing:   0,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(tt.args.xmlFilePath)
			if err != nil {
				t.Fatal("Failed to open file", err)
			}
			defer file.Close()

			got, err := parseXML(file)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseXML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseXML() = %v, want %v", got, tt.want)
			}
		})
	}
}

var getMock func(url string) (*http.Response, error)

type httpClientMock struct{}

func (client httpClientMock) Get(url string) (*http.Response, error) {
	return getMock(url)
}

func Test_getBusStopsFromDFT(t *testing.T) {
	type args struct {
		response  *http.Response
		shouldErr bool
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Gets a zip file from department for transport",
			args: args{
				response: &http.Response{
					Status: "200 OK",
					StatusCode: 200,
					Proto: "HTTP/1.0",
					ProtoMajor: 1,
					ProtoMinor: 0,
					Header: http.Header{},
					Body: ioutil.NopCloser(strings.NewReader("hello")),
					ContentLength: 100,
					TransferEncoding: nil,
					Close: false,
					Uncompressed: false,
					Trailer: http.Header{},
					Request: &http.Request{},
					TLS: nil,
				},
				shouldErr: false,
			},
			want: []byte("hello"),
			wantErr: false,
		},
		{
			name: "Department for transport responses with a 500 error",
			args: args{
				response: &http.Response{
					Status: "500",
					StatusCode: 500,
					Proto: "HTTP/1.0",
					ProtoMajor: 1,
					ProtoMinor: 0,
					Header: http.Header{},
					Body: ioutil.NopCloser(strings.NewReader("")),
					ContentLength: 100,
					TransferEncoding: nil,
					Close: false,
					Uncompressed: false,
					Trailer: http.Header{},
					Request: &http.Request{},
					TLS: nil,
				},
				shouldErr: true,
			},
			want: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var clientMock httpClientMock
			getMock = func(url string) (*http.Response, error) {
				if tt.args.shouldErr {
					return &http.Response{}, errors.New("")
				}
				return tt.args.response, nil
			}

			got, err := getBusStopsFromDFT(clientMock)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBusStopsFromDFT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getBusStopsFromDFT() = %v, want %v", got, tt.want)
			}
		})
	}
}
