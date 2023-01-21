package integration_test

import (
	"encoding/json"
	gohit "github.com/Eun/go-hit"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/entity"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/utils/utils"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
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
		err = gohit.Do(gohit.Get(healthPath), gohit.Expect().Status().Equal(http.StatusOK))
		if err == nil {
			return nil
		}

		log.Printf("Integration tests: url %s is not available, attempts left: %d", healthPath, attempts)

		time.Sleep(time.Second)

		attempts--
	}

	return err
}

// HTTP POST: /event/update .
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

	gohit.Test(t,
		gohit.Description("Update Event Success"),
		gohit.Post(basePath+"/event/update"),
		gohit.Send().Headers("Content-Type").Add("application/json"),
		gohit.Send().Body().Bytes(body),
		gohit.Expect().Status().Equal(http.StatusOK),
		gohit.Expect().Body().JSON().JQ(".id").Equal("7.000000"),
		gohit.Expect().Body().JSON().JQ(".title").Equal("Test title"),
		gohit.Expect().Body().JSON().JQ(".desc").Equal("Test description"),
		gohit.Expect().Body().JSON().JQ(".user_id").Equal("42.000000"),
	)
}

// HTTP GET: /event/daily .
func TestHTTPDaily(t *testing.T) {
	date := time.Now()
	start := date.Format("2006-01-02")
	gohit.Test(t,
		gohit.Description("Get Daily Events"),
		gohit.Get(basePath+"/event/daily?uid=42&date="+start),
		gohit.Expect().Status().Equal(http.StatusOK),
		gohit.Expect().Body().JSON().JQ(".events[0].id").Equal("7.000000"),
	)
}
