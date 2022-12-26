package usecase

// EventUseCase -.
type EventUseCase struct {
	repo Storage
}

// New -.
func New(s Storage) *EventUseCase {
	return &EventUseCase{
		repo: s,
	}
}
