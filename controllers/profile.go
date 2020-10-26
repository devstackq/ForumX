package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/devstackq/ForumX/models"
	util "github.com/devstackq/ForumX/utils"
)

//GetUserProfile  current -> user page
func GetUserProfile(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/profile") {

		if r.Method == "GET" {
			access, _ := util.CheckForCookies(w, r)
			if !access {
				http.Redirect(w, r, "/signin", 302)
				return
			}
			cookie, _ := r.Cookie("_cookie")
			//if userId now, createdPost uid equal -> show
			likedpost, posts, comments, user, err := models.GetUserProfile(r, w, cookie)
			if err != nil {
				panic(err)
			}

			util.DisplayTemplate(w, "header", util.IsAuth(r))
			util.DisplayTemplate(w, "profile", user)
			util.DisplayTemplate(w, "likedpost", likedpost)
			util.DisplayTemplate(w, "postuser", posts)
			util.DisplayTemplate(w, "commentuser", comments)

			//delete coookie db
			go func() {
				for now := range time.Tick(299 * time.Minute) {
					util.СheckCookieLife(now, cookie, w, r)
					//next logout each 300 min
					time.Sleep(299 * time.Minute)
				}
			}()
		}
	}
}

//GetAnotherProfile  other user page
func GetAnotherProfile(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/user/id") {

		if r.Method == "POST" {

			uid := models.User{Temp: r.FormValue("uid")}
			posts, user, err := uid.GetAnotherProfile(r)
			if err != nil {
				panic(err)
			}
			util.DisplayTemplate(w, "header", util.IsAuth(r))
			util.DisplayTemplate(w, "user", user)
			util.DisplayTemplate(w, "postuser", posts)
		}
	}
}

//update profile
func UpdateProfile(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/edit/user") {

		if r.Method == "GET" {
			util.DisplayTemplate(w, "header", util.IsAuth(r))
			util.DisplayTemplate(w, "updateuser", "")
		}

		if r.Method == "POST" {

			access, s := util.CheckForCookies(w, r)
			if !access {
				http.Redirect(w, r, "/signin", 302)
				return
			}

			DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).
				Scan(&s.UserID)

			is, _ := strconv.Atoi(r.FormValue("age"))

			p := models.User{
				FullName: r.FormValue("fullname"),
				Age:      is,
				Sex:      r.FormValue("sex"),
				City:     r.FormValue("city"),
				Image:    util.FileByte(r),
				ID:       s.UserID,
			}

			err = p.UpdateProfile()

			if err != nil {
				panic(err.Error())
			}
		}
		http.Redirect(w, r, "/profile", http.StatusFound)
	}
}