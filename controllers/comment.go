package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/devstackq/ForumX/models"
	util "github.com/devstackq/ForumX/utils"
)

//LeaveComment function
func LeaveComment(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/comment") {

		if r.Method == "POST" {

			pid, _ := strconv.Atoi(r.FormValue("curr"))
			commentInput := r.FormValue("comment-text")

			access, s := util.IsCookie(w, r)
			if !access {
				return
			}

			DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).Scan(&s.UserID)

			if util.CheckLetter(commentInput) {

				comment := models.Comment{
					Content: commentInput,
					PostID:  pid,
					UserID:  s.UserID,
				}

				err = comment.LeaveComment()

				if err != nil {
					log.Println(err.Error())
				}
			}
		}
		http.Redirect(w, r, "/post?id="+r.FormValue("curr"), 302)
	}
}

//UpdateComment func
func UpdateComment(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/edit/comment") {

		access, _ := util.IsCookie(w, r)
		if !access {
			http.Redirect(w, r, "/signin", 200)
			return
		}
		cid, _ := strconv.Atoi(r.FormValue("id"))

		if r.Method == "GET" {

			var comment models.Comment
			err = DB.QueryRow("SELECT * FROM comments WHERE id = ?", cid).Scan(&comment.ID, &comment.Content, &comment.PostID, &comment.UserID, &comment.CreatedTime, &comment.Like, &comment.Dislike)
			if err != nil {
				fmt.Println(err)
			}
			util.DisplayTemplate(w, "header", util.IsAuth(r))
			fmt.Println(comment, "FF")
			util.DisplayTemplate(w, "update_comment", comment)
		}
		if r.Method == "POST" {

			comment := models.Comment{
				ID:      cid,
				Content: r.FormValue("content"),
			}

			comment.UpdateComment()
		}
		http.Redirect(w, r, "/profile", 302)
	}
}

//DeleteComment dsa
func DeleteComment(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/delete/comment") {

		access, _ := util.IsCookie(w, r)
		if !access {
			http.Redirect(w, r, "/signin", 200)
			return
		}
		models.DeleteComment(r.FormValue("id"))
	}
	http.Redirect(w, r, "/profile", 302)
}

func ReplyAnswer(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/reply/answer/comment") {

		access, s := util.IsCookie(w, r)
		if !access {
			http.Redirect(w, r, "/signin", 200)
			return
		}

		content := r.FormValue("answer")
		asnwerID := r.FormValue("answerID")
		pid := r.FormValue("pid")
		cid := r.FormValue("cid")

		DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).Scan(&s.UserID)
		var toWhom int
		DB.QueryRow("SELECT fromWho FROM replyComment WHERE id = ?", cid).Scan(&toWhom)
		//	DB.QueryRow("select ")
		fmt.Println(content, asnwerID, pid, cid, "LKS", toWhom)

		replyAnswerPrepare, err := DB.Prepare(`INSERT INTO replyAnswer(content, post_id, reply_comment_id, fromWho, toWhom, created_time) VALUES(?, ?, ?, ?, ?, ?)`)
		if err != nil {
			log.Println(err)
		}
		defer replyAnswerPrepare.Close()
		_, err = replyAnswerPrepare.Exec(content, pid, asnwerID, s.UserID, toWhom, time.Now())
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/post?id="+pid, 302)
	}
}

//AnswerComment func
func AnswerComment(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/answer/comment") {

		access, s := util.IsCookie(w, r)
		if !access {
			http.Redirect(w, r, "/signin", 200)
			return
		}

		answer := r.FormValue("answerComment")
		cid := r.FormValue("commentID")
		pid := r.FormValue("postId")

		DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).Scan(&s.UserID)
		var toWhom int
		DB.QueryRow("SELECT creator_id FROM comments WHERE id = ?", cid).Scan(&toWhom)

		replyCommentPrepare, err := DB.Prepare(`INSERT INTO  replyComment(content, post_id, comment_id, fromWho, toWhom, created_time) VALUES(?, ?, ?, ?, ?, ?)`)
		if err != nil {
			log.Println(err)
		}
		defer replyCommentPrepare.Close()
		_, err = replyCommentPrepare.Exec(answer, pid, cid, s.UserID, toWhom, time.Now())
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/post?id="+pid, 302)
	}
}
