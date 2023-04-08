package entity

type (
	Event struct {
		ID           int    `json:"id"`
		Title        string `json:"title"`
		Desc         string `json:"desc"`
		UserID       int    `json:"userId"`
		StartTime    string `json:"startTime"`
		EndTime      string `json:"endTime"`
		Notification string `json:"notification"`
	}
)
