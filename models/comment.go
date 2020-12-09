package models

import (
	"log"
	"time"

	util "github.com/devstackq/ForumX/utils"
)

//Comment ID -> foreign key -> postID
type Comment struct {
	ID          int
	Content     string
	PostID      int
	UserID      int
	Author      string
	Like        int
	Dislike     int
	TitlePost   string
	Time        time.Time
	CreatedTime string
	ToWhom      int
}

//LeaveComment for post by id
func (c *Comment) LeaveComment() error {

	q, err := DB.Exec("INSERT INTO comments(content, post_id, creator_id, created_time) VALUES(?,?,?,?)",
		c.Content, c.PostID, c.UserID, time.Now())
	if err != nil {
		return err
	}
	//comnet contetn
	err = DB.QueryRow("SELECT creator_id FROM posts WHERE id=?", c.PostID).Scan(&c.ToWhom)
	if err != nil {
		log.Println(err)
	}
	lid, err := q.LastInsertId()
	if err != nil {
		log.Println(err)
	}
	//fmt.Println(c.ToWhom, "comment to whom lost")
	util.SetCommentNotify(c.PostID, c.UserID, c.ToWhom, lid)
	return nil
}

//UpdateComment func
func (c *Comment) UpdateComment() {

	_, err := DB.Exec("UPDATE  comments SET  content=? WHERE id =?",
		c.Content, c.ID)

	if err != nil {
		log.Println(err)
	}
}

// DeleteComment func
func DeleteComment(id string) {

	_, err = DB.Exec("DELETE FROM  comments  WHERE id =?", id)
	if err != nil {
		log.Println(err)
	}
	_, err = DB.Exec("DELETE FROM notify  WHERE comment_id =?", id)
	if err != nil {
		log.Println(err)
	}
	_, err = DB.Exec("DELETE FROM voteState  WHERE comment_id =?", id)
	if err != nil {
		log.Println(err)
	}
}
