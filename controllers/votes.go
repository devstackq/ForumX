package controllers

import (
	"ForumX/models"
	"ForumX/utils"
	"net/http"
)

//VotesPost func Post
func VotesPost(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/votes/post") {

		pid := r.URL.Query().Get("id")
		lukas := r.FormValue("like")
		dislike := r.FormValue("dislike")

		if r.Method == "POST" {

			if lukas == "1" {
				models.VoteLike(w, r, pid, "post", session)
			}
			if dislike == "1" {
				models.VoteDislike(w, r, pid, "post", session)
			}
		}
		http.Redirect(w, r, "/post?id="+pid, 302)
	}
}

//VotesComment function
func VotesComment(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/votes/comment") {

		commentID := r.URL.Query().Get("commentID")
		commentDis := r.FormValue("commentDislike")
		commentLike := r.FormValue("commentLike")

		if r.Method == "POST" {
			if commentLike == "1" {
				models.VoteLike(w, r, commentID, "comment", session)
			}
			if commentDis == "1" {
				models.VoteDislike(w, r, commentID, "comment", session)
			}
			http.Redirect(w, r, "/post?id="+r.FormValue("pidc"), 302)
		}
	}
}
