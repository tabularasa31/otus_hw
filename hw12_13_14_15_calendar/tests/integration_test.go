package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	gohit "github.com/Eun/go-hit"
	"github.com/stretchr/testify/suite"
	proto "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/api"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/entity"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/utils/dateconv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
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

// HTTP testing
func TestHTTP(t *testing.T) {
	date := time.Now()
	event := entity.Event{
		Title:        "Test title",
		Desc:         "Test description",
		UserID:       42,
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
		gohit.Expect().Body().JSON().JQ(".user_id").Equal("42.000000"),
	)

	start := date.Format("2006-01-02")
	gohit.Test(t,
		gohit.Description("Get Daily Events"),
		gohit.Get(basePath+"/event/daily?uid=42&date="+start),
		gohit.Expect().Status().Equal(http.StatusOK),
		gohit.Expect().Body().JSON().JQ(".events[0].user_id").Equal("42.000000"),
	)
}

// GRPC testing
type CalSuite struct {
	suite.Suite
	ctx    context.Context
	conn   *grpc.ClientConn
	client proto.EventServiceClient
}

func (s *CalSuite) SetupSuite() {
	grpcHost := os.Getenv("GRPC_ADDR")
	if grpcHost == "" {
		grpcHost = "localhost:50051"
	}
	var err error
	s.conn, err = grpc.Dial(grpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)

	s.ctx = context.Background()
	s.client = proto.NewEventServiceClient(s.conn)
}

type testCase struct {
	description string
	request     *proto.Event
	response    *proto.Event
	expectedErr error
}

func (s *CalSuite) Test_CreateEvent() {
	date := time.Now()
	testCases := []testCase{
		{
			description: "success case",
			request: &proto.Event{
				Title:        "Test title",
				Desc:         "Test description",
				UserId:       43,
				Start:        dateconv.TimeToString(date.Add(time.Hour)),
				End:          dateconv.TimeToString(date.Add(2 * time.Hour)),
				Notification: dateconv.TimeToString(date),
			},
			response: &proto.Event{
				Title:        "Test title",
				Desc:         "Test description",
				UserId:       43,
				Start:        dateconv.TimeToString(date.Add(time.Hour)),
				End:          dateconv.TimeToString(date.Add(2 * time.Hour)),
				Notification: dateconv.TimeToString(date),
			},
			expectedErr: nil,
		},
		{
			description: "empty user id",
			request: &proto.Event{
				Title:        "Test title",
				Desc:         "Test description",
				Start:        dateconv.TimeToString(date.Add(time.Hour)),
				End:          dateconv.TimeToString(date.Add(2 * time.Hour)),
				Notification: dateconv.TimeToString(date),
			},
			response:    &proto.Event{},
			expectedErr: status.Errorf(codes.InvalidArgument, "empty event user id"),
		},
		{
			description: "empty start time",
			request: &proto.Event{
				Title:        "Test title",
				Desc:         "Test description",
				UserId:       44,
				End:          dateconv.TimeToString(date.Add(2 * time.Hour)),
				Notification: dateconv.TimeToString(date),
			},
			response:    &proto.Event{},
			expectedErr: status.Errorf(codes.InvalidArgument, "empty event time"),
		},
	}

	for _, tc := range testCases {
		fmt.Printf("\n Test Case: %s \n", tc.description)
		resp, err := s.client.CreateEvent(s.ctx, tc.request)
		s.Require().Equal(tc.response.GetTitle(), resp.GetTitle())
		s.Require().Equal(tc.response.GetDesc(), resp.GetDesc())
		s.Require().Equal(tc.expectedErr, err)
	}
}

func TestCalSuite(t *testing.T) {
	suite.Run(t, new(CalSuite))
}
