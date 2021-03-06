# Religious Ecologies Data API

This repository provides an API to access data stored in a PostgreSQL database. It is a component of the [American Religious Ecologies](http://religiousecologies.org) project and other projects at the [Roy Rosenzweig Center for History and New Media](https://rrchnm.org).

## Endpoints

The following endpoints are offered by the API.

### AHCB counties and states

Spatial polygons from the [Atlas of Historic County Boundaries](https://publications.newberry.org/ahcbp/) are available by date. The results will always be filtered by the date provided. Use the `id`, `state-terr-id`, or `state-code` to filter geographically.

```
GET /ahcb/states/:date/
GET /ahcb/counties/:date/
GET /ahcb/counties/:date/id/:id/
GET /ahcb/counties/:date/state-terr-id/:state_terr_id/
GET /ahcb/counties/:date/state-code/:state_code/
```

Parameters:

- `date`: The date of the historic boundaries, specified as an ISO-8601 string (e.g., `1848-07-05`). If the date requested is before or after the minimum or after the maximum dates for that type of geometry, the minimum or maximum will be silently returned.
- `id`: A comma-separated list of AHCB IDs for counties (e.g., `vas_fairfax`).
- `state_terr_id`: A comma-separated list of AHCB IDs for states and territories (e.g., `va_state`).
- `state_code`: A comma-separated list two-letter codes for states, roughly corresponding to postal codes (.e.g, `va`).


Response:

A GeoJSON feature collection in EPSG 4326 with one feature per state or county. The properties of each feature will include information such as the square mileage of the polygons.

### Catholic dioceses

Historical spatial point locations for Catholic dioceses in the United States, Canada, and Mexico. 

```
GET /catholic-dioceses/
```

Parameters:

- None

Response:

A JSON array of objects, each object representing a Catholic diocese. The date that the diocese was created (`date_erected`) is an ISO-8601 string, and the date that the diocese became an archdiocese (`date_metropolitan`) is either an ISO-8601 string or an empty string if the diocese did not become an archdiocese.

### Catholic dioceses per decade

Counts of the number of dioceses established per decade

```
GET /catholic-dioceses/per-decade/
```

Parameters:

- None

Response:

A JSON array of objects, each object representing the count of dioceses per decade.

### North America

Country polygons from [Natural Earth](https://www.naturalearthdata.com) for North America. 

```
GET /ne/northamerica/
```

Parameters:

- None

Response:

A GeoJSON feature collection in EPSG 4326 with one feature per country. 

### Presbyterians

Presbyterian membership data per year for 1826-1926.

```
GET /presbyterians/
```

Parameters:

- None

Response:

A JSON array of objects, each object representing a year of membership data.

### Populated places

Lists of populated places from the USGS [Geographic Names Information System](https://www.usgs.gov/core-science-systems/ngp/board-on-geographic-names/domestic-names) in a state or county.

```
GET /pop-places/state/:state/county/
GET /pop-places/state/:state/place/
GET /pop-places/county/:county/place/
GET /pop-places/place/:place_id/
```

Parameters:

- `state`: The state as a two-character abbreviation (e.g., `CA`).
- `county`: An ACHB county ID (e.g., `vas_fairfax`).
- `place_id`: The ID for a place (e.g., `611119`).

Response:

A JSON array of objects where the objrects are all the counties in a state, or all the places in a state or county. The objects have both the unique identifiers and human-readable names.

For example, you should be able to query `/pop-places/state/va/county` to get all the counties in Virginia, and then query `/pop-places/county/vas_fairfax/place` to get all the places in Fairfax county. If you had the ID for a place like Centreville, you could look up its details by querying `/pop-places/place/1491083/`.

### Endpoints

```
GET /
```

Parameters:

- None

Response:

A JSON array of objects. Each object is an endpoint for the API, with a sample URL for that endpoint.

## Compiling

Using a version of Go that supports Go modules, you should be able to run `go build` in the project root to install dependencies.

There is a Makefile in `cmd/dataapi/` that can be used for compiling and for running the service locally.

## Configuration

Set the following environment variables to configure the server:

- `DATAAPI_DBHOST` (default: `localhost`)
- `DATAAPI_DBPORT` (default: `5432`)
- `DATAAPI_DBNAME` (default: none)
- `DATAAPI_DBUSER` (default: none)
- `DATAAPI_DBPASS` (default: none)
- `DATAAPI_SSL` (default: `require`; see [pq docs](https://godoc.org/github.com/lib/pq))
- `DATAAPI_INTERFACE` (default: `0.0.0.0`)
- `DATAAPI_PORT` (default: `8090`)
- `DATAAPI_LOGGING` (default: `on`)

If logging is on, then access logs will be written to stdout in the Apache Common Log format. Errors and status messages will always be written to stderr.

Obviously this service requires that you be able to access the database.

## Testing

You can run the tests with `go test -v` inside the directory that contains the package for the command: `cmd/dataapi`.
