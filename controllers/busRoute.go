package controllers

import (
	"server/utils"
	"server/models"
	"encoding/json"
	"io/ioutil"
	"time"
	"net/url"
	"strings"
	"os"
	"log"
	"fmt"
	"encoding/xml"
	"strconv"
	"errors"
)

func GetRoute(lineName string, direction string, operatorID string) (models.Route, error) {
	route, err := models.GetBusRoute(lineName, direction, operatorID)
	if err != nil {
		return models.Route{}, err
	}

	return route, nil
}

func UpdateRoute(datasetID uint, httpClient httpClient, busRoute models.BusRoute) error {
	baseUrl := "https://data.bus-data.dft.gov.uk/api/v1/dataset"
	v := url.Values{}
	v.Set("api_key", os.Getenv("DFT_SECRET"))

	datasetIDString := strconv.FormatUint(uint64(datasetID), 10)

	log.Println(baseUrl + "/" + datasetIDString + "?" + v.Encode())

	resp, err := httpClient.Get(baseUrl + "/" + datasetIDString + "?" + v.Encode())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't get dft timetable", err)

		return err
	}

	if resp.StatusCode != 200 {
		fmt.Fprintln(os.Stderr, "DFT returned non 200 status of: ", resp.StatusCode)
		return errors.New("DFT returned non 200 status")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read DFT response", err)
		return errors.New("Failed to read DFT response")
	}

	timetable := timetableResults{}
	if err := json.Unmarshal(body, &timetable); err != nil {
		fmt.Fprintln(os.Stderr, "Unmarshal failed", err)
		return err
	}

	if strings.ToUpper(timetable.Extension) != "ZIP" {
		fmt.Fprintln(os.Stderr, "Folder extension was not ZIP")
		return errors.New("Folder extension was not ZIP")
	}

	return parseTimetables([]string{ timetable.URL }, httpClient, busRoute)
}

func UpdateRoutes(offset uint, limit uint, httpClient httpClient, busRoute models.BusRoute) error {
	t := time.Now()
	timeString := fmt.Sprintf("%d-%02d-%02dT00:00:00", t.Year(), t.Month(), t.Day())

	baseUrl := "https://data.bus-data.dft.gov.uk/api/v1/dataset"
	v := url.Values{}
	v.Set("api_key", os.Getenv("DFT_SECRET"))
	v.Set("limit", fmt.Sprint(offset))
	v.Set("offset", fmt.Sprint(limit))
	v.Set("status", "published")
	v.Set("startDateEnd", timeString)
	v.Set("endDateStart", timeString)

	resp, err := httpClient.Get(baseUrl + "?" + v.Encode())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't get dft timetable", err)
		return err
	}

	if resp.StatusCode != 200 {
		fmt.Fprintln(os.Stderr, "DFT returned non 200 status of: ", resp.StatusCode)
		return errors.New("DFT returned non 200 status")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read DFT response", err)
		return errors.New("Failed to read DFT response")
	}

	timetable := timetableResponse{}
	if err := json.Unmarshal(body, &timetable); err != nil {
		fmt.Fprintln(os.Stderr, "Unmarshal failed", err)
		return err
	}

	timetableUrls := make([]string, 0)
	for _, result := range timetable.Results {
		if strings.ToUpper(result.Extension) == "ZIP" {
			timetableUrls = append(timetableUrls, result.URL)
		}
	}

	return parseTimetables(timetableUrls, httpClient, busRoute)
}

func parseTimetables(timetableURIs []string, httpClient httpClient, busRoute models.BusRoute) error {
	for _, timetableURL := range timetableURIs {
		zippedFolder, err := getTimetable(timetableURL, httpClient)
		if err != nil {
			return err
		}

		zippedFiles, err := utils.UnZipFile(zippedFolder)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to unzip folder", err)

			return err
		}

		for _, zippedFile := range zippedFiles {
			rawFile, err := zippedFile.Open()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to open file", err)

				return err
			}

			file, err := ioutil.ReadAll(rawFile)
			if err != nil {
				rawFile.Close()

				fmt.Fprintln(os.Stderr, "Failed to read file", err)

				return err
			}
			rawFile.Close()

			transXChange, err := parseTransXChange(2.4, file)
			if err != nil {
				return err
			}

			if err := updateRouteTable(transXChange, busRoute); err != nil {
				fmt.Fprintln(os.Stderr, "Failed to update tables", err)
				return err
			}
		}

		time.Sleep(2 * time.Second)
	}

	return nil
}

func updateRouteTable(transXChange parsedTransXChange, busRoute models.BusRoute) error {
	if err := busRoute.InsertOperators(transXChange.operators); err != nil {
		return err
	}
	if err := busRoute.InsertLines(transXChange.lines); err != nil {
		return err
	}
	if err := busRoute.InsertJourneys(transXChange.journeys); err != nil {
		return err
	}
	if err := busRoute.InsertJourneyStops(transXChange.journeyStops); err != nil {
		return err
	}

	return nil
}

func getTimetable(url string, httpClient httpClient) ([]byte, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't get dft timetable", err)

		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Fprintln(os.Stderr, "DFT returned non 200 status of: ", resp.StatusCode)
		return nil, errors.New("DFT returned non 200 status")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read DFT response", err)
		return nil, errors.New("Failed to read DFT response")
	}

	return body, nil
}

type timetableResponse struct {
	Count    uint               `json:"count"`
	Next     string             `json:"next"`
	Previous string             `json:"previous"`
	Results  []timetableResults `json:"results"`
}

type timetableResults struct {
	ID             uint   `json:"id"`
	Created        string `json:"created"`
	Modified       string `json:"modified"`
	OperatorName   string `json:"operatorName"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Comment        string `json:"comment"`
	Status         string `json:"status"`
	URL            string `json:"url"`
	Extension      string `json:"extension"`
	FirstStartDate string `json:"firstStartDate"`
	FirstEndDate   string `json:"firstEndDate"`
	LastEndDate    string `json:"lastEndDate"`
}

type parsedTransXChange struct {
	operators    []models.Operator
	lines        []models.Line
	journeys     []models.Journey
	journeyStops []models.JourneyStop
}

func parseTransXChange(version float32, rawXML []byte) (parsedTransXChange, error) {
	if version != 2.4 {
		return parsedTransXChange{}, errors.New("Unsupported transXChange version. Supports: 2.4")
	}

	var transXChange transXChange
	err := xml.Unmarshal(rawXML, &transXChange)
	if err != nil {
		return parsedTransXChange{}, err
	}

	if version != transXChange.Version {
		return parsedTransXChange{}, errors.New("Incorrect transXChange version")
	}

	// operators
	operators := make([]models.Operator, 0)
	for _, operator := range transXChange.Operators {
		operators = append(operators, models.Operator{
			ID: operator.OperatorID,
			Name: operator.Name,
			ShortName: operator.ShortName,
		})
	}

	// services and lines
	lines := make([]models.Line, 0)
	journeys := make([]models.Journey, 0)
	for _, service := range transXChange.Services {
		// Use the LocalOperatorID to get the operatorID
		operatorID := ""
		operatorIndex := 0
		for operatorID == "" && operatorIndex < len(transXChange.Operators) {
			if transXChange.Operators[operatorIndex].LocalID == service.LocalOperatorID {
				operatorID = transXChange.Operators[operatorIndex].OperatorID
			}

			operatorIndex++
		}

		if operatorID == "" {
			fmt.Fprintln(os.Stderr, "Couldn't find a operatorID", transXChange.FileName)

			return parsedTransXChange{}, errors.New("Couldn't find a operatorID ")
		}

		lines = append(lines, models.Line{
			ID:         service.Line.ID,
			OperatorID: operatorID,
			Name:       service.Line.Name,
		})

		for _, journeyPattern := range service.JourneyPattern {
			journeys = append(journeys, models.Journey{
				LineID: service.Line.ID,
				RouteID: journeyPattern.RouteID,
				Direction: strings.ToUpper(journeyPattern.Direction),
				Description: journeyPattern.Description,
			})
		}
	}

	journeyStops := make([]models.JourneyStop, 0)
	for _, routeSection := range transXChange.RouteSections {
		// get a route ID
		var routeID string
		routeIndex := 0
		for routeID == "" && routeIndex < len(transXChange.Routes) {
			if transXChange.Routes[routeIndex].RouteSectionID == routeSection.ID {
				routeID = transXChange.Routes[routeIndex].ID
			}

			routeIndex++
		}

		if routeID == "" {
			fmt.Fprintln(os.Stderr, "Couldn't find a routeID", routeSection.ID, transXChange.FileName)

			return parsedTransXChange{}, errors.New("Couldn't find a routeID ")
		}

		// find the LineID from the Journey with the same RouteID
		var lineID string
		journeyIndex := 0
		for lineID == "" && journeyIndex < len(journeys) {
			if journeys[journeyIndex].RouteID == routeID {
				lineID = journeys[journeyIndex].LineID
			}

			journeyIndex++
		}

		if lineID == "" {
			fmt.Fprintln(os.Stderr, "Couldn't find a LineID", routeSection.ID, transXChange.FileName)

			return parsedTransXChange{}, errors.New("Couldn't find a LineID ")
		}

		previousStopDestination := ""
		for routeIndex, routeLink := range routeSection.RouteLinks {
			if routeIndex > 0 && routeLink.From != previousStopDestination {
				fmt.Fprintln(os.Stderr, "Previous stop destination does not match current stop", transXChange.FileName)

				return parsedTransXChange{}, errors.New("Previous stop destination does not match current stop")
			}

			journeyStops = append(journeyStops, models.JourneyStop{
				LineID: lineID,
				RouteID: routeID,
				StopNumber: uint(routeIndex),
				BusStopID: routeLink.From,
			})

			if routeIndex == len(routeSection.RouteLinks) - 1 {
				journeyStops = append(journeyStops, models.JourneyStop{
					LineID: lineID,
					RouteID: routeID,
					StopNumber: uint(routeIndex + 1),
					BusStopID: routeLink.To,
				})
			}
			previousStopDestination = routeLink.To
		}
	}

	return parsedTransXChange{
		operators:    operators,
		lines:        lines,
		journeys:     journeys,
		journeyStops: journeyStops,
	}, nil
}

// TransXChange
type transXChange struct {
	XMLName       xml.Name                   `xml:"TransXChange"`
	FileName      string                     `xml:"FileName"` // used for debugging
	CreatedAt     string                     `xml:"CreationDateTime,attr"`
	UpdatedAt     string                     `xml:"ModificationDateTime,attr"`
	Version       float32                    `xml:"SchemaVersion,attr"`
	Revision      uint                       `xml:"RevisionNumber,attr"`
	RouteSections []transXChangeRouteSection `xml:"RouteSections>RouteSection"`
	Routes        []transXChangeRoute        `xml:"Routes>Route"`
	Operators     []transXChangeOperator     `xml:"Operators>Operator"`
	Services      []transXChangeService      `xml:"Services>Service"`
}

type transXChangeRoute struct {
	XML            xml.Name `xml:"Route"`
	ID             string   `xml:"id,attr"`
	CreatedAt      string   `xml:"CreationDateTime,attr"`
	UpdatedAt      string   `xml:"ModificationDateTime,attr"`
	Revision       uint     `xml:"RevisionNumber,attr"`
	RouteSectionID string   `xml:"RouteSectionRef"`
	Description    string   `xml:"Description"`
}

type transXChangeRouteSection struct {
	XML        xml.Name                `xml:"RouteSection"`
	ID         string                  `xml:"id,attr"`
	RouteLinks []transXChangeRouteLink `xml:"RouteLink"`
}

type transXChangeRouteLink struct {
	XML  xml.Name `xml:"RouteLink"`
	From string   `xml:"From>StopPointRef"`
	To   string   `xml:"To>StopPointRef"`
}

type transXChangeOperator struct {
	XML        xml.Name `xml:"Operator"`
	OperatorID string   `xml:"NationalOperatorCode"`
	LocalID    string   `xml:"id,attr"` // local ID to this file
	ShortName  string   `xml:"OperatorShortName"`
	Name       string   `xml:"TradingName"`
}

type transXChangeService struct {
	XML             xml.Name                     `xml:"Service"`
	LocalOperatorID string                       `xml:"RegisteredOperatorRef"`
	Line            transXChangeServiceLine    	 `xml:"Lines>Line"`
	Origin          string                       `xml:"StandardService>Origin"`
	Destination     string                       `xml:"StandardService>Destination"`
	JourneyPattern  []transXChangeJourneyPattern `xml:"StandardService>JourneyPattern"`
}

type transXChangeServiceLine struct {
	XML      xml.Name	                       `xml:"Line"`
	ID       string                          `xml:"id,attr"`
	Name     string                          `xml:"LineName"`
	OutBound transXChangeServiceLineOutBound `xml:"OutboundDescription"`
	InBound  transXChangeServiceLineInBound  `xml:"InboundDescription"`
}

type transXChangeServiceLineOutBound struct {
	XML         xml.Name `xml:"OutboundDescription"`
	Origin      string   `xml:"Origin"`
	Destination string   `xml:"Destination"`
	Description string   `xml:"Description"`
}

type transXChangeServiceLineInBound struct {
	XML         xml.Name `xml:"InboundDescription"`
	Origin      string   `xml:"Origin"`
	Destination string   `xml:"Destination"`
	Description string   `xml:"Description"`
}

type transXChangeJourneyPattern struct {
	XML                xml.Name `xml:"JourneyPattern"`
	DestinationDisplay string   `xml:"DestinationDisplay"`
	Direction          string   `xml:"Direction"`
	Description        string   `xml:"Description"`
	RouteID            string   `xml:"RouteRef"`
}