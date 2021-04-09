# API-Server

[![CI/CD](https://github.com/University-of-Kent-VR-Transport/api-server/actions/workflows/continuous-integration-and-delivery.yml/badge.svg?branch=master)](https://github.com/University-of-Kent-VR-Transport/api-server/actions/workflows/continuous-integration-and-delivery.yml)
[![codecov](https://codecov.io/gh/University-of-Kent-VR-Transport/api-server/branch/master/graph/badge.svg?token=PL3PK3N5RC)](https://codecov.io/gh/University-of-Kent-VR-Transport/api-server)
[![Icicle Coverage](https://codecov.io/gh/University-of-Kent-VR-Transport/api-server/branch/master/graphs/icicle.svg?token=PL3PK3N5RC)](https://app.codecov.io/gh/University-of-Kent-VR-Transport/api-server/branch/master)

## Contents

- [About](#About)
- [API Documentation](#API-Documentation)
- [Getting Started](#Getting-Started)
	- [Prerequisites](#Prerequisites)
		- [Required Secrets](#Required-Secrets)
		- [Development](#Development)
	- [Installation](#Installation)
		- [Development](#Development)
		- [Production](#Production)
- [Versioning](#Versioning)
- [Authors](#Authors)

## About

This is a final year group project for the module
[CO600 at the University of Kent](https://www.kent.ac.uk/courses/modules/module/CO600).
The project uses virtual reality to show realtime bus locations in selected
areas.

This is the API server for the project. The client repository can be found on
[GitHub](https://github.com/University-of-Kent-VR-Transport/vr-client).

## API Documentation

API documentation can be found in the [docs](./docs)

## Getting Started

These instructions will get you a copy of the project up and running on your
local machine for development and testing purposes.

### Prerequisites

You'll need to install the following software:

```
Git v^2.0.0
Docker v^20.10.0
Docker Compose v^1.21.0
```

#### Required Secrets

Create an `.env` file in the root of the repository containing your secrets:
- `DFT_SECRET` [Department for Transport](https://data.bus-data.dft.gov.uk/account/settings/)
- `MAPBOX_TOKEN` [Mapbox](https://account.mapbox.com/access-tokens)
- `DB_PASSWORD` This is the main database password for the user root
- `DB_DOCKER_PASSWORD` This is the database password for the user docker
- `ADMIN_TOKEN` This is the admin token for protected routes.
eg. **`UPDATE`** `/api/bus-stop`

**Important: Do not commit your `.env` file**

#### Development

For development you'll also need:
```
go v1.15
```

### Installation

#### Development

A step by step series of examples that tell you how to get a development
environment running.

Once the repo has been cloned. Build the project:
```
go build -o bin/server server.go
```

To run the project using Docker Compose:

```
docker-compose up --detach --build
```

This will start all the images required for the project detached. This enables
the ability to rebuild the web service when you make a change in development:

```
go build -o bin/server server.go && docker-compose restart web
```

By default the project exposes port [`5050`](http://localhost:5050/).

To stop and tear down the containers run:

```
docker-compose down
```

#### Production

A step by step series of examples that tell you how to get a production
environment running.

Once the repo has been cloned. Build and run the project using docker-compose.
This will start all the images required for the project detached:

```
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up --build --detach --no-deps --remove-orphans
```

By default the project exposes port [`80`](http://localhost:80/).

To stop and tear down the containers run:

```
docker-compose down
```

## Versioning

We use [SemVer](https://semver.org/) for versioning. For the versions available,
see the
[tags on this repository](https://github.com/University-of-Kent-VR-Transport/vr-client/tags).

## Authors

- **Henry Brown** [HenryBrown0](https://github.com/HenryBrown0) `hb317@kent.ac.uk`
- **Joshua Lewis-Powell** [Wildcastle117](https://github.com/Wildcastle117) `jl715@kent.ac.uk`
- **Alex Fry** [the-dark-beat](https://github.com/the-dark-beat) `af491@kent.ac.uk`
