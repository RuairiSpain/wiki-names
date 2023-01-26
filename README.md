# Simple WikiMedia API endpoint to search for people and output short description as JSON

## Features

- Docker file that hosts on PORT 8080
- All REST APIs (GET)

## Start server

```bash or zsh
go build
./wiki-names &
```

## API

[GIN-debug] GET /search/:name --> wiki-names/controllers.GetContentSummary (4 handlers)
[GIN-debug] GET /extract/:name --> wiki-names/controllers.GetExtract (4 handlers)
[GIN-debug] GET /extract/:name/:locale --> wiki-names/controllers.GetExtract (4 handlers)
[GIN-debug] GET /swagger/\*any --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (4 handlers)

### Get All Users

URL: GET localhost:8080/firstname_lastname

**Response:**

```json
{
  "short_description": "Leanne Graham"
}
```
