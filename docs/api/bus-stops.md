# Bus Stops

**/**  [docs/api](../)  **/**  [bus-stops](#Bus-Stops)

## Contents

- [Get](#GET)
- [Put](#PUT)
- [Options](#OPTIONS)

## GET

Returns all buses stops within a box defined by min/max coordinates.

### Endpoint

**`GET`** `/api/bus-stops`

### Query parameters

| Parameter    | Type    | Example |
| ------------ | ------- | ------- |
| minLongitude | float32 | 1.0511  |
| minLatitude  | float32 | 51.2672 |
| maxLongitude | float32 | 1.1207  |
| maxLatitude  | float32 | 51.2943 |

### Example request

```curl
curl -X GET https://bus.henrybrown0.com/api/bus-stops?minLongitude=1.0511&minLatitude=51.2672&maxLongitude=1.1207&maxLatitude=51.2943
```

### Example Response

```json
{
	"BusStops": [
		{
			"ID": "2400105752",
			"Name": "St Dunstan's Church",
			"Longitude": 1.0704081,
			"Latitude": 51.28378,
			"Bearing": 225
		},
		{
			"ID": "2400100621",
			"Name": "Hanscomb House",
			"Longitude": 1.0712043,
			"Latitude": 51.28496,
			"Bearing": 45
		},
	]
}
```

## PUT

Updates all bus stops using the Department for Transport National Public
Transport Access Node database. It returns a running
[background job](./jobs.md#Get).

### Endpoint

**`PUT`** `/api/bus-stops`

### Authorization Header

You must provide the admin Authorization Bearer token for this request.

### Example request

```curl
curl -H "Authorization: Bearer admin-token" -X PUT https://bus.henrybrown0.com/api/bus-stops
```

### Example Response

```json
{
	"Job": {
		"ID": 3,
		"URI": "/api/job/3",
		"Type": "UPDATE NATIONAL PUBLIC TRANSPORT ACCESS NODES",
		"Status": "RUNNING",
		"CreatedAt": "2021-04-06T21:33:48.089822Z",
		"UpdatedAt": "2021-04-06T21:33:48.089822Z"
	}
}
```

## OPTIONS

Returns the options for the bus stops endpoint.

### Endpoint

**`OPTIONS`** `/api/bus-stops`

### Example request

```curl
curl -X OPTIONS https://bus.henrybrown0.com/api/bus-stops
```

### Example Response Header

| KEY             | Value                             |
| --------------- | --------------------------------- |
| Accept          | `application/json; charset=utf-8` |
| Accept-Encoding | `gzip`                            |
| Allow           | `GET, PUT, OPTIONS`               |