package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/devstackq/ForumX/models"
	util "github.com/devstackq/ForumX/utils"
)

//create comment
func CreateComment(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/comment" {
		util.DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}

	if r.Method == "POST" {

		access, s := util.CheckForCookies(w, r)
		if !access {
			http.Redirect(w, r, "/signin", 302)
			return
		}

		DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).Scan(&s.UserID)

		pid, _ := strconv.Atoi(r.FormValue("curr"))
		comment := r.FormValue("comment-text")

		if util.CheckLetter(comment) {

			com := models.Comment{
				Content: comment,
				PostID:  pid,
				UserID:  s.UserID,
			}

			err = com.LeaveComment()

			if err != nil {
				log.Println(err.Error())
			}
		}
	}
	http.Redirect(w, r, "post?id="+r.FormValue("curr"), 301)
}
