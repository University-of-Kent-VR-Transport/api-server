# Bus Routes

**/**  [docs/api](../)  **/**  [bus-routes](#Bus-Routes)

## Contents

- [Get](#GET)
- [Put](#PUT)
- [Options](#OPTIONS)

## GET

Returns a bus route by line name, direction and operator ID.

### Endpoint

**`GET`** `/api/bus-route`

### Query parameters

| Parameter  | Type             | Example  |
| ---------- | ---------------- | -------- |
| lineName   | string           | Uni1     |
| direction  | INBOUND/OUTBOUND | OUTBOUND |
| operatorID | string           | SCEK     |

### Example request

```curl
curl -X GET https://bus.henrybrown0.com/api/bus-routes?lineName=Uni1&direction=OUTBOUND&operatorID=SCEK
```

### Example Response

```json
{
    "Route": {
        "LineID": "SCEK:PK0000098:314_Uni1_Uni1V:Uni1:",
        "RouteID": "RT197",
        "OperatorID": "SCEK",
        "OperatorName": "Stagecoach in East Kent",
        "Name": "Uni1",
        "Direction": "OUTBOUND",
        "Description": "City Centre - University",
        "Origin": "Bus Station",
        "Destination": "Darwin College",
        "Stops": [
            {
                "ID": "240098906   ",
                "Name": "Bus Station",
                "Longitude": 1.0813389,
                "Latitude": 51.276302,
                "Bearing": 0
            },
						...
            {
                "ID": "240095612   ",
                "Name": "Darwin College",
                "Longitude": 1.0713621,
                "Latitude": 51.29914,
                "Bearing": 0
            }
        ]
    }
}
```

## PUT

Updates all bus routes within a dataset using the Department for Transport
timetable API. It returns a running [background job](./jobs.md#Get).

### Endpoint

**`PUT`** `/api/bus-routes/:datasetID`

### Authorization Header

You must provide the admin Authorization Bearer token for this request.

### Path parameters

| Parameter   | Type   | Example |
| ----------- | ------ | ------- |
| datasetID   | uint32 | 2022    |

### Example request

```curl
curl -H "Authorization: Bearer admin-token" -X PUT https://bus.henrybrown0.com/api/bus-routes/2022
```

### Example Response

```json
{
	"Job": {
		"ID": 3,
		"URI": "/api/job/3",
		"Type": "UPDATE ROUTES BY DATASET ID",
		"Status": "RUNNING",
		"CreatedAt": "2021-04-06T21:33:48.089822Z",
		"UpdatedAt": "2021-04-06T21:33:48.089822Z"
	}
}
```

## OPTIONS

Returns the options for the bus routes endpoint.

### Endpoint

**`OPTIONS`** `/api/bus-routes`

### Example request

```curl
curl -X OPTIONS https://bus.henrybrown0.com/api/bus-routes
```

### Example Response Header

| KEY             | Value                             |
| --------------- | --------------------------------- |
| Accept          | `application/json; charset=utf-8` |
| Accept-Encoding | `gzip`                            |
| Allow           | `GET, PUT, OPTIONS`               |