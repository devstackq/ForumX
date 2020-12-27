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
				ID:      cid,
				Content: r.FormValue("content"),
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
func AnswerComment(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/answer/comment") {

		answer := r.FormValue("answerComment")
		currentCommentID := r.FormValue("commentID")
		pid := r.FormValue("postId")
		var toWhom int
		var lastInsertCommentID int64

		DB.QueryRow("SELECT creator_id FROM comments WHERE id = ?", currentCommentID).Scan(&toWhom)

		if utils.CheckLetter(answer) {

			comment := models.Comment{
				Content: answer,
				PostID:  pid,
				UserID:  session.UserID,
			}
			lastInsertCommentID = comment.LeaveComment()
		}
		fmt.Println(toWhom, answer, currentCommentID, "last inserted comment ID", lastInsertCommentID)

		//if have flag -> answered bool, -> show Naswer by comment
		// || comments -> table RepliesComments - child each  Comment

		replyCommentPrepare, err := DB.Prepare(`INSERT INTO commentBridge( post_id, comment_id, reply_comment_id, fromWhoId, toWhoId, create_time) VALUES(?, ?, ?, ?, ?, ?)`)
		if err != nil {
			log.Println(err)
		}
		_, err = replyCommentPrepare.Exec(pid, currentCommentID, lastInsertCommentID, session.UserID, toWhom, time.Now())
		if err != nil {
			log.Println(err)
		}
		defer replyCommentPrepare.Close()

		http.Redirect(w, r, "/post?id="+pid, 302)
	}
}
