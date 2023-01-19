package integration_test

import (
	"encoding/json"
	"github.com/Eun/go-hit"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/entity"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/utils/utils"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	// Attempts connection
	host       = "localhost:8080"
	healthPath = "http://" + host + "/healthz"
	attempts   = 20

	// HTTP REST
	basePath = "http://" + host + "/api/v1"
)

func TestMain(m *testing.M) {
	err := healthCheck(attempts)
	if err != nil {
		log.Fatalf("Integration tests: host %s is not available: %s", host, err)
	}

	log.Printf("Integration tests: host %s is available", host)

	code := m.Run()
	os.Exit(code)
}

func healthCheck(attempts int) error {
	var err error

	for attempts > 0 {
		err = hit.Do(hit.Get(healthPath), hit.Expect().Status().Equal(http.StatusOK))
		if err == nil {
			return nil
		}

		log.Printf("Integration tests: url %s is not available, attempts left: %d", healthPath, attempts)

		time.Sleep(time.Second)

		attempts--
	}

	return err
}

// HTTP POST: /event/create .
func TestHTTPUpdate(t *testing.T) {
	date := time.Now()
	event := entity.Event{
		ID:           7,
		Title:        "Test title",
		Desc:         "Test description",
		UserID:       42,
		StartTime:    utils.TimeToString(date.Add(time.Hour)),
		EndTime:      utils.TimeToString(date.Add(2 * time.Hour)),
		Notification: utils.TimeToString(date),
	}

	body, err := json.Marshal(event)
	if err != nil {
		log.Fatal("error marshal json: %w", err)
	}

	hit.Test(t,
		hit.Description("Update Event Success"),
		hit.Post(basePath+"/event/update"),
		hit.Send().Headers("Content-Type").Add("application/json"),
		hit.Send().Body().Bytes(body),
		hit.Expect().Status().Equal(http.StatusOK),
		hit.Expect().Body().JSON().JQ(".id").Equal("7.000000"),
		hit.Expect().Body().JSON().JQ(".title").Equal("Test title"),
		hit.Expect().Body().JSON().JQ(".desc").Equal("Test description"),
		hit.Expect().Body().JSON().JQ(".user_id").Equal("42.000000"),
	)
}

// HTTP GET: /event/daily .
func TestHTTPDaily(t *testing.T) {
	date := time.Now()
	start := date.Format("2006-01-02")
	hit.Test(t,
		hit.Description("Get Daily Events"),
		hit.Get(basePath+"/event/daily?uid=42&date="+start),
		hit.Expect().Status().Equal(http.StatusOK),
		hit.Expect().Body().JSON().JQ(".events[0].id").Equal("7.000000"),
	)
}
