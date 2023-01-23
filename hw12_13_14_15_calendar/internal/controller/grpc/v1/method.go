package grpcv1

import (
	"context"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/utils"
	"go.uber.org/zap"
	"time"

	proto "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/api"
	errapp "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/controller/repo"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/entity"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CalendarGRPCService struct {
	u usecase.EventUseCase
	l zap.SugaredLogger
	proto.UnimplementedEventServiceServer
}

func NewCalendarGRPCService(u usecase.EventUseCase, l zap.SugaredLogger) *CalendarGRPCService {
	return &CalendarGRPCService{u: u, l: l}
}

func (g *CalendarGRPCService) CreateEvent(ctx context.Context, req *proto.Event) (*proto.Event, error) {
	result, err := g.u.Create(
		ctx,
		entity.Event{
			Title:        req.GetTitle(),
			Desc:         req.GetDesc(),
			UserID:       int(req.GetUserId()),
			StartTime:    req.GetStart(),
			EndTime:      req.GetEnd(),
			Notification: req.GetNotification(),
		},
	)
	if err != nil {
		if err == errapp.ErrEventTimeBusy {
			g.l.Error("grpc - v1 - create - ErrEventTimeBusy")
			return nil, status.Errorf(codes.InvalidArgument, "this event time is already busy")
		}
		g.l.Error(err, "grpc - v1 - create")
		return nil, status.Errorf(codes.Internal, "event creating problems")

	}

	return &proto.Event{
		Id:           int64(result.ID),
		Title:        result.Title,
		Desc:         result.Desc,
		UserId:       int64(result.UserID),
		Start:        result.StartTime,
		End:          result.EndTime,
		Notification: result.Notification,
	}, nil
}

func (g *CalendarGRPCService) UpdateEvent(ctx context.Context, req *proto.Event) (*proto.Event, error) {
	result, err := g.u.Update(
		ctx,
		entity.Event{
			ID:           int(req.GetId()),
			Title:        req.GetTitle(),
			Desc:         req.GetDesc(),
			UserID:       int(req.GetUserId()),
			StartTime:    req.GetStart(),
			EndTime:      req.GetEnd(),
			Notification: req.GetNotification(),
		},
	)
	if err != nil {
		if err == errapp.ErrEventTimeBusy {
			g.l.Error("grpc - v1 - update - ErrEventTimeBusy")
			return nil, status.Errorf(codes.InvalidArgument, "this event time is already busy")
		}
		g.l.Error(err, "grpc - v1 - update")
		return nil, status.Errorf(codes.Internal, "event updating problems")

	}

	return &proto.Event{
		Id:           int64(result.ID),
		Title:        result.Title,
		Desc:         result.Desc,
		UserId:       int64(result.UserID),
		Start:        result.StartTime,
		End:          result.EndTime,
		Notification: result.Notification,
	}, nil
}

func (g *CalendarGRPCService) DeleteEvent(ctx context.Context, uid *proto.UID) (*proto.Response, error) {
	err := g.u.Delete(ctx, int(uid.GetUid()))
	if err != nil {
		g.l.Error(err, "grpc - v1 - DeleteEvent")
		return &proto.Response{Status: "event deleting problems"}, status.Errorf(codes.Internal, "event deleting problems")
	}
	return &proto.Response{Status: "OK"}, nil
}

func (g *CalendarGRPCService) GetDailyEvents(ctx context.Context, in *proto.GetEventsRequest) (*proto.GetEventsResponse, error) {
	uid := int(in.GetUserId())
	start, err := dateconv.StringToDay(in.GetStart())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad event date")
	}

	end := start.Add(24 * time.Hour)

	result, err := g.u.EventsByDates(ctx, uid, start, end)
	if err != nil {
		g.l.Error("grpc - v1 - monthly - EventsByDates: %w", err)
		return nil, status.Errorf(codes.Internal, "getting daily events by date problems")
	}

	events := make([]*proto.Event, 0, len(result))
	for _, event := range result {
		events = append(events, &proto.Event{
			Id:           int64(event.ID),
			Title:        event.Title,
			Desc:         event.Desc,
			Start:        event.StartTime,
			End:          event.EndTime,
			Notification: event.Notification,
		})
	}

	return &proto.GetEventsResponse{Events: events}, nil
}

func (g *CalendarGRPCService) GetWeeklyEvents(ctx context.Context, in *proto.GetEventsRequest) (*proto.GetEventsResponse, error) {
	uid := int(in.GetUserId())
	start, err := dateconv.StringToDay(in.GetStart())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad event date")
	}

	end := start.Add(7 * 24 * time.Hour)

	result, err := g.u.EventsByDates(ctx, uid, start, end)
	if err != nil {
		g.l.Error("grpc - v1 - weekly - EventsByDates: %w", err)
		return nil, status.Errorf(codes.Internal, "getting weekly events by date problems")
	}

	events := make([]*proto.Event, 0, len(result))
	for _, event := range result {
		events = append(events, &proto.Event{
			Id:           int64(event.ID),
			Title:        event.Title,
			Desc:         event.Desc,
			Start:        event.StartTime,
			End:          event.EndTime,
			Notification: event.Notification,
		})
	}

	return &proto.GetEventsResponse{Events: events}, nil
}

func (g *CalendarGRPCService) GetMonthlyEvents(ctx context.Context, in *proto.GetEventsRequest) (*proto.GetEventsResponse, error) {
	uid := int(in.GetUserId())
	start, err := dateconv.StringToDay(in.GetStart())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad event date")
	}

	end := start.Add(30 * 24 * time.Hour)

	result, err := g.u.EventsByDates(ctx, uid, start, end)
	if err != nil {
		g.l.Error("grpc - v1 - monthly - EventsByDates: %w", err)
		return nil, status.Errorf(codes.Internal, "getting monthly events by date problems")
	}

	events := make([]*proto.Event, 0, len(result))
	for _, event := range result {
		events = append(events, &proto.Event{
			Id:           int64(event.ID),
			Title:        event.Title,
			Desc:         event.Desc,
			Start:        event.StartTime,
			End:          event.EndTime,
			Notification: event.Notification,
		})
	}

	return &proto.GetEventsResponse{Events: events}, nil
}

func (g *CalendarGRPCService) GetNotificationEvents(ctx context.Context, in *proto.Time) (*proto.GetEventsResponse, error) {
	start, err := dateconv.StringToTime(in.GetStart())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad event date")
	}

	result, err := g.u.EventsByTime(ctx, start)
	if err != nil {
		g.l.Error("grpc - v1 - notification - EventsByTime: %w", err)
		return nil, status.Errorf(codes.Internal, "getting events by time problems")
	}

	events := make([]*proto.Event, 0, len(result))
	for _, event := range result {
		events = append(events, &proto.Event{
			Id:           int64(event.ID),
			Title:        event.Title,
			Desc:         event.Desc,
			UserId:       int64(event.UserID),
			Start:        event.StartTime,
			End:          event.EndTime,
			Notification: event.Notification,
		})
	}

	return &proto.GetEventsResponse{Events: events}, nil
}
func (g *CalendarGRPCService) DeleteOldEvents(ctx context.Context, in *proto.Time) (*proto.Response, error) {
	start, err := dateconv.StringToDay(in.GetStart())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad start date")
	}

	e := g.u.DeleteOldEvents(ctx, start)
	if e != nil {
		g.l.Error(e, "grpc - v1 - DeleteOldEvent")
		return &proto.Response{Status: "old events deleting problems"}, status.Errorf(codes.Internal, "old events deleting problems")
	}
	return &proto.Response{Status: "OK"}, nil
}
