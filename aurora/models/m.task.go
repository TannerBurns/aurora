package models

type Task struct {
	ID         int        `json:"id"`
	Created    NullTime   `json:"created"`
	Modified   NullTime   `json:"modified"`
	OwnerID    int        `json:"owner_id"`
	Restricted bool       `json:"restricted"`
	Status     NullString `json:"status"`
	Title      NullString `json:"title"`
	Body       NullString `json:"body"`
	Comments   []Comment  `json:"comments"`
}
