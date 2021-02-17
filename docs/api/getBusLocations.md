# Get Bus Locations
Returns all buses within a box defined by two coordinates.

## Endpoint
**`GET`** `/api/get-bus-locations`

## Query parameters
| Parameter   | Type                                              | Example         |
| ------------| ------------------------------------------------- | --------------- |
| topLeft     | Coordinate (Longitude float32, Latitude  float32) | 0.2654,80.2119  |
| bottomRight | Coordinate (Longitude float32, Latitude  float32) | 50.3020,-1.8579 |

## Example request
**`GET`** `/api/get-bus-locations?topLeft=0.2654,80.2119&bottomRight=50.3020,-1.8579`

## Example Response
```json
[
    {
        "ID": "878311f6-8c42-4267-b2ed-2ea9aaffb338",
        "Route": {
            "ID": "16",
            "Name": "16"
        },
        "Location": {
            "Longitude": 1.1774309,
            "Latitude": 51.07938
        },
        "Bearing": 222,
        "LastUpdated": "2021-02-17T10:20:07Z"
    },
    {
        "ID": "2d2f78fa-3b92-40f3-b074-d64432b4453b",
        "Route": {
            "ID": "6",
            "Name": "6"
        },
        "Location": {
            "Longitude": 1.1274358,
            "Latitude": 51.371384
        },
        "Bearing": 252,
        "LastUpdated": "2021-02-17T10:20:09Z"
    }
]
```