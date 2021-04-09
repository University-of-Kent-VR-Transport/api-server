package models

import (
	"os"
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
)

// Operator
// | ID        | ShortName  | Name                    |
// | --------- | ---------- | ----------------------- |
// | PK String | String     | String                  |
// | SCEK      | Stagecoach | Stagecoach in East Kent |

// Line
// | ID        | Name   | OperatorID    |
// | --------- | ------ | ------------- |
// | PK String | String | FK OperatorID |
// | SCEK:PK...| 953    | SCEK          |

// Journey
// | LineID       | RouteID   | Direction        | Description |
// | ------------ | --------- | ---------------- | ----------- |
// | PK FK LineID | PK String | INBOUND/OUTBOUND | String      |
// | SCEK:PK...   | RT132     | INBOUND          | Bus Stat... |

// Journey Stop
// | LineID              | RouteID               | StopNumber | BusStopID    |
// | ------------------- | --------------------- | ---------- | ------------ |
// | PK FK JourneyLineID | PK FK JourneyRouteID  | PK Uint    | FK BusStopID |
// | SCEK:PK...          | RT132                 | 0          | 240098892    |

type Route struct {
	LineID       string
	RouteID      string
	OperatorID   string
	OperatorName string
	Name         string
	Direction    string
	Description  string
	Origin       string
	Destination  string
	Stops        []BusStop
}

type BusRoutes struct {}
type BusRoute interface {
	InsertOperators(operators []Operator) error
	InsertLines(lines []Line) error
	InsertJourneys(journeys []Journey) error
	InsertJourneyStops(journeyStops []JourneyStop) error
}

const getRouteByLineDirectionOperator = `SELECT
	journey_stop.route_id AS routeID,
	bus_stop.id AS stopID,
	bus_stop.name AS stopName,
	bus_stop.longitude AS longitude,
	bus_stop.latitude AS latitude,
	bus_stop.bearing AS bearing,
	journey.direction AS direction,
	journey.description AS description,
	line.name AS lineName,
	line.id AS lineID,
	operator.name AS operatorName,
	operator.short_name AS operatorShortName
FROM
	journey_stop
INNER JOIN bus_stop ON journey_stop.bus_stop_id = bus_stop.id
INNER JOIN journey ON journey_stop.line_id = journey.line_id AND journey_stop.route_id = journey.route_id
INNER JOIN line ON journey_stop.line_id = line.id
INNER JOIN operator ON line.operator_id = operator.id
WHERE line.name=$1 AND journey.direction=$2 AND operator.id=$3
ORDER BY journey_stop.route_id, journey_stop.stop_number`

func GetBusRoute(lineName string, direction string, operatorID string) (Route, error) {
	connectionString := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to db", err)

		return Route{}, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getRouteByLineDirectionOperator)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to prepare select route statement", err)

		return Route{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(lineName, direction, operatorID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to execute select route statement", err)

		return Route{}, err
	}
	defer rows.Close()

	route := Route{}
	busStops := make([]BusStop, 0)

	for rows.Next() {
		var routeId, stopID, stopName, direction, description, lineName, lineID, operatorName, operatorShortName string
		var longitude, latitude, bearing float32

    err = rows.Scan(&routeId, &stopID, &stopName, &longitude, &latitude, &bearing, &direction, &description, &lineName, &lineID, &operatorName, &operatorShortName)
    if err != nil {
			fmt.Fprintln(os.Stderr, err)

			return Route{}, err
		}

		route.LineID = lineID
		route.RouteID = routeId
		route.OperatorID = operatorID

		if operatorName != "" {
			route.OperatorName = operatorName
		} else {
			route.OperatorName = operatorShortName
		}

		route.Name = lineName
		route.Direction = direction
		route.Description = description

		if route.Origin == "" {
			route.Origin = stopName
		}

		route.Destination = stopName

		busStops = append(busStops, BusStop{
			ID: stopID,
			Name: stopName,
			Longitude: longitude,
			Latitude: latitude,
			Bearing: bearing,
		})
	}

	route.Stops = busStops

	return route, nil
}

const insertOperatorSQL string = "INSERT INTO operator(id, name, short_name) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING"

func (BusRoutes *BusRoutes) InsertOperators(operators []Operator) error {
	connectionString := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to db", err)

		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertOperatorSQL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to prepare insert operator statement", err)

		return err
	}
	defer stmt.Close()

	for _, operator := range operators {
		_, err := stmt.Exec(operator.ID, operator.Name, operator.ShortName)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to execute insert operator statement", err)

			return err
		}
	}

	return nil
}

const insertLineSQL string = "INSERT INTO line(id, name, operator_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING"

func (BusRoutes *BusRoutes) InsertLines(lines []Line) error {
	connectionString := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to db", err)

		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertLineSQL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to prepare insert line statement", err)

		return err
	}
	defer stmt.Close()

	for _, line := range lines {
		_, err := stmt.Exec(line.ID, line.Name, line.OperatorID)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to execute insert line statement", err)

			return err
		}
	}

	return nil
}

const insertJourneysSQL string = "INSERT INTO journey(line_id, route_id, direction, description) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING"

func (BusRoutesBR *BusRoutes) InsertJourneys(journeys []Journey) error {
	connectionString := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to db", err)

		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertJourneysSQL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to prepare insert journey statement", err)

		return err
	}
	defer stmt.Close()

	for _, journey := range journeys {
		_, err := stmt.Exec(journey.LineID, journey.RouteID, journey.Direction, journey.Description)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to execute insert journey statement", err)

			return err
		}
	}

	return nil
}

const insertJourneyStopsSQL string = "INSERT INTO journey_stop(line_id, route_id, stop_number, bus_stop_id) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING"

func (BusRoutes *BusRoutes) InsertJourneyStops(journeyStops []JourneyStop) error {
	connectionString := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to db", err)

		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertJourneyStopsSQL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to prepare insert journey stop statement", err)

		return err
	}
	defer stmt.Close()

	for _, journeyStop := range journeyStops {
		_, err := stmt.Exec(journeyStop.LineID, journeyStop.RouteID, journeyStop.StopNumber, journeyStop.BusStopID)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to execute insert journey stop statement", err)

			// Consume error. Best to continue than exit
		}
	}

	return nil
}

type Operator struct {
	ID        string
	Name      string
	ShortName string
}

// Line will have an outbound and inbound service
type Line struct {
	ID         string
	OperatorID string
	Name       string
}

// LineID, RouteID used as primary keys
type Journey struct {
	LineID      string
	RouteID     string
	Direction   string
	Description string
}

// LineID, RouteID and StopNumber used as primary keys
type JourneyStop struct {
	LineID     string
	RouteID    string
	StopNumber uint
	BusStopID  string
}