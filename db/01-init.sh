#!/bin/bash
set -e
export PGPASSWORD=$POSTGRES_PASSWORD;
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE USER $APP_DB_USER WITH PASSWORD '$APP_DB_PASS';

  CREATE DATABASE $APP_DB_NAME;

	GRANT CREATE ON DATABASE $APP_DB_NAME TO $APP_DB_USER;

  \connect $APP_DB_NAME $APP_DB_USER
  BEGIN;
    CREATE TABLE IF NOT EXISTS bus_stop (
			id CHAR(12) NOT NULL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			longitude DOUBLE PRECISION NOT NULL,
			latitude DOUBLE PRECISION NOT NULL,
			bearing DOUBLE PRECISION NOT NULL
		);

		CREATE TYPE job_type AS ENUM ('UPDATE NATIONAL PUBLIC TRANSPORT ACCESS NODES');
		CREATE TYPE status AS ENUM ('RUNNING', 'COMPLETE', 'FAILED');
		CREATE TABLE IF NOT EXISTS background_job (
			id SERIAL NOT NULL PRIMARY KEY,
			type job_type NOT NULL,
			status status NOT NULL DEFAULT 'RUNNING',
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
  COMMIT;

	GRANT SELECT ON TABLE bus_stop TO $APP_DB_USER;
	GRANT INSERT ON TABLE bus_stop TO $APP_DB_USER;
	GRANT TRUNCATE ON TABLE bus_stop TO $APP_DB_USER;

	GRANT SELECT ON TABLE background_job TO $APP_DB_USER;
	GRANT INSERT ON TABLE background_job TO $APP_DB_USER;
	GRANT UPDATE ON TABLE background_job TO $APP_DB_USER;
EOSQL