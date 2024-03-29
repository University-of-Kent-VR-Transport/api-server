package controllers

import (
	"server/utils"
	"server/models"
	"io"
	"log"
	"os"
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/xml"
	"database/sql"
	_ "github.com/lib/pq"
)

const naptanURL = "https://naptan.app.dft.gov.uk/Datarequest/naptan.ashx"

// UpdateBusStops updates all bus stops using the NaPTAN database and returns
// a background job
func UpdateBusStops() (models.BackgroundJob, error) {
	// Set job in db
	connectionString := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Println("Couldn't connect to db", err)

		return models.BackgroundJob{}, err
	}
	defer db.Close()

	job, err := models.CreateBackgroundJob("UPDATE NATIONAL PUBLIC TRANSPORT ACCESS NODES", db)
	if err != nil {
		return models.BackgroundJob{}, err
	}

	go runUpdate(job.ID)

	return job, nil
}

func runUpdate(jobID uint) {
	connectionString, _ := os.LookupEnv("DATABASE_URL")
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Println("Couldn't connect to db", err)
		models.UpdateBackgroundJob(jobID, "FAILED", db)
		return
	}
	defer db.Close()

	log.Println("Connected to DB")

	// Get NaPTAN from naptanURL
	zippedFolder, err := getBusStopsFromDFT(&http.Client{})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get NaPTAN from DFT")
		models.UpdateBackgroundJob(jobID, "FAILED", db)
		return
	}

	// UnZip folder
	rawFile, err := utils.UnZipFile(zippedFolder)
	if err != nil {
		log.Println("Failed to unzip folder")
		models.UpdateBackgroundJob(jobID, "FAILED", db)
		return
	}

	file, err := rawFile[0].Open()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to unzip folder")
		models.UpdateBackgroundJob(jobID, "FAILED", db)
		return
	}
	defer file.Close()

	// Parse xml
	busStops, err := parseXML(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse xml")
		models.UpdateBackgroundJob(jobID, "FAILED", db)
		return
	}

	// Insert using model
	if err := models.UpdateBusStops(busStops, db); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to rebuild bus stops")
		models.UpdateBackgroundJob(jobID, "FAILED", db)
		return
	}

	// Complete background job
	models.UpdateBackgroundJob(jobID, "COMPLETE", db)
}

type httpClient interface {
	Get(url string) (*http.Response, error)
}

func getBusStopsFromDFT(client httpClient) ([]byte, error) {
	resp, err := client.Get(naptanURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to make request to DFT", err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "DFT returned non 200 status of: %d/n", resp.StatusCode)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read DFT response", err)
		return nil, err
	}

	return body, nil
}

type stopPoint struct {
	XML        xml.Name `xml:"StopPoint"`
	ID         string   `xml:"AtcoCode"`
	Status     string   `xml:"Status,attr"`
	Name       string   `xml:"Descriptor>CommonName"`
	Longitude  float32  `xml:"Place>Location>Translation>Longitude"`
	Latitude   float32  `xml:"Place>Location>Translation>Latitude"`
	Bearing    float32  `xml:"StopClassification>OnStreet>Bus>MarkedPoint>Bearing>Degrees"`
}

func parseXML(xmlFile io.ReadCloser) ([]models.BusStop, error) {
	decoder := xml.NewDecoder(xmlFile)
	stopPoints := make([]stopPoint, 0)

	reachedEndOfFile := false

	for !reachedEndOfFile {
		token, err := decoder.Token()
		if err != nil {
			if err != io.EOF {
				fmt.Fprintln(os.Stderr, "Failed to decode xml token", err)
				return nil, err
			}
			reachedEndOfFile = true
		}

		switch currentElement := token.(type) {
			case xml.StartElement:
				if currentElement.Name.Local == "StopPoint" {
					var stopPoint stopPoint

					if err := decoder.DecodeElement(&stopPoint, &currentElement); err != nil {
						fmt.Fprintln(os.Stderr, "Failed to decode StopPoint", err)
						return nil, err
					}

					if stopPoint.Status == "active" {
						stopPoints = append(stopPoints, stopPoint)
					}
				}
			default:
		}
	}

	log.Println("Stop Points: ", len(stopPoints))

	busStops := make([]models.BusStop, 0)

	for _, stopPoint := range stopPoints {
		busStop := models.BusStop{
			ID: strings.TrimSpace(stopPoint.ID),
			Name: strings.TrimSpace(stopPoint.Name),
			Longitude: stopPoint.Longitude,
			Latitude: stopPoint.Latitude,
			Bearing: stopPoint.Bearing,
		}

		busStops = append(busStops, busStop)
	}

	return busStops, nil
}