package models

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"server/types"
	"strconv"
	"strings"
)

// GetBusLocation returns the live bus location within two coordinates
func GetBusLocation(
	coordinateTopLeft types.Coordinate,
	coordinateBottomRight types.Coordinate,
) (types.Siri, error) {
	requestURL, err := buildRequestURL(coordinateTopLeft, coordinateBottomRight)
	if err != nil {
		return types.Siri{}, err
	}

	resp, err := http.Get(requestURL)
	if err != nil {
		log.Println("Request to DFT failed", err)
		return types.Siri{}, errors.New("Request to DFT failed")
	}

	if resp.StatusCode != 200 {
		log.Printf("DFT returned non 200 status of: %d", resp.StatusCode)
		return types.Siri{}, errors.New("DFT returned non 200 status")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read DFT response", err)
		return types.Siri{}, errors.New("Failed to read DFT response")
	}

	siriData, err := parseRawResponse(body)
	if err != nil {
		return types.Siri{}, err
	}

	return siriData, nil
}

func buildRequestURL(
	coordinateTopLeft types.Coordinate,
	coordinateBottomRight types.Coordinate,
) (string, error) {
	urlBuilder := strings.Builder{}

	if _, err := urlBuilder.WriteString(
		"https://data.bus-data.dft.gov.uk/api/v1/datafeed?status=published&api_key=",
	); err != nil {
		log.Println("Failed to build request URL", err)
		return "", errors.New("Failed to build request URL")
	}

	if _, err := urlBuilder.WriteString(os.Getenv("DFT_SECRET")); err != nil {
		log.Println("Failed to build request URL", err)
		return "", errors.New("Failed to build request URL")
	}

	if _, err := urlBuilder.WriteString("&boundingBox="); err != nil {
		log.Println("Failed to build request URL", err)
		return "", errors.New("Failed to build request URL")
	}

	if _, err := urlBuilder.WriteString(
		strconv.FormatFloat(float64(coordinateTopLeft.Longitude), 'g', -1, 32),
	); err != nil {
		log.Println("Failed to build request URL", err)
		return "", errors.New("Failed to build request URL")
	}

	if _, err := urlBuilder.WriteString(","); err != nil {
		log.Println("Failed to build request URL", err)
		return "", errors.New("Failed to build request URL")
	}

	if _, err := urlBuilder.WriteString(
		strconv.FormatFloat(float64(coordinateBottomRight.Latitude), 'g', -1, 32),
	); err != nil {
		log.Println("Failed to build request URL", err)
		return "", errors.New("Failed to build request URL")
	}

	if _, err := urlBuilder.WriteString(","); err != nil {
		log.Println("Failed to build request URL", err)
		return "", errors.New("Failed to build request URL")
	}

	if _, err := urlBuilder.WriteString(
		strconv.FormatFloat(float64(coordinateBottomRight.Longitude), 'g', -1, 32),
	); err != nil {
		log.Println("Failed to build request URL", err)
		return "", errors.New("Failed to build request URL")
	}

	if _, err := urlBuilder.WriteString(","); err != nil {
		log.Println("Failed to build request URL", err)
		return "", errors.New("Failed to build request URL")
	}

	if _, err := urlBuilder.WriteString(
		strconv.FormatFloat(float64(coordinateTopLeft.Latitude), 'g', -1, 32),
	); err != nil {
		log.Println("Failed to build request URL", err)
		return "", errors.New("Failed to build request URL")
	}

	return urlBuilder.String(), nil
}

func parseRawResponse(rawResponse []byte) (types.Siri, error) {
	var siriData types.Siri
	err := xml.Unmarshal(rawResponse, &siriData)

	if err != nil {
		log.Println("Failed to parse bus location XML from DFT", err)
		return types.Siri{}, errors.New("Failed to parse bus location XML from DFT")
	}

	return siriData, nil
}
