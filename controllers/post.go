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

//receive request, from client, query params, category ID, then query DB, depends catID, get Post this catID
func GetAllPosts(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" && r.URL.Path != "/science" && r.URL.Path != "/love" && r.URL.Path != "/sapid" {
		util.DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}

	fv := models.Filter{
		Like:     r.FormValue("likes"),
		Date:     r.FormValue("date"),
		Category: r.FormValue("cats"),
	}

	posts, endpoint, category, err := fv.GetAllPost(r)

	if err != nil {
		log.Fatal(err)
	}

	util.DisplayTemplate(w, "header", util.IsAuth(r))

	if endpoint == "/" {
		util.DisplayTemplate(w, "index", posts)
	} else {
		//send category
		msg := []byte(fmt.Sprintf("<h2 id='category'> `Category: %s` </h2>", category))
		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
		util.DisplayTemplate(w, "catTemp", posts)
	}
}

//get 1 post by id
func GetPostById(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/post" {
		util.DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}
	id, _ := strconv.Atoi(r.FormValue("id"))
	p := models.Posts{ID: id}
	comments, post, err := p.GetPostById(r)

	if err != nil {
		log.Println(err)
	}
	util.DisplayTemplate(w, "header", util.IsAuth(r))
	util.DisplayTemplate(w, "posts", post)
	util.DisplayTemplate(w, "comment", comments)
}

//create post
func CreatePost(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/create/post" {
		util.DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}

	msg = ""

	switch r.Method {
	case "GET":
		util.DisplayTemplate(w, "header", util.IsAuth(r))
		util.DisplayTemplate(w, "create", &msg)
	case "POST":
		access, session := util.CheckForCookies(w, r)
		log.Println(access, "access status")
		if !access {
			http.Redirect(w, r, "/signin", 302)
			return
		}
		r.ParseMultipartForm(10 << 20)
		f, _, _ := r.FormFile("uploadfile")
		f2, _, _ := r.FormFile("uploadfile")
		categories, _ := r.Form["input"]

		post := models.Posts{
			Title:      r.FormValue("title"),
			Content:    r.FormValue("content"),
			Categories: categories,
			FileS:      f,
			FileI:      f2,
			Session:    session,
		}
		post.CreatePost(w, r)
		http.Redirect(w, r, "/", 200)
	}
	//	http.Redirect(w, r, "/", http.StatusOK)
}

//update post
func UpdatePost(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		pid, _ := strconv.Atoi(r.URL.Query().Get("id"))
		p := models.Posts{}
		p.PostIDEdit = pid
		util.DisplayTemplate(w, "updatepost", p)
	}

	if r.Method == "POST" {

		access, _ := util.CheckForCookies(w, r)
		if !access {
			http.Redirect(w, r, "/signin", 302)
			return
		}
		imgBytes := util.FileByte(r)
		pid, _ := strconv.Atoi(r.FormValue("pid"))

		p := models.Posts{
			Title:   r.FormValue("title"),
			Content: r.FormValue("content"),
			Image:   imgBytes,
			ID:      pid,
		}

		err = p.UpdatePost()

		if err != nil {
			defer log.Println(err, "upd post err")
		}
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

//delete post
func DeletePost(w http.ResponseWriter, r *http.Request) {

	pid, _ := strconv.Atoi(r.URL.Query().Get("id"))
	p := models.Posts{ID: pid}

	access, _ := util.CheckForCookies(w, r)
	if !access {
		http.Redirect(w, r, "/signin", 302)
		return
	}

	err = p.DeletePost()

	if err != nil {
		panic(err.Error())
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

//search
func Search(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/search" {
		util.DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}

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
