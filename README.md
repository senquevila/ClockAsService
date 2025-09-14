# ClockAsService
Create a service to retrieve alarms and the elapsed time.

## API Endpoints

The service runs on port 8080 by default.

### Create an Alarm
```
POST /alarms/create
Content-Type: application/json
{
  "id": "alarm1",
  "target": "2025-09-07T12:00:00Z"
}
```

### Get Alarm Countdown
```
GET /alarms/countdown?id=alarm1
Response: { "id": "alarm1", "countdown": 3599.99 }
```

### Create an Event
```
POST /events/create
Content-Type: application/json
{
  "id": "event1"
}
```

### Get Event Elapsed Time
```
GET /events/elapsed?id=event1
Response: { "id": "event1", "elapsed": 2.01 }
```

## Running the Service

1. Build and run the server:
   ```sh
   go run src/api.go
   ```
2. The SQLite database will be created as `clock.db` in the project directory.

## Requirements
- Go 1.20+
- SQLite (handled via go-sqlite3)

## Notes
- All alarms and events are persisted in the SQLite database.
- Time values are in seconds.
