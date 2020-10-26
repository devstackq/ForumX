package config

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/devstackq/ForumX/controllers"
	"github.com/devstackq/ForumX/models"
	util "github.com/devstackq/ForumX/utils"
)

func Init() {
	// create DB and columns
	db, err := sql.Open("sqlite3", "forumx2.db")
	if err != nil {
		log.Fatalln(err)
	}

	postCategoryBridge, err := db.Prepare(`CREATE TABLE IF NOT EXISTS post_cat_bridge(id INTEGER PRIMARY KEY AUTOINCREMENT, post_id INTEGER, category TEXT, FOREIGN KEY(post_id) REFERENCES posts(id) )`)
	comment, err := db.Prepare(`CREATE TABLE IF NOT EXISTS comments(id	INTEGER PRIMARY KEY AUTOINCREMENT, content TEXT, post_id	INTEGER, user_idx	INTEGER, created_time	datetime DEFAULT current_timestamp,  com_like	INTEGER DEFAULT 0, com_dislike	INTEGER DEFAULT 0, FOREIGN KEY(post_id) REFERENCES posts(id), FOREIGN KEY(user_idx) REFERENCES users(id) )`)
	like, err := db.Prepare(`CREATE   TABLE IF NOT EXISTS likes (id INTEGER PRIMARY KEY AUTOINCREMENT, 	state_id INTEGER, 	post_id	INTEGER, user_id	INTEGER,  	comment_id	INTEGER,	FOREIGN KEY(post_id) REFERENCES posts(id), 	FOREIGN KEY(user_id) REFERENCES users(id) )`)
	post, err := db.Prepare(`CREATE TABLE  IF NOT EXISTS "posts" ("id"	INTEGER PRIMARY KEY AUTOINCREMENT, "title"	TEXT, "content"	TEXT, "creator_id"	INTEGER,  "created_time"	datetime DEFAULT current_timestamp, "image"	BLOB NOT NULL, "count_like"	INTEGER DEFAULT 0, "count_dislike"	INTEGER DEFAULT 0, FOREIGN KEY("creator_id") REFERENCES "users"("id"))`)
	session, err := db.Prepare(`CREATE TABLE IF NOT EXISTS "session" ("id"	INTEGER PRIMARY KEY AUTOINCREMENT, "uuid"	TEXT, "user_id"	INTEGER UNIQUE,	FOREIGN KEY("user_id") REFERENCES  "users"("id") )`)
	user, err := db.Prepare(`CREATE TABLE IF NOT EXISTS "users" ("id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "full_name"	TEXT NOT NULL, "email"	TEXT NOT NULL UNIQUE, "password"	TEXT NOT NULL, "isAdmin"	INTEGER DEFAULT 0, "age"	INTEGER, 	"sex"	TEXT, 	"created_time"	datetime DEFAULT current_timestamp, 	"city"	TEXT,	"image"	BLOB NOT NULL	)`)

	if err != nil {
		log.Println(err)
	}

	postCategoryBridge.Exec()
	session.Exec()
	post.Exec()
	comment.Exec()
	like.Exec()
	user.Exec()

	fmt.Println("Сукцесс коннект")
	//add connection - controllers/models & utils
	controllers.DB = db
	models.DB = db
	util.DB = db

	//handlers
	http.Handle("/statics/", http.StripPrefix("/statics/", http.FileServer(http.Dir("./statics/"))))

	http.HandleFunc("/", controllers.GetAllPosts)
	http.HandleFunc("/sapid", controllers.GetAllPosts)
	http.HandleFunc("/love", controllers.GetAllPosts)
	http.HandleFunc("/science", controllers.GetAllPosts)

	http.HandleFunc("/post", controllers.GetPostByID)
	http.HandleFunc("/create/post", controllers.CreatePost)
	http.HandleFunc("/edit/post", controllers.UpdatePost)
	http.HandleFunc("/delete/post", controllers.DeletePost)

	http.HandleFunc("/comment", controllers.LeaveComment)

	http.HandleFunc("/votes", controllers.LostVotes)
	http.HandleFunc("/votes/comment", controllers.LostVotesComment)

	http.HandleFunc("/search", controllers.Search)

	http.HandleFunc("/profile", controllers.GetUserProfile)
	http.HandleFunc("/user/id/", controllers.GetAnotherProfile)
	http.HandleFunc("/edit/user", controllers.UpdateProfile)

	http.HandleFunc("/signup", controllers.Signup)
	http.HandleFunc("/signin", controllers.Signin)
	http.HandleFunc("/logout", controllers.Logout)

	// http.HandleFunc("/chat", routing.StartChat)
	log.Fatal(http.ListenAndServe(":6969", nil))
}