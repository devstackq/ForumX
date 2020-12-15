package main

import (
	"github.com/devstackq/ForumX/config"
	"github.com/devstackq/ForumX/controllers"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config.Init()
	controllers.Init()
}

1 table create RepliesComment, FK(reply_id) References comments(id) -> Comment -> []ReplyComments
form inside Client(answer comment )
Client - form Comment, form each Comments inside comment -> ReplyForm todo
// //----------------------

// replyComment, err := db.Prepare(`CREATE TABLE IF NOT EXISTS replyComment(id INTEGER PRIMARY KEY AUTOINCREMENT, content TEXT, post_id INTEGER, comment_id INTEGER, fromWho INTEGER, toWhom INTEGER,  created_time	datetime, FOREIGN KEY(comment_id) REFERENCES comments(id),  FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE )`)
// - Nested array in 1  comment
// replyComment := Comment{}

// for replyCommentQuery.Next() {

// 	err = replyCommentQuery.Scan(&replyComment.ID, &replyComment.Content, &replyComment.PostID, &replyComment.ReplyID, &replyComment.FromWhom, &replyComment.ToWhom, &replyComment.Time)
// 	if err != nil {
// 		log.Println(err.Error())
// 	}
// 	fmt.Println("/", replyComment.ID, "ReplCom ID")

// 	replyComment.CreatedTime = replyComment.Time.Format("2006 Jan _2 15:04:05")
// 	DB.QueryRow("SELECT full_name FROM users WHERE id = ?", replyComment.FromWhom).Scan(&replyComment.Author)
// 	//write answer by comment - answer answer
// 	c.RepliesComments = append(c.RepliesComments, replyComment)
// }
// comments = append(comments, c)

// create comment -> Client side, fromWho = session.UserID /13, toWho = queryDb- comment.Author /20,
// current comment_id =  CommentID(lastInsertId() / 45) -> reply_comment_id = Client(form - action/commentID 33),
// post_id = pid(form action)

// get all comment - post -> show -> comment, show COmments and answer

//-------------------

//show/hidden by ID -> comment Field textarea
//global variable
// 	DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).Scan(&s.UserID)
// var toWhom int
// DB.QueryRow("SELECT creator_id FROM comments WHERE id = ?", cid).Scan(&toWhom)

//toggle - windows under comment JS
//answer - COmments -> by userNickname -> ?

//each comment By Id-> show comments
//query - out -> models
//try todo  answer -> to another comment
// interest func - adv feat -> search, pagination

//try - event -> add sound & confetti -Login
// save photo, like - source DB refactor
//config, router refactorr
//if cookie = 0, notify message  user, logout etc
