package models

import (
	"context"
	"time"
	"fmt"
	"errors"
	"log"
	"database/sql"
	_ "github.com/lib/pq"
)

type BackgroundJob struct {
	ID        uint
	URI       string
	Type      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

const selectRunningJob string = "SELECT id FROM background_job WHERE type = $1 AND status = 'RUNNING'"
const insertNewJob string = "INSERT INTO background_job(type) VALUES($1) RETURNING id, created_at"
const selectJob string = "SELECT type, status, created_at, updated_at FROM background_job WHERE id = $1"
const updateJob string = "UPDATE background_job SET status = $1 WHERE id = $2"

func CreateBackgroundJob(jobType string, db *sql.DB) (BackgroundJob, error) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Couldn't create database transaction")
		return BackgroundJob{}, err
	}

	var runningJobID uint
	jobTypeRunning := true

	if err := tx.QueryRow(selectRunningJob, jobType).Scan(&runningJobID); err != nil {
		if err != sql.ErrNoRows {
			log.Println("Couldn't select running jobs", err)
			tx.Rollback()
			return BackgroundJob{}, errors.New("Couldn't select running jobs")
		}

		jobTypeRunning = false
	}

	if jobTypeRunning {
		log.Println("Job type running")
		tx.Rollback()
		return BackgroundJob{}, errors.New("Job type running")
	}

	var jobID uint
	var createdAt time.Time

	if err := tx.QueryRow(insertNewJob, jobType).Scan(&jobID, &createdAt); err != nil {
		log.Println("Error inserting background job to db", err)
		tx.Rollback()
		return BackgroundJob{}, errors.New("Error inserting into background_job")
	}

	if err := tx.Commit(); err != nil {
		log.Println("Transaction failed", err)
		tx.Rollback()
		return BackgroundJob{}, errors.New("Transaction failed")
	}

	return BackgroundJob{
		ID: jobID,
		URI: fmt.Sprintf("/api/jobs/%v", jobID),
		Type: jobType,
		Status: "RUNNING",
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}, nil
}

func UpdateBackgroundJob(id uint, status string, db *sql.DB) error {
	if _, err := db.Query(updateJob, status, id); err != nil {
		log.Println("Error updating background job in db", err)
		return errors.New("Error updating background_job")
	}

	return nil
}

type sqlDB interface {
	QueryRow(query string, args ...interface{}) *sql.Row
}

func GetBackgroundJob(jobID uint, db sqlDB) (BackgroundJob, error) {
	var jobType string
	var status string
	var createdAt time.Time
	var updatedAt time.Time

	err := db.QueryRow(selectJob, jobID).Scan(&jobType, &status, &createdAt, &updatedAt)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Println("Error getting background job from db", err)
		}

		return BackgroundJob{}, err
	}

	job := BackgroundJob{
		ID: jobID,
		URI: fmt.Sprintf("/api/jobs/%v", jobID),
		Type: jobType,
		Status: status,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return job, nil
}