package todo

import "time"

type Todo struct {
	ID         int64
	UserID     int64
	Title      string
	IsDone     bool
	CreatedAt  time.Time
	ModifiedAt time.Time
}

type CreateTodoRequest struct {
	Title string `json:"title"`
}

type UpdateTodoRequest struct {
	Title  *string `json:"title"`
	IsDone *bool   `json:"is_done"`
}
