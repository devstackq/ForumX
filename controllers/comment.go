package controllers

import (
	"ForumX/models"
	"ForumX/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

//LeaveComment function
func LeaveComment(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/comment") {

		commentInput := r.FormValue("comment-text")

		if utils.CheckLetter(commentInput) {
			comment := models.Comment{
				Content: commentInput,
				PostID:  r.FormValue("curr"),
				UserID:  session.UserID,
			}
			comment.LeaveComment()
		}
		http.Redirect(w, r, "/post?id="+r.FormValue("curr"), 302)
	}
}

//UpdateComment func
func UpdateComment(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/edit/comment") {

		cid, _ := strconv.Atoi(r.FormValue("id"))

		if r.Method == "GET" {

			var comment models.Comment
			err = DB.QueryRow("SELECT * FROM comments WHERE id = ?", cid).Scan(&comment.ID, &comment.Content, &comment.PostID, &comment.UserID, &comment.CreatedTime, &comment.UpdatedTime, &comment.Like, &comment.Dislike)
			if err != nil {
				fmt.Println(err)
			}

			utils.RenderTemplate(w, "header", utils.IsAuth(r))
			utils.RenderTemplate(w, "update_comment", comment)
		}
		if r.Method == "POST" {

			comment := models.Comment{
				ID:          cid,
				Content:     r.FormValue("content"),
				UpdatedTime: time.Now(),
			}

			comment.UpdateComment()
		}
		http.Redirect(w, r, "/profile", 302)
	}
}

//DeleteComment dsa
func DeleteComment(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/delete/comment") {
		models.DeleteComment(r.FormValue("id"))
	}
	http.Redirect(w, r, "/profile", 302)
}

//AnswerComment func replyComment
func ReplyComment(w http.ResponseWriter, r *http.Request) {

	//post 8, 5 commentId, Reply 1 id, from Uid13, To 24, currentReplyId, answerReplyId
	if utils.URLChecker(w, r, "/reply/comment/replyId") {

		content := r.FormValue("answerComment")
		currentReplyID := r.FormValue("replyId")

		cID, _ := strconv.Atoi(r.FormValue("commentId"))

		pid := r.FormValue("postId")
		var toWhom int
		//var lastInsertCommentID int64

		//if answer comment -> show UserID,  commentCreate, else replies FromWho
		//DB.QueryRow("SELECT creator_id FROM comments WHERE id = ?", currentCommentID).Scan(&toWhom)
		DB.QueryRow("SELECT fromWho FROM replies WHERE id = ?", currentReplyID).Scan(&toWhom)
		//else get comment table, user_id, if 1 answer - in comment
		//if no reply, create First reply -> init reply this comment
		if utils.CheckLetter(content) {

			comment := models.Comment{
				CommentID: cID,
				Content:   content,
				PostID:    pid,
				FromWhom:  session.UserID,
				ToWhom:    toWhom,
			}
			fmt.Print(comment, "reply commen ", currentReplyID, toWhom)
			//lastInsertCommentID = comment.LeaveComment()
			//user_id INTEGER, content TEXT, post_id INTEGER, comment_id INTEGER, fromWhoId INTEGER, toWhomId INTEGER,  created_time datetime
			commentPrepare, err := DB.Prepare(`INSERT INTO replies(content, post_id, comment_id, fromWho, toWho, created_time) VALUES(?,?,?,?,?,?)`)
			if err != nil {
				log.Println(err)
			}
			commentExec, err := commentPrepare.Exec(comment.Content, comment.PostID, comment.CommentID, comment.FromWhom, comment.ToWhom, time.Now())
			if err != nil {
				log.Println(err)
			}
			commentExec.LastInsertId()
		}
		// 1,3,1,1

		//fmt.Println(toWhom, content, cID, "last inserted comment ID", lastInsertCommentID)

		//if have flag -> answered bool, -> show Naswer by comment
		// || comments -> table RepliesComments - child each  Comment

		// replyCommentPrepare, err := DB.Prepare(`INSERT INTO commentBridge( post_id, comment_id, reply_comment_id, fromWhoId, toWhoId, create_time) VALUES(?, ?, ?, ?, ?, ?)`)
		// if err != nil {
		// 	log.Println(err)
		// }
		// _, err = replyCommentPrepare.Exec(pid, cID, lastInsertCommentID, session.UserID, toWhom, time.Now())
		// if err != nil {
		// 	log.Println(err)
		// }
		// defer replyCommentPrepare.Close()

		http.Redirect(w, r, "/post?id="+pid, 302)
	}
}
