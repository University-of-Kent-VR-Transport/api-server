package controllers

import (
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"server/models"
	"testing"
	"bytes"
	"fmt"
)

var insertOperatorsMock func(operators []models.Operator) error
var insertLinesMock func (lines []models.Line) error
var insertJourneysMock func (journeys []models.Journey) error
var insertJourneyStopsMock func (journeyStops []models.JourneyStop) error

type busRouteMock struct{}

func (busRoute busRouteMock) InsertOperators(operators []models.Operator) error {
	return insertOperatorsMock(operators)
}
func (busRoute busRouteMock) InsertLines(lines []models.Line) error {
	return insertLinesMock(lines)
}
func (busRoute busRouteMock) InsertJourneys(journeys []models.Journey) error {
	return insertJourneysMock(journeys)
}
func (busRoute busRouteMock) InsertJourneyStops(journeyStops []models.JourneyStop) error {
	return insertJourneyStopsMock(journeyStops)
}

func TestUpdateRoutes(t *testing.T) {
	type httpResponse struct {
		StatusCode int
		BodyDir    string
	}
	type args struct {
		offset                uint
		limit                 uint
		getResponse           []httpResponse
		getError              bool
		insertOperatorsErr    bool
		insertLinesErr        bool
		insertJourneysErr     bool
		insertJourneyStopsErr bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Updates routes",
			args: args{
				offset: 0,
				limit: 1,
				getResponse: []httpResponse{
					httpResponse{
						StatusCode: 200,
						BodyDir: "testdata/dft-timetable-query.json",
					},
					httpResponse{
						StatusCode: 200,
						BodyDir: "testdata/dft-timetable.zip",
					},
				},
				getError: false,
				insertOperatorsErr: false,
				insertLinesErr: false,
				insertJourneysErr: false,
				insertJourneyStopsErr: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var httpClient httpClientMock
			requestCount := 0
			getMock = func(url string) (*http.Response, error) {
				if tt.args.getError {
					return &http.Response{}, errors.New("")
				}

				body, err := ioutil.ReadFile(tt.args.getResponse[requestCount].BodyDir)
				if err != nil {
					t.Fatal(err)
				}

				response := &http.Response{
					Status: "200 OK",
					StatusCode: tt.args.getResponse[requestCount].StatusCode,
					Proto: "HTTP/1.0",
					ProtoMajor: 1,
					ProtoMinor: 0,
					Header: http.Header{},
					Body: ioutil.NopCloser(bytes.NewReader(body)),
					ContentLength: 100,
					TransferEncoding: nil,
					Close: false,
					Uncompressed: false,
					Trailer: http.Header{},
					Request: &http.Request{},
					TLS: nil,
				}

				fmt.Println(requestCount)

				requestCount++
				return response, nil
			}

			var busRoute busRouteMock
			insertOperatorsMock = func(operators []models.Operator) error {
				if tt.args.insertOperatorsErr {
					return errors.New("")
				}
				return nil
			}
			insertLinesMock = func(lines []models.Line) error {
				if tt.args.insertLinesErr {
					return errors.New("")
				}
				return nil
			}
			insertJourneysMock = func(journey []models.Journey) error {
				if tt.args.insertJourneysErr {
					return errors.New("")
				}
				return nil
			}
			insertJourneyStopsMock = func(journeyStop []models.JourneyStop) error {
				if tt.args.insertJourneyStopsErr {
					return errors.New("")
				}
				return nil
			}

			if err := UpdateRoutes(tt.args.offset, tt.args.limit, httpClient, busRoute); (err != nil) != tt.wantErr {
				t.Errorf("UpdateRoutes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_parseTransXChange(t *testing.T) {
	type args struct {
		version float32
		xmlFile string
	}
	tests := []struct {
		name    string
		args    args
		want    parsedTransXChange
		wantErr bool
	}{
		{
			name: "Gets two routes from SCEK-953",
			args: args{
				version: 2.4,
				xmlFile: "./testdata/dft-timetable.xml",
			},
			want: parsedTransXChange{
				operators: []models.Operator{
					models.Operator{
						ID:        "SCEK",
						Name:      "Stagecoach in East Kent",
						ShortName: "Stagecoach",
					},
				},
				lines: []models.Line{
					models.Line{
						ID:         "SCEK:PK0000098:84_953_953:953:",
						OperatorID: "SCEK",
						Name:       "953",
					},
				},
				journeys: []models.Journey{
					models.Journey{
						LineID:      "SCEK:PK0000098:84_953_953:953:",
						RouteID:     "RT131",
						Direction:   "OUTBOUND",
						Description: "Bus Station - High School",
					},
					models.Journey{
						LineID:      "SCEK:PK0000098:84_953_953:953:",
						RouteID:     "RT132",
						Direction:   "INBOUND",
						Description: "Bus Station - High School",
					},
				},
				journeyStops: []models.JourneyStop{
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:84_953_953:953:",
						RouteID:    "RT131",
						StopNumber: 0,
						BusStopID:  "240098892",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:84_953_953:953:",
						RouteID:    "RT131",
						StopNumber: 1,
						BusStopID:  "2400A049530A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:84_953_953:953:",
						RouteID:    "RT131",
						StopNumber: 2,
						BusStopID:  "240096713",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:84_953_953:953:",
						RouteID:    "RT131",
						StopNumber: 3,
						BusStopID:  "2400A050490A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:84_953_953:953:",
						RouteID:    "RT131",
						StopNumber: 4,
						BusStopID:  "2400A050530A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:84_953_953:953:",
						RouteID:    "RT132",
						StopNumber: 0,
						BusStopID:  "2400A050510A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:84_953_953:953:",
						RouteID:    "RT132",
						StopNumber: 1,
						BusStopID:  "2400A050500A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:84_953_953:953:",
						RouteID:    "RT132",
						StopNumber: 2,
						BusStopID:  "240096711",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:84_953_953:953:",
						RouteID:    "RT132",
						StopNumber: 3,
						BusStopID:  "2400A049550A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:84_953_953:953:",
						RouteID:    "RT132",
						StopNumber: 4,
						BusStopID:  "2400A039640A",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Gets two routes from uni1",
			args: args{
				version: 2.4,
				xmlFile: "./testdata/dft-timetable-uni1.xml",
			},
			want: parsedTransXChange{
				operators: []models.Operator{
					models.Operator{
						ID:        "SCEK",
						Name:      "Stagecoach in East Kent",
						ShortName: "Stagecoach",
					},
				},
				lines: []models.Line{
					models.Line{
						ID:         "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						OperatorID: "SCEK",
						Name:       "Uni1",
					},
				},
				journeys: []models.Journey{
					models.Journey{
						LineID:      "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:     "RT197",
						Direction:   "OUTBOUND",
						Description: "City Centre - University",
					},
					models.Journey{
						LineID:      "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:     "RT199",
						Direction:   "INBOUND",
						Description: "City Centre - University",
					},
					models.Journey{
						LineID:      "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:     "RT198",
						Direction:   "INBOUND",
						Description: "City Centre - University",
					},
					models.Journey{
						LineID:      "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:     "RT199",
						Direction:   "INBOUND",
						Description: "City Centre - University",
					},
					models.Journey{
						LineID:      "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:     "RT199",
						Direction:   "INBOUND",
						Description: "City Centre - University",
					},
					models.Journey{
						LineID:      "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:     "RT199",
						Direction:   "INBOUND",
						Description: "City Centre - University",
					},
				},
				journeyStops: []models.JourneyStop{
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT197",
						StopNumber: 0,
						BusStopID:  "240098906",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT197",
						StopNumber: 1,
						BusStopID:  "2400A049530A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT197",
						StopNumber: 2,
						BusStopID:  "2400A048110A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT197",
						StopNumber: 3,
						BusStopID:  "2400A048140A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT197",
						StopNumber: 4,
						BusStopID:  "2400A048160A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT197",
						StopNumber: 5,
						BusStopID:  "2400A048170A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT197",
						StopNumber: 6,
						BusStopID:  "2400A050260A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT197",
						StopNumber: 7,
						BusStopID:  "2400A050270A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT197",
						StopNumber: 8,
						BusStopID:  "2400A050290A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT197",
						StopNumber: 9,
						BusStopID:  "2400100704",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT197",
						StopNumber: 10,
						BusStopID:  "2400A040360A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT197",
						StopNumber: 11,
						BusStopID:  "240095612",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT198",
						StopNumber: 0,
						BusStopID:  "240095612",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT198",
						StopNumber: 1,
						BusStopID:  "2400A040370A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT198",
						StopNumber: 2,
						BusStopID:  "240075428",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT198",
						StopNumber: 3,
						BusStopID:  "240097602",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT198",
						StopNumber: 4,
						BusStopID:  "240097597",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT198",
						StopNumber: 5,
						BusStopID:  "2400100702",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT198",
						StopNumber: 6,
						BusStopID:  "2400A050300A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT198",
						StopNumber: 7,
						BusStopID:  "2400A050280A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT198",
						StopNumber: 8,
						BusStopID:  "2400A048180A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT198",
						StopNumber: 9,
						BusStopID:  "2400A048150A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT198",
						StopNumber: 10,
						BusStopID:  "2400105752",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT198",
						StopNumber: 11,
						BusStopID:  "2400A040380A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT198",
						StopNumber: 12,
						BusStopID:  "240096711",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT198",
						StopNumber: 13,
						BusStopID:  "2400A049550A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT198",
						StopNumber: 14,
						BusStopID:  "2400A039640A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT199",
						StopNumber: 0,
						BusStopID:  "240095612",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT199",
						StopNumber: 1,
						BusStopID:  "2400A040370A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT199",
						StopNumber: 2,
						BusStopID:  "2400100702",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT199",
						StopNumber: 3,
						BusStopID:  "2400A050300A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT199",
						StopNumber: 4,
						BusStopID:  "2400A050280A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT199",
						StopNumber: 5,
						BusStopID:  "2400A048180A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT199",
						StopNumber: 6,
						BusStopID:  "2400A048150A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT199",
						StopNumber: 7,
						BusStopID:  "2400105752",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT199",
						StopNumber: 8,
						BusStopID:  "2400A040380A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT199",
						StopNumber: 9,
						BusStopID:  "240096711",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT199",
						StopNumber: 10,
						BusStopID:  "2400A049550A",
					},
					models.JourneyStop{
						LineID:     "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
						RouteID:    "RT199",
						StopNumber: 11,
						BusStopID:  "2400A039640A",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xml, err := ioutil.ReadFile(tt.args.xmlFile)
			if err != nil {
				t.Fatal(err)
			}

			got, err := parseTransXChange(tt.args.version, xml)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTransXChange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got.operators, tt.want.operators) {
				t.Errorf("parseTransXChange() = Operators{ %v }, want Operators{ %v }", got.operators, tt.want.operators)
			}

			if !reflect.DeepEqual(got.lines, tt.want.lines) {
				t.Errorf("parseTransXChange() = Lines{ %v }, want Lines{ %v }", got.lines, tt.want.lines)
			}

			if !reflect.DeepEqual(got.journeys, tt.want.journeys) {
				t.Errorf("parseTransXChange() = Journeys{ %v }, want Journeys{ %v }", got.journeys, tt.want.journeys)
			}

			if !reflect.DeepEqual(got.journeyStops, tt.want.journeyStops) {
				t.Errorf("parseTransXChange() = JourneyStops{ %v }, want JourneyStops{ %v }", got.journeyStops, tt.want.journeyStops)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTransXChange() = %v, want %v", got, tt.want)
			}
		})
	}
}