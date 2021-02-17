# API-Server

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
Create an `.env` file in the root of the repository containing your DFT_SECRET.
A DFT_SECRET can be obtained from the 
[Department for Transport website](https://data.bus-data.dft.gov.uk/account/signup/).
**Important: Do not commit your `.env` file**

### Installation

A step by step series of examples that tell you how to get a development
environment running.

Once the repo has been cloned. Build and run the project using Docker Compose:

```
docker-compose up --detach --build
```

This will start all the images required for the project detached. This enables
the ability to rebuild the web service when you make a change in development:

```
docker-compose up --detach --build web
```

By default the project expose port [`5050`](http://localhost:5050/).

To stop the containers simple run:

```
docker-compose down
```

## Versioning

We use [SemVer](https://semver.org/) for versioning. For the versions available,
see the
[tags on this repository](https://github.com/University-of-Kent-VR-Transport/vr-client/tags).

## Authors

* **Henry Brown** [HenryBrown0](https://github.com/HenryBrown0) `hb317@kent.ac.uk`
* **Joshua Lewis-Powell** [Wildcastle117](https://github.com/Wildcastle117) `jl715@kent.ac.uk`
* **Alex Fry** [the-dark-beat](https://github.com/the-dark-beat) `af491@kent.ac.uk`
