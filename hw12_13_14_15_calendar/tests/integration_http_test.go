package integration_test

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	gohit "github.com/Eun/go-hit"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/entity"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/utils/dateconv"
)

const (
	attempts = 5
)

func getHost() string {
	host := os.Getenv("HTTP_ADDR")
	if host == "" {
		host = "localhost:8080"
	}
	return host
}

func TestMain(m *testing.M) {
	host := getHost()
	healthPath := "http://" + host + "/healthz"

	err := healthCheck(healthPath, attempts)
	if err != nil {
		log.Fatalf("Integration tests: host %s is not available: %s", host, err)
	}

	log.Printf("Integration tests: host %s is available", host)

	code := m.Run()
	os.Exit(code)
}

func healthCheck(healthPath string, attempts int) error {
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

// HTTP testing
func TestHTTP(t *testing.T) {
	host := getHost()
	basePath := "http://" + host + "/api/v1"

	date := time.Now()
	uid := rand.Intn(100)

	event := entity.Event{
		Title:        "Test title",
		Desc:         "Test description",
		UserID:       uid,
		StartTime:    dateconv.TimeToString(date.Add(time.Hour)),
		EndTime:      dateconv.TimeToString(date.Add(2 * time.Hour)),
		Notification: dateconv.TimeToString(date),
	}

	body, err := json.Marshal(event)
	if err != nil {
		log.Fatal("error marshal json: %w", err)
	}

	gohit.Test(t,
		gohit.Description("Create Event Success"),
		gohit.Post(basePath+"/event/create"),
		gohit.Send().Headers("Content-Type").Add("application/json"),
		gohit.Send().Body().Bytes(body),
		gohit.Expect().Status().Equal(http.StatusCreated),
		gohit.Expect().Body().JSON().JQ(".title").Equal("Test title"),
		gohit.Expect().Body().JSON().JQ(".desc").Equal("Test description"),
		gohit.Expect().Body().JSON().JQ(".userId").Equal(uid),
	)

	start := date.Format("2006-01-02")
	u := strconv.Itoa(uid)
	gohit.Test(t,
		gohit.Description("Get Daily Events"),
		gohit.Get(basePath+"/event/daily?uid="+u+"&date="+start),
		gohit.Expect().Status().Equal(http.StatusOK),
		gohit.Expect().Body().JSON().JQ(".events[0].userId").Equal(uid),
	)

	gohit.Test(t,
		gohit.Delete(basePath+"/event/deletebyuid/"+u),
		gohit.Expect().Status().Equal(http.StatusOK),
	)
}
