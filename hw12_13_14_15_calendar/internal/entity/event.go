package entity

type (
	Event struct {
		ID           int    `json:"id"`
		Title        string `json:"title"`
		Desc         string `json:"desc"`
		UserID       int    `json:"user_id"`
		StartTime    string `json:"start_time"`
		EndTime      string `json:"end_time"`
		Notification string `json:"notification"`
	}
)
