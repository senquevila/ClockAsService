# ClockAsService
Create a service to retrieve alarms and the elapsed time.

## Running the Service

### Prerequisites
- Go 1.20 or higher
- SQLite (handled via go-sqlite3)

### Install Dependencies
```sh
go mod download
```

### Run the Application

#### Option 1: Run directly with Go
```sh
go run src/api.go
```

#### Option 2: Build and run the binary
```sh
# Build the binary
go build -o clock-service src/api.go

# Run the binary
./clock-service
```

The server will start on `http://localhost:8080` and the SQLite database will be created as `clock.db` in the project directory.

### Run Tests
```sh
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test ./src/services
```

## Requirements
- Go 1.20+
- SQLite (handled via go-sqlite3)


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
Response: {
  "id": "alarm1",
  "countdown": 3599.99,
  "countdown_detailed": "59 minutes, 59 seconds"
}
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
Response: {
  "id": "event1",
  "elapsed": 2.01,
  "elapsed_detailed": "2 seconds"
}
```

## Notes
- All alarms and events are persisted in the SQLite database.
- Time values are in seconds and also provided in a human-readable format.
