# Health Check

**/**  [docs/api](../)  **/**  [health-check](#Health-Check)

## Contents

- [Get](#GET)
- [Options](#OPTIONS)

## GET

Returns the health of the service.

### Endpoint

**`GET`** `/api/health-check`

### Example request

```curl
curl -X GET https://bus.henrybrown0.com/api/health-check
```

### Example Response

```json
{
  "database": true
}
```

## OPTIONS

Returns the options for the health check endpoint.

### Endpoint

**`OPTIONS`** `/api/bus-locations`

### Example request

```curl
curl -X OPTIONS https://bus.henrybrown0.com/api/bus-locations
```

### Example Response Header

| KEY             | Value                             |
| --------------- | --------------------------------- |
| Accept          | `application/json; charset=utf-8` |
| Accept-Encoding | `gzip`                            |
| Allow           | `GET, OPTIONS`                    |