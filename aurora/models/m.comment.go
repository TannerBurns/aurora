package models

type Comment struct {
	ID       int        `json:"id"`
	Created  NullTime   `json:"created"`
	Modified NullTime   `json:"modified"`
	OwnerID  int        `json:"owner_id"`
	TaskID   int        `json:"task_id"`
	Body     NullString `json:"body"`
	Tags     []Tag      `json:"tags"`
}
