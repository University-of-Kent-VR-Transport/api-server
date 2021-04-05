package handlers

import (
	"encoding/json"
	"server/models"
	"io"
	"log"
	"archive/zip"
	"bytes"
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

type getUpdateBusStopsResponse struct {
	Job  models.BackgroundJob
}

func UpdateBusStops(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	// Set job in db
	connectionString, _ := os.LookupEnv("DATABASE_URL")

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Println("Couldn't connect to db", err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}
	defer db.Close()

	job, err := models.CreateBackgroundJob("UPDATE NATIONAL PUBLIC TRANSPORT ACCESS NODES", db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, http.StatusText(http.StatusInternalServerError))

		return
	}

	// Response with 202
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	response := getUpdateBusStopsResponse{Job: job}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error while encoding update naptan response: %s", err.Error())
	}

	go runUpdate(job.ID)
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
		log.Println("Failed to get NaPTAN from DFT")
		models.UpdateBackgroundJob(jobID, "FAILED", db)
		return
	}

	log.Println("Fetched NaPTAN from DFT")

	// UnZip folder
	rawFile, err := unZipFile(zippedFolder)
	if err != nil {
		log.Println("Failed to unzip folder")
		models.UpdateBackgroundJob(jobID, "FAILED", db)
		return
	}
	defer rawFile.Close()

	log.Println("Unzipped folder")

	// Parse xml
	busStops, err := parseXML(rawFile)
	if err != nil {
		log.Println("Failed to parse xml")
		models.UpdateBackgroundJob(jobID, "FAILED", db)
		return
	}

	log.Println("Parsed XML")

	fmt.Printf("Found %v stop points\n", len(busStops))

	// Insert using model
	if err := models.RebuildBusStops(busStops, db); err != nil {
		log.Println("Failed to rebuild bus stops")
		models.UpdateBackgroundJob(jobID, "FAILED", db)
		return
	}

	log.Println("Updated table")

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

func unZipFile(zippedFolder []byte) (io.ReadCloser, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(zippedFolder), int64(len(zippedFolder)))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create zip reader", err)
		return nil, err
	}

	fmt.Println(len(zipReader.File))

	fileToUnzip := zipReader.File[0]
	fmt.Println("Reading file:", fileToUnzip.Name)

	file, err := fileToUnzip.Open()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open zip", err)
		return nil, err
	}

	return file, nil
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