package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/devstackq/ForumX/models"
	"github.com/devstackq/ForumX/routing"
	util "github.com/devstackq/ForumX/utils"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	//init call func router, and dbconfig file
	CreateDB()

	http.Handle("/statics/", http.StripPrefix("/statics/", http.FileServer(http.Dir("./statics/"))))
	http.HandleFunc("/", routing.GetAllPosts)
	http.HandleFunc("/sapid", routing.GetAllPosts)
	http.HandleFunc("/love", routing.GetAllPosts)
	http.HandleFunc("/science", routing.GetAllPosts)

	http.HandleFunc("/post", routing.GetPostById)
	http.HandleFunc("/profile", routing.GetUserProfile)
	http.HandleFunc("/user/id/", routing.GetAnotherProfile)

	http.HandleFunc("/comment", routing.CreateComment)
	http.HandleFunc("/create/post", routing.CreatePost)

	http.HandleFunc("/edit/post", routing.UpdatePost)
	http.HandleFunc("/delete/post", routing.DeletePost)
	http.HandleFunc("/edit/user", routing.UpdateProfile)

	http.HandleFunc("/signup", routing.Signup)
	http.HandleFunc("/signin", routing.Signin)
	http.HandleFunc("/logout", routing.Logout)

	http.HandleFunc("/votes", routing.LostVotes)
	http.HandleFunc("/votes/comment", routing.LostVotesComment)
	http.HandleFunc("/search", routing.Search)
	// http.HandleFunc("/chat", routing.StartChat)
	log.Fatal(http.ListenAndServe(":6969", nil))
}

//connect Db, create table if not exist
func CreateDB() {

	db, err := sql.Open("sqlite3", "forumx2.db")
	if err != nil {
		log.Fatalln(err)
		//	panic(err)
	}

	postCategoryBridge, err := db.Prepare(`CREATE TABLE IF NOT EXISTS post_cat_bridge(id INTEGER PRIMARY KEY AUTOINCREMENT, post_id INTEGER, category TEXT, FOREIGN KEY(post_id) REFERENCES posts(id) )`)
	comments, err := db.Prepare(`CREATE TABLE IF NOT EXISTS comments(id	INTEGER PRIMARY KEY AUTOINCREMENT, content TEXT, post_id	INTEGER, user_idx	INTEGER, created_time	datetime DEFAULT current_timestamp,  com_like	INTEGER DEFAULT 0, com_dislike	INTEGER DEFAULT 0, FOREIGN KEY(post_id) REFERENCES posts(id), FOREIGN KEY(user_idx) REFERENCES users(id) )`)
	likes, err := db.Prepare(`CREATE   TABLE IF NOT EXISTS likes (id INTEGER PRIMARY KEY AUTOINCREMENT, 	state_id INTEGER, 	post_id	INTEGER, user_id	INTEGER,  	comment_id	INTEGER,	FOREIGN KEY(post_id) REFERENCES posts(id), 	FOREIGN KEY(user_id) REFERENCES users(id) )`)
	posts, err := db.Prepare(`CREATE TABLE  IF NOT EXISTS "posts" ("id"	INTEGER PRIMARY KEY AUTOINCREMENT, "title"	TEXT, "content"	TEXT, "creator_id"	INTEGER,  "created_time"	datetime DEFAULT current_timestamp, "image"	BLOB NOT NULL, "count_like"	INTEGER DEFAULT 0, "count_dislike"	INTEGER DEFAULT 0, FOREIGN KEY("creator_id") REFERENCES "users"("id"))`)
	session, err := db.Prepare(`CREATE TABLE IF NOT EXISTS "session" ("id"	INTEGER PRIMARY KEY AUTOINCREMENT, "uuid"	TEXT, "user_id"	INTEGER UNIQUE,	FOREIGN KEY("user_id") REFERENCES  "users"("id") )`)
	users, err := db.Prepare(`CREATE TABLE IF NOT EXISTS "users" ("id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "full_name"	TEXT NOT NULL, "email"	TEXT NOT NULL UNIQUE, "password"	TEXT NOT NULL, "isAdmin"	INTEGER DEFAULT 0, "age"	INTEGER, 	"sex"	TEXT, 	"created_time"	datetime DEFAULT current_timestamp, 	"city"	TEXT,	"image"	BLOB NOT NULL	)`)

	if err != nil {
		log.Fatal(err)
	}

	postCategoryBridge.Exec()
	session.Exec()
	posts.Exec()
	comments.Exec()
	likes.Exec()
	users.Exec()

	fmt.Println("Сукцесс конект")
	routing.DB = db
	models.DB = db
	util.DB = db
}

//Обнаруживает ли проект неправильный адрес электронной почты или пароль? - password min 8 symbols, and lowerBig number pwd
//active session -> show other browser || show notification
//if cookie = 0, notify message  user, logout etc
//Представлен ли в проекте скрипт для создания образов и контейнеров? (используя скрипт для упрощения сборки)
//like, dislike - refactor
//обработать ошикбки, log
//errors check http etc

//name func, variable - norm
//slice - controller/model/view
//conf file
//check contreoller -> middleware -> check data from CLient todo

//like, create post -slow work ?
//redirect logout not work,  &create post

//checkCookieLife(now, w, r, cookie) - del cookie client and backend - redirect main page

//photo not required || set defauklt photo
//refactor function  -> Single responsibility, DRY

//design refactor
// pagination for posts
//google acc signin -> -> back signin ? what??

//start Auth
//google token, client id, event signin Google, -> get data User,
//Name. email, photo, -> then save Db. -> authorized Forum
// Logout event, logout system, delete cookie, logout Google
//272819090705-qu6arlmkvs66hc5fuvalv6liuf2n9fj8.apps.googleusercontent.com   || W42c6sfYqhPc4O5wXMobY3av
