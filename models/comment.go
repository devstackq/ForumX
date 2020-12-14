package models

import (
	"log"
	"time"

	util "github.com/devstackq/ForumX/utils"
)

//Comment ID -> foreign key -> postID
type Comment struct {
	ID              int       `json:"id"`
	Content         string    `json:"content"`
	PostID          int       `json:"postId"`
	UserID          int       `json:"userId"`
	Author          string    `json:"author"`
	Like            int       `json:"like"`
	Dislike         int       `json:"dislike"`
	TitlePost       string    `json:"titlePost"`
	Time            time.Time `json:"time"`
	CreatedTime     string    `json:"createdTime"`
	ToWhom          int       `json:"toWhom"`
	FromWhom        int       `json:"fromWhom"`
	CommentID       int       `json:"commentId"`
	ReplyID         int       `json:"replyId"`
	RepliesComments []Comment
	RepliesAnswer   []Comment
}

//LeaveComment for post by id
func (c *Comment) LeaveComment() error {

	commentPrepare, err := DB.Prepare(`INSERT INTO comments(content, post_id, creator_id, created_time) VALUES(?,?,?,?)`)
	if err != nil {
		log.Println(err)
	}
	defer commentPrepare.Close()
	commentExec, err := commentPrepare.Exec(c.Content, c.PostID, c.UserID, time.Now())
	if err != nil {
		log.Println(err)
		return err
	}
	//commet content
	err = DB.QueryRow("SELECT creator_id FROM posts WHERE id=?", c.PostID).Scan(&c.ToWhom)
	if err != nil {
		log.Println(err)
	}
	lid, err := commentExec.LastInsertId()
	if err != nil {
		log.Println(err)
	}
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

	_, err = DB.Exec("DELETE FROM notify  WHERE comment_id =?", id)
	if err != nil {
		log.Println(err)
	}
	_, err = DB.Exec("DELETE FROM voteState  WHERE comment_id =?", id)
	if err != nil {
		log.Println(err)
	}
	_, err = DB.Exec("DELETE FROM  comments  WHERE id =?", id)
	if err != nil {
		log.Println(err)
	}
}
