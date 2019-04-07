package models

type Tag struct {
	ID        int        `json:"id"`
	Created   NullTime   `json:"created"`
	OwnerID   int        `json:"owner_id"`
	CommentID int        `json:"comment_id"`
	Body      NullString `json:"body"`
}
