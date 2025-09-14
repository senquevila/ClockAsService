package main

	"encoding/json"
	"net/http"
	"time"
	"ClockAsService/src/common"
	"github.com/google/uuid"
)

var service *common.Service

// AlarmRequest is used for creating/updating alarms
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Target      time.Time `json:"target"`
}

// EventRequest is used for creating/updating events
	Name        string `json:"name"`
	Description string `json:"description"`
}

func createAlarmHandler(w http.ResponseWriter, r *http.Request) {
	var req AlarmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	id := uuid.New().String()
	if err := service.CreateAlarm(id, req.Name, req.Description, req.Target); err != nil {
		http.Error(w, "Failed to create alarm", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func getAlarmCountdownHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if countdown, ok := service.GetAlarmCountdown(id); ok {
		json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "countdown": countdown.Seconds()})
	} else {
		http.Error(w, "Alarm not found", http.StatusNotFound)
	}
}

func createEventHandler(w http.ResponseWriter, r *http.Request) {
	var req EventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	id := uuid.New().String()
	if err := service.CreateEvent(id, req.Name, req.Description); err != nil {
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func getEventElapsedHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if elapsed, ok := service.GetEventElapsed(id); ok {
		json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "elapsed": elapsed.Seconds()})
	} else {
		http.Error(w, "Event not found", http.StatusNotFound)
	}
}

func main() {
	var err error
	service, err = common.NewService("clock.db")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/alarms/create", createAlarmHandler)
	http.HandleFunc("/alarms/countdown", getAlarmCountdownHandler)
	http.HandleFunc("/events/create", createEventHandler)
	http.HandleFunc("/events/elapsed", getEventElapsedHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
