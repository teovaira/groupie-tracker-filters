# API Contract

Base URL: `https://groupietrackers.herokuapp.com/api`

---

## GET /artists

Returns a list of all artists.

**Response** `200 OK`
```json
[
  {
    "id": 1,
    "name": "Queen",
    "image": "https://groupietrackers.herokuapp.com/api/images/queen.jpeg",
    "members": ["Freddie Mercury", "Brian May", "Roger Taylor", "John Deacon"],
    "creationDate": 1970,
    "firstAlbum": "14-12-1973",
    "locations": "https://groupietrackers.herokuapp.com/api/locations/1",
    "concertDates": "https://groupietrackers.herokuapp.com/api/dates/1",
    "relations": "https://groupietrackers.herokuapp.com/api/relation/1"
  }
]
```

| Field          | Type       | Description                        |
|----------------|------------|------------------------------------|
| `id`           | `int`      | Unique artist identifier           |
| `name`         | `string`   | Artist or band name                |
| `image`        | `string`   | URL to artist image                |
| `members`      | `[]string` | List of band members               |
| `creationDate` | `int`      | Year the band was formed           |
| `firstAlbum`   | `string`   | Date of first album (DD-MM-YYYY)   |
| `locations`    | `string`   | URL to artist's locations resource |
| `concertDates` | `string`   | URL to artist's dates resource     |
| `relations`    | `string`   | URL to artist's relations resource |

---

## GET /locations

Returns all concert locations grouped by artist.

**Response** `200 OK`
```json
{
  "index": [
    {
      "id": 1,
      "locations": ["london-uk", "paris-france"]
    }
  ]
}
```

| Field       | Type       | Description                      |
|-------------|------------|----------------------------------|
| `index`     | `[]object` | List of location entries         |
| `id`        | `int`      | Artist ID                        |
| `locations` | `[]string` | List of concert location strings |

---

## GET /dates

Returns all concert dates grouped by artist.

**Response** `200 OK`
```json
{
  "index": [
    {
      "id": 1,
      "dates": ["*06-03-2020", "07-03-2020"]
    }
  ]
}
```

| Field   | Type       | Description                                                       |
|---------|------------|-------------------------------------------------------------------|
| `index` | `[]object` | List of date entries                                              |
| `id`    | `int`      | Artist ID                                                         |
| `dates` | `[]string` | List of concert dates (DD-MM-YYYY), past dates prefixed with `*`  |

---

## GET /relation

Returns the mapping of locations to dates for each artist.

**Response** `200 OK`
```json
{
  "index": [
    {
      "id": 1,
      "datesLocations": {
        "london-uk": ["06-03-2020", "07-03-2020"],
        "paris-france": ["10-03-2020"]
      }
    }
  ]
}
```

| Field            | Type                  | Description                              |
|------------------|-----------------------|------------------------------------------|
| `index`          | `[]object`            | List of relation entries                 |
| `id`             | `int`                 | Artist ID                                |
| `datesLocations` | `map[string][]string` | Map of location to list of concert dates |

---

## GET /api/search?q=

Search artists by name, member, location, or creation year.

**Query Parameters**

| Parameter | Type     | Required | Description      |
|-----------|----------|----------|------------------|
| `q`       | `string` | yes      | Search term      |

**Response** `200 OK`
```json
[
  {
    "id": 1,
    "name": "Queen",
    "image": "https://groupietrackers.herokuapp.com/api/images/queen.jpeg",
    "creationDate": 1970
  }
]
```

| Field          | Type     | Description              |
|----------------|----------|--------------------------|
| `id`           | `int`    | Unique artist identifier |
| `name`         | `string` | Artist or band name      |
| `image`        | `string` | URL to artist image      |
| `creationDate` | `int`    | Year the band was formed |

**Status Codes**

| Code | Meaning                           |
|------|-----------------------------------|
| 200  | Success — results or empty array  |
| 400  | Missing q parameter               |
| 500  | Internal server error             |

**Notes**
- Returns empty array `[]` when no artists match — never `null`
- Field names are exact — do not rename without team sync
