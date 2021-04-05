package models

import (
	"context"
	"server/types"
	"os"
	"log"
	"fmt"
	"database/sql"
	"github.com/lib/pq"
)

// | ID        | Name         | Longitude      | Latitude       | Bearing |
// | ----------| ------------ | -------------- | -------------- | ------- |
// | PK String | String       | Float          | Float          | Float   |
// | 010000001 | Cassell Road | -2.51701423067 | 51.4843326109  | 225     |
// | 010000002 | The Centre   | -2.59725334008 | 51.45306504329 | 0       |

type BusStop struct {
	ID        string
	Name      string
	Longitude float32
	Latitude  float32
	Bearing   float32
}

type NaptanLocation struct {
	name string
	coordinate types.Coordinate
	bearing float32
}

const selectNaptanByID = "SELECT name, longitude, latitude, bearing FROM naptan WHERE atcoCode = $1"
const selectStopsWithinBounds = "SELECT id, name, longitude, latitude, bearing FROM bus_stop WHERE longitude >= $1 AND latitude >= $2 AND longitude <= $3 AND latitude <= $4 LIMIT 200"

func GetBusStopWithinBounds(minLongitude float32, minLatitude float32, maxLongitude float32, maxLatitude float32) ([]BusStop, error) {
	connectionString, _ := os.LookupEnv("DATABASE_URL")

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(selectStopsWithinBounds, minLongitude, minLatitude, maxLongitude, maxLatitude)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	defer rows.Close()

	busStops := make([]BusStop, 0)

	for rows.Next() {
		var id, name string
		var longitude, latitude, bearing float32

    err = rows.Scan(&id, &name, &longitude, &latitude, &bearing)
    if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil, err
		}

		busStops = append(busStops, BusStop{
			ID: id,
			Name: name,
			Longitude: longitude,
			Latitude: latitude,
			Bearing: bearing,
		})
	}

	return busStops, nil
}

func GetLocationFromNaPTAN(id string) (NaptanLocation, error) {
	connectionString, _ := os.LookupEnv("DATABASE_URL")

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return NaptanLocation{}, err
	}
	defer db.Close()

	var (
		name string
		longitude float64
		latitude float64
		bearing float64
	)
	if err := db.QueryRow(selectNaptanByID, id).Scan(&name, &longitude, &latitude, &bearing); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return NaptanLocation{}, err
	}

	return NaptanLocation{
		name: name,
		coordinate: types.Coordinate{
			Longitude: float32(longitude),
			Latitude: float32(latitude),
		},
		bearing: float32(bearing),
	}, nil
}

// Insert shit tones
func RebuildBusStops(busStops []BusStop, db *sql.DB) error {
	ctx := context.Background()
	txn, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Couldn't create database transaction", err)
		return err
	}

	log.Println("Created transaction")

	if _, err := txn.Query("TRUNCATE TABLE bus_stop"); err != nil {
		log.Fatal(err)
		txn.Rollback()
		return err
	}

	log.Println("Truncated bus_stop table")

	stmt, err := txn.Prepare(pq.CopyIn("bus_stop", "id", "name", "longitude", "latitude", "bearing"))
	if err != nil {
		log.Fatal(err)
		txn.Rollback()
		return err
	}

	log.Println("Build prepared sql statement")

	for _, busStop := range busStops {
		_, err := stmt.Exec(busStop.ID, busStop.Name, busStop.Longitude, busStop.Latitude, busStop.Bearing)
		if err != nil {
			log.Fatal(err)
			stmt.Close()
			txn.Rollback()
			return err
		}
	}

	log.Println("Filled prepared statement with busStops")

	if _, err := stmt.Exec(); err != nil {
		log.Fatal(err)
		stmt.Close()
		txn.Rollback()
		return err
	}

	log.Println("Executed prepared statements")

	if err := stmt.Close(); err != nil {
		log.Fatal(err)
		txn.Rollback()
		return err
	}

	log.Println("Closed prepared statements")

	if err := txn.Commit(); err != nil {
		log.Fatal(err)
		txn.Rollback()
		return err
	}

	log.Println("Commited transaction")

	return nil
}