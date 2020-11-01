package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	structure "github.com/devstackq/ForumX/general"
	"github.com/devstackq/ForumX/models"
	util "github.com/devstackq/ForumX/utils"
)

var (
	err  error
	DB   *sql.DB
	msg  = structure.API.Message
	auth = structure.API.Authenticated
)

//GetAllPosts  by category || all posts
func GetAllPosts(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" && r.URL.Path != "/science" && r.URL.Path != "/love" && r.URL.Path != "/sapid" {
		util.DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}

	filterValue := models.Filter{
		Like:     r.FormValue("likes"),
		Date:     r.FormValue("date"),
		Category: r.FormValue("cats"),
	}

	posts, endpoint, category, err := filterValue.GetAllPost(r)

	if err != nil {
		log.Fatal(err)
	}

	util.DisplayTemplate(w, "header", util.IsAuth(r))

	if endpoint == "/" {
		util.DisplayTemplate(w, "index", posts)
	} else {
		//send category value
		msg := []byte(fmt.Sprintf("<h2 id='category'> `Category: %s` </h2>", category))
		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
		util.DisplayTemplate(w, "category_post_template", posts)
	}
}

//GetPostByID  1 post by id
func GetPostByID(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/post") {

		id, _ := strconv.Atoi(r.FormValue("id"))
		pid := models.Post{ID: id}
		comments, post, err := pid.GetPostByID(r)

		if err != nil {
			log.Println(err)
		}
		util.DisplayTemplate(w, "header", util.IsAuth(r))
		util.DisplayTemplate(w, "posts", post)
		util.DisplayTemplate(w, "comment_post", comments)
	}
}

//CreatePost  function
func CreatePost(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/create/post") {

		switch r.Method {
		case "GET":
			util.DisplayTemplate(w, "header", util.IsAuth(r))
			util.DisplayTemplate(w, "create_post", &msg)
		case "POST":
			access, session := util.IsCookie(w, r)
			log.Println(access, "access status")
			if !access {
				http.Redirect(w, r, "/signin", 302)
				return
			}
			r.ParseMultipartForm(10 << 20)
			f, _, _ := r.FormFile("uploadfile")
			f2, _, _ := r.FormFile("uploadfile")

			//if file == nil, no set photo -> client delete img tag

			// if f != nil && f2 != nil {
			// 	IsPhoto = true
			// }
			categories, _ := r.Form["input"]

			post := models.Post{
				Title:      r.FormValue("title"),
				Content:    r.FormValue("content"),
				Categories: categories,
				FileS:      f,
				FileI:      f2,
				Session:    session,
				IsPhoto:    true,
			}
			post.CreatePost(w, r)
			http.Redirect(w, r, "/", 200)
		}
	}
}

//UpdatePost function
func UpdatePost(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/edit/post") {

		if r.Method == "GET" {
			pid, _ := strconv.Atoi(r.URL.Query().Get("id"))
			p := models.Post{PostIDEdit: pid}
			util.DisplayTemplate(w, "update_post", p)
		}
		if r.Method == "POST" {
			access, _ := util.IsCookie(w, r)
			if !access {
				http.Redirect(w, r, "/signin", 302)
				return
			}
			imgBytes := util.FileByte(r, "post")
			pid, _ := strconv.Atoi(r.FormValue("pid"))

			p := models.Post{
				Title:   r.FormValue("title"),
				Content: r.FormValue("content"),
				Image:   imgBytes,
				ID:      pid,
			}

			err = p.UpdatePost()

			if err != nil {
				//try hadnler all error
				defer log.Println(err, "upd post err")
			}
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

//DeletePost function
func DeletePost(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/delete/post") {

		access, _ := util.IsCookie(w, r)
		if !access {
			http.Redirect(w, r, "/signin", 302)
			return
		}
		pid, _ := strconv.Atoi(r.URL.Query().Get("id"))
		p := models.Post{ID: pid}

		err = p.DeletePost()

		if err != nil {
			panic(err.Error())
		}
		http.Redirect(w, r, "/", 200)
	}
}

//search
func Search(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/search") {

		if r.Method == "GET" {
			util.DisplayTemplate(w, "search", http.StatusFound)
		}

		if r.Method == "POST" {

			foundPosts, err := models.Search(w, r)

			if err != nil {
				panic(err)
			}
			util.DisplayTemplate(w, "header", util.IsAuth(r))
			util.DisplayTemplate(w, "index", foundPosts)
		}
	}
}
