package models

import (
	"time"
)

//comment ID -> foreign key -> postID
type Comment struct {
	ID          int
	Content     string
	PostID      int
	UserID      int
	CreatedTime time.Time
	Author      string
	Like        int
	Dislike     int
	TitlePost   string
}

//get data from client, put data in Handler, then models -> query db
func (c *Comment) LeaveComment() error {
	_, err := DB.Exec("INSERT INTO comments( content, post_id, user_idx) VALUES(?,?,?)",
		c.Content, c.PostID, c.UserID)
	if err != nil {
		return err
	}
	return nil
}
