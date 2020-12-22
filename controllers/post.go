package controllers

import (
	"ForumX/general"
	"ForumX/models"
	"ForumX/utils"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var (
	err  error
	DB   *sql.DB
	msg  = general.API.Message
	auth = general.API.Authenticated
)

//GetAllPosts  by category || all posts
func GetAllPosts(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" && r.URL.Path != "/science" && r.URL.Path != "/love" && r.URL.Path != "/sapid" {
		utils.DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}

	filterValue := models.Filter{
		Like:     r.FormValue("likes"),
		Date:     r.FormValue("date"),
		Category: r.FormValue("cats"),
	}

	p, endpoint, category, err := filterValue.GetAllPost(r, r.FormValue("next"), r.FormValue("prev"))

	if err != nil {
		log.Fatal(err)
	}
	//s := utils.IsAuth(r)
	utils.DisplayTemplate(w, "header", utils.IsAuth(r))

	if endpoint == "/" {
		posts := models.AllPosts{Posts: p}
		if session.UserID > 0 {
			posts.Auth = "success"
		}
		utils.DisplayTemplate(w, "index", posts)
	} else {
		//send category value
		msg := []byte(fmt.Sprintf("<h3 id='category'> %s </h3>", category))
		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
		utils.DisplayTemplate(w, "category_post_template", p)
	}
}
//GetPostByID  1 post by id
func GetPostByID(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/post") {

		count := 0
		err = DB.QueryRow("SELECT count(*) FROM posts").Scan(&count)
		if err != nil {
			log.Println(err)
		}

		id, _ := strconv.Atoi(r.FormValue("id"))

		if id > 0 && id <= count {
			pid := models.Post{ID: id}
			comments, post, err := pid.GetPostByID(r)
			if err != nil {
				log.Println(err)
			}
			utils.DisplayTemplate(w, "header", utils.IsAuth(r))
			utils.DisplayTemplate(w, "posts", post)
			utils.DisplayTemplate(w, "comment_post", comments)
			//utils.DisplayTemplate(w, "reply_comment", repliesComment)
		} else {
			utils.DisplayTemplate(w, "404page", http.StatusBadRequest)
		}
	}
}

//CreatePost  function
func CreatePost(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/create/post") {

		//switch r.Method {
		if r.Method == "GET" {
			utils.DisplayTemplate(w, "header", utils.IsAuth(r))
			utils.DisplayTemplate(w, "create_post", &msg)
		}
		if r.Method == "POST" {
			//r.ParseMultipartForm(10 << 20)
			f, _, _ := r.FormFile("uploadfile")
			f2, _, _ := r.FormFile("uploadfile")

			categories, _ := r.Form["input"]

			photoFlag := false
			if f != nil && f2 != nil {
				photoFlag = true
			}
			post := models.Post{
				Title:      r.FormValue("title"),
				Content:    r.FormValue("content"),
				Categories: categories,
				FileS:      f,
				FileI:      f2,
				Session:    session,
				IsPhoto:    photoFlag,
			}
			post.CreatePost(w, r)
		}
	}
}

//UpdatePost function
func UpdatePost(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/edit/post") {

		pid, _ := strconv.Atoi(r.FormValue("id"))

		if r.Method == "GET" {

			var p models.Post
			DB.QueryRow("SELECT * FROM posts WHERE id = ?", pid).Scan(&p.ID, &p.Title, &p.Content, &p.CreatorID, &p.CreatedTime, &p.Image, &p.Like, &p.Dislike)
			p.ImageHTML = base64.StdEncoding.EncodeToString(p.Image)

			utils.DisplayTemplate(w, "header", utils.IsAuth(r))
			utils.DisplayTemplate(w, "update_post", p)
		}

		if r.Method == "POST" {

			p := models.Post{
				Title:   r.FormValue("title"),
				Content: r.FormValue("content"),
				Image:   utils.IsImage(r),
				ID:      pid,
			}

			err = p.UpdatePost()

			if err != nil {
				//try hadnler all error
				defer log.Println(err, "upd post err")
			}
		}
		http.Redirect(w, r, "/profile", 302)
		//http.Redirect(w, r, "/post?id="+strconv.Itoa(int(pid)), 302)
	}
}

//DeletePost function
func DeletePost(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/delete/post") {

		pid, _ := strconv.Atoi(r.URL.Query().Get("id"))
		p := models.Post{ID: pid}

		err = p.DeletePost()

		if err != nil {
			log.Println(err.Error())
		}
		http.Redirect(w, r, "/profile", 302)
	}
}

//Search
func Search(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/search") {

		if r.Method == "GET" {
			utils.DisplayTemplate(w, "search", http.StatusFound)
		}

		if r.Method == "POST" {

			foundPosts, err := models.Search(w, r)

			if err != nil {
				log.Println(err)
			}
			if foundPosts == nil {
				utils.DisplayTemplate(w, "header", utils.IsAuth(r))
				msg := []byte(fmt.Sprintf("<h2 id='notFound'> Nihuya ne naideno </h2>"))
				w.Header().Set("Content-Type", "application/json")
				w.Write(msg)
				utils.DisplayTemplate(w, "index", nil)
			} else {
				utils.DisplayTemplate(w, "header", utils.IsAuth(r))
				utils.DisplayTemplate(w, "index", foundPosts)
			}
		}
	}
}
