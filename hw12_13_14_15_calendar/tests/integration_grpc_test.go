package integration_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	proto "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/api"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/utils/dateconv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

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
	uid := rand.Int63n(100)

	testCases := []testCase{
		{
			description: "success create event",
			request: &proto.Event{
				Title:        "Test title",
				Desc:         "Test description",
				UserId:       uid,
				Start:        dateconv.TimeToString(date.Add(time.Hour)),
				End:          dateconv.TimeToString(date.Add(2 * time.Hour)),
				Notification: dateconv.TimeToString(date),
			},
			response: &proto.Event{
				Title:        "Test title",
				Desc:         "Test description",
				UserId:       uid,
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
			expectedErr: status.Errorf(codes.InvalidArgument, "empty user id"),
		},
		{
			description: "empty start time",
			request: &proto.Event{
				Title:        "Test title",
				Desc:         "Test description",
				UserId:       uid + 1,
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
		s.Require().Equal(tc.expectedErr, err)
		s.Require().Equal(tc.response.GetTitle(), resp.GetTitle())
		s.Require().Equal(tc.response.GetDesc(), resp.GetDesc())
	}

	_, err := s.client.DeleteByUserID(s.ctx, &proto.UID{Uid: uid})
	s.Require().Equal(nil, err)
}

func (s *CalSuite) Test_UpdateEvent() {
	date := time.Now()
	uid := rand.Int63n(100)

	requestCreate := &proto.Event{
		Title:        "Test title",
		Desc:         "Test description",
		UserId:       uid,
		Start:        dateconv.TimeToString(date.Add(time.Hour)),
		End:          dateconv.TimeToString(date.Add(2 * time.Hour)),
		Notification: dateconv.TimeToString(date),
	}

	response := &proto.Event{
		Title:        "Updated title",
		Desc:         "Updated description",
		UserId:       uid,
		Start:        dateconv.TimeToString(date.Add(time.Hour)),
		End:          dateconv.TimeToString(date.Add(2 * time.Hour)),
		Notification: dateconv.TimeToString(date),
	}

	fmt.Printf("\n Test Case: success update event \n")
	temp, err := s.client.CreateEvent(s.ctx, requestCreate)
	s.Require().Equal(nil, err)

	temp.Title = response.Title
	temp.Desc = response.Desc
	resp, err := s.client.UpdateEvent(s.ctx, temp)
	s.Require().Equal(nil, err)
	s.Require().Equal(response.GetTitle(), resp.GetTitle())
	s.Require().Equal(response.GetDesc(), resp.GetDesc())

	fmt.Printf("\n Test Case: failed update event \n")
	temp.UserId = 0
	_, err = s.client.UpdateEvent(s.ctx, temp)
	s.Require().Equal(status.Errorf(codes.InvalidArgument, "empty user id"), err)

	_, err = s.client.DeleteByUserID(s.ctx, &proto.UID{Uid: uid})
	s.Require().Equal(nil, err)
}

func (s *CalSuite) Test_DeleteEvent() {
	date := time.Now()
	requestCreate := &proto.Event{
		Title:        "Test title",
		Desc:         "Test description",
		UserId:       45,
		Start:        dateconv.TimeToString(date.Add(time.Hour)),
		End:          dateconv.TimeToString(date.Add(2 * time.Hour)),
		Notification: dateconv.TimeToString(date),
	}

	fmt.Printf("\n Test Case: success delete event \n")
	temp, err := s.client.CreateEvent(s.ctx, requestCreate)
	s.Require().Equal(nil, err)

	id := &proto.ID{Id: temp.GetId()}
	resp, err := s.client.DeleteEvent(s.ctx, id)
	s.Require().Equal(nil, err)
	s.Require().Equal("OK", resp.Status)
}

func (s *CalSuite) Test_GetDailyEvents() {
	date := time.Now()
	uid := rand.Int63n(100)
	fmt.Printf("\n Test Case: success get daily events \n")
	for i := 0; i < 5; i++ {
		_, err := s.client.CreateEvent(s.ctx, &proto.Event{
			Title:        "Test title",
			Desc:         "Test description",
			UserId:       uid,
			Start:        dateconv.TimeToString(date.Add(time.Hour)),
			End:          dateconv.TimeToString(date.Add(1 * time.Hour)),
			Notification: dateconv.TimeToString(date),
		})
		s.Require().Equal(nil, err)
	}

	events, err := s.client.GetDailyEvents(s.ctx, &proto.GetEventsRequest{UserId: uid, Start: date.Format("2006-01-02")})
	s.Require().Equal(nil, err)
	s.Require().Equal(5, len(events.Events))

	_, err = s.client.DeleteByUserID(s.ctx, &proto.UID{Uid: uid})
	s.Require().Equal(nil, err)
}

func TestCalSuite(t *testing.T) {
	suite.Run(t, new(CalSuite))
}
