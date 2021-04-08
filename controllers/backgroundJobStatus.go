package controllers

import (
	"server/models"
	"fmt"
	"os"
	"database/sql"
	_ "github.com/lib/pq"
)

// GetBackgroundJobStatus gets the job with the id of jobID.
func GetBackgroundJobStatus(jobID uint) (models.BackgroundJob, bool, error) {
	connectionString := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't connect to db", err)

		return models.BackgroundJob{}, false, err
	}
	defer db.Close()

	job, err := models.GetBackgroundJob(jobID, db)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.BackgroundJob{}, false, nil
		}
		
		return models.BackgroundJob{}, false, err
	}

	return job, true, nil
}