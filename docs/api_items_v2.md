````markdown
# mSpace DSpace-like API Documentation

Base URL: `http://localhost:8080`

---

## Health Check

**Endpoint:** `GET /health`

```bash
curl -X GET http://localhost:8080/health
````

**Response:**

```json
{
  "status": "ok"
}
```

---

## Communities

### List Communities

**Endpoint:** `GET /api/communities`

```bash
curl -X GET http://localhost:8080/api/communities
```

**Response Example:**

```json
[
  {
    "id": 1,
    "name": "Computer Science",
    "description": "CS community",
    "collections": []
  }
]
```

### Create Community

**Endpoint:** `POST /api/communities`

```bash
curl -X POST http://localhost:8080/api/communities \
  -H "Content-Type: application/json" \
  -d '{"name": "Computer Science", "description": "CS community"}'
```

**Response:** Newly created community JSON.

---

## Collections

### List Collections

**Endpoint:** `GET /api/collections`

```bash
curl -X GET http://localhost:8080/api/collections
```

### Create Collection

**Endpoint:** `POST /api/collections`

```bash
curl -X POST http://localhost:8080/api/collections \
  -H "Content-Type: application/json" \
  -d '{"name": "Algorithms", "description": "Algorithm research", "community_id": 1}'
```

**Response:** Newly created collection JSON.

---

## Items

### List Items

**Endpoint:** `GET /api/items`

```bash
curl -X GET http://localhost:8080/api/items
```

### Create Item

**Endpoint:** `POST /api/items`

```bash
curl -X POST http://localhost:8080/api/items \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Sorting Algorithms",
    "author": "Alice",
    "abstract": "Study on sorting algorithms",
    "collection_id": 1,
    "visibility": "PUBLIC"
}'
```

**Response:** Newly created item JSON with `Status: DRAFT`.

### Get Item by ID

**Endpoint:** `GET /api/items/{id}`

```bash
curl -X GET http://localhost:8080/api/items/1
```

### Upload File to Item

**Endpoint:** `POST /api/items/{id}/file`

```bash
curl -X POST http://localhost:8080/api/items/1/file \
  -F "file=@/path/to/file.pdf"
```

**Response:**

```json
{
  "message": "file uploaded",
  "file_url": "http://localhost:9000/repository/item-1-v1-1694000000.pdf"
}
```

### Publish Item

**Endpoint:** `POST /api/items/{id}/publish`

```bash
curl -X POST http://localhost:8080/api/items/1/publish
```

**Response:** Item JSON with `Status: PUBLISHED`.

### Reject Item

**Endpoint:** `POST /api/items/{id}/reject`

```bash
curl -X POST http://localhost:8080/api/items/1/reject
```

**Response:** Item JSON with `Status: REJECTED`.

---

## Metadata

### Add Metadata to Item

**Endpoint:** `POST /api/items/{id}/metadata`

```bash
curl -X POST http://localhost:8080/api/items/1/metadata \
  -H "Content-Type: application/json" \
  -d '{"key": "keywords", "value": "golang, repository"}'
```

### Get Metadata of Item

**Endpoint:** `GET /api/items/{id}/metadata`

```bash
curl -X GET http://localhost:8080/api/items/1/metadata
```

**Response Example:**

```json
[
  {"key": "keywords", "value": "golang, repository"}
]
```

---

> **Notes:**
>
> * Ensure MinIO is running and the bucket exists for file uploads.
> * Items must be published to be publicly visible.
> * Collections must belong to a community.
> * Supports versioning for files; each upload increments the version.

```
```
