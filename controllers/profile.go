package controllers

import (
	"ForumX/models"
	"ForumX/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

//GetUserProfile  current -> user page
func GetUserProfile(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/profile") {

		if r.Method == "GET" {
			//if userId now, createdPost uid equal -> show
			dislikedPost, likedPost, posts, comments, user, err := models.GetUserProfile(r, w, CookieBrowser)
			if err != nil {
				log.Println(err)
			}
			//check if current cookie equal - cookie
			utils.DisplayTemplate(w, "header", utils.IsAuth(r))
			utils.DisplayTemplate(w, "profile", user)
			utils.DisplayTemplate(w, "created_post", posts)
			utils.DisplayTemplate(w, "favorited_post", likedPost)
			utils.DisplayTemplate(w, "disliked_post", dislikedPost)
			utils.DisplayTemplate(w, "comment_user", comments)

			//delete coookie db, 20 min
			go func() {
				for range time.Tick(19 * time.Minute) {
					utils.IsCookieExpiration(w, r)
					fmt.Println("del cookie in DB")
					//time.Sleep(1 * time.Minute)
					break
				}
			}()
		}
	}
}

//GetAnotherProfile  other user page
func GetAnotherProfile(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/user/id") {

		if r.Method == "POST" {

			uid := models.User{Temp: r.FormValue("uid")}
			posts, user, err := uid.GetAnotherProfile(r)
			if err != nil {
				log.Println(err)
			}

			utils.DisplayTemplate(w, "header", utils.IsAuth(r))
			utils.DisplayTemplate(w, "another_user", user)
			utils.DisplayTemplate(w, "created_post", posts)
		}
	}
}

//GetUserActivities func
func GetUserActivities(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/activity") {

		notifies := models.GetUserActivities(w, r)

		if err != nil {
			log.Println(err)
		}
		utils.DisplayTemplate(w, "header", utils.IsAuth(r))

		if notifies != nil {
			utils.DisplayTemplate(w, "activity", notifies)
		}
	}
}

//UpdateProfile function
func UpdateProfile(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/edit/user") {

		if r.Method == "GET" {
			utils.DisplayTemplate(w, "header", utils.IsAuth(r))
			utils.DisplayTemplate(w, "profile_update", "")
		}

		if r.Method == "POST" {

			is, _ := strconv.Atoi(r.FormValue("age"))

			p := models.User{
				FullName: r.FormValue("fullname"),
				Age:      is,
				Sex:      r.FormValue("sex"),
				City:     r.FormValue("city"),
				Image:    utils.FileByte(r, "user"),
				ID:       session.UserID,
			}

			err = p.UpdateProfile()

			if err != nil {
				log.Println(err.Error())
			}
		}
		http.Redirect(w, r, "/profile", http.StatusFound)
	}
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/delete/account") {

		if r.Method == "POST" {

			var p models.User

			err := json.NewDecoder(r.Body).Decode(&p.ID)
			if err != nil {
				log.Println(err)
			}

			p.DeleteAccount(w, r)
			fmt.Println("delete account by ID", p.ID)
		}
		http.Redirect(w, r, "/", 302)
	}
}
