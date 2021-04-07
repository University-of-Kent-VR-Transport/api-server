# Background Jobs

**/**  [docs/api](../)  **/**  [jobs](#Background-Jobs)

## Contents

- [Get](#GET)
- [Options](#OPTIONS)

## GET

Returns the background job with the provided ID.

### Endpoint

**`GET`** `/api/job`

### Authorization Header

You must provide the admin Authorization Bearer token for this request.

### Path parameters

| Parameter   | Type   | Example |
| ----------- | ------ | ------- |
| jobID       | uint32 | 1       |

### Example request

```curl
curl -H "Authorization: Bearer admin-token" -X GET https://bus.henrybrown0.com/api/job/1
```

### Example Response

```json
{
	"Job": {
		"ID": 1,
		"URI": "/api/job/1",
		"Type": "UPDATE NATIONAL PUBLIC TRANSPORT ACCESS NODES",
		"Status": "RUNNING",
		"CreatedAt": "2021-04-06T21:33:48.089822Z",
		"UpdatedAt": "2021-04-06T21:33:48.089822Z"
	}
}
```

## OPTIONS

Returns the options for the jobs endpoint.

### Endpoint

**`OPTIONS`** `/api/job`

### Authorization Header

You must provide the admin Authorization Bearer token for this request.

### Example request

```curl
curl -H "Authorization: Bearer admin-token" -X OPTIONS https://bus.henrybrown0.com/api/job
```

### Example Response Header

| KEY             | Value                             |
| --------------- | --------------------------------- |
| Accept          | `application/json; charset=utf-8` |
| Accept-Encoding | `gzip`                            |
| Allow           | `GET, OPTIONS`                    |