# Simple WikiMedia API endpoint to search for people and output short description as JSON.

## Features

- Docker file that hosts on PORT 8080
- All REST APIs (GET)

## Start server

```bash
air run main.go
```

## API

### Get All Users

URL: GET localhost:8080/firstname_lastname

**Response:**

```json
[
  {
    "short_description": "Leanne Graham"
  }
]
```
