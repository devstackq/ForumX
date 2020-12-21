package config

import (
	"ForumX/controllers"
	"ForumX/models"
	"ForumX/utils"
	"database/sql"
	"fmt"
	"log"
)

var (
	db  *sql.DB
	err error
)

//admin@mail.kz

//Init Db
func Init() {
	// create DB and columns
	db, err = sql.Open("sqlite3", "forumx.db")
	if err != nil {
		log.Println(err)
	}
	db.Exec("PRAGMA foreign_keys=ON")

	postCategoryBridge, err := db.Prepare(`CREATE TABLE IF NOT EXISTS post_cat_bridge(id INTEGER PRIMARY KEY AUTOINCREMENT, post_id INTEGER, category_id INTEGER, FOREIGN KEY(category_id) REFERENCES category(id), FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE )`)
	if err != nil {
		log.Println(err)
	}
	comment, err := db.Prepare(`CREATE TABLE IF NOT EXISTS comments(id INTEGER PRIMARY KEY AUTOINCREMENT, content TEXT, post_id	INTEGER, creator_id	INTEGER, created_time	datetime,  count_like INTEGER DEFAULT 0, count_dislike  INTEGER DEFAULT 0, CONSTRAINT fk_key_post_comment FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE )`)
	if err != nil {
		log.Println(err)
	}
	post, err := db.Prepare(`CREATE TABLE IF NOT EXISTS posts(id INTEGER PRIMARY KEY AUTOINCREMENT, title	TEXT, content	TEXT, creator_id INTEGER,  created_time	datetime, image	BLOB NOT NULL, count_like INTEGER DEFAULT 0, count_dislike INTEGER DEFAULT 0, FOREIGN KEY(creator_id) REFERENCES users(id) ON DELETE CASCADE ) `)
	if err != nil {
		log.Println(err)
	}
	session, err := db.Prepare(`CREATE TABLE IF NOT EXISTS "session"("id" INTEGER PRIMARY KEY AUTOINCREMENT, "uuid"	TEXT, "user_id"	INTEGER UNIQUE )`)
	if err != nil {
		log.Println(err)
	}
	user, err := db.Prepare(`CREATE TABLE IF NOT EXISTS "users"("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "full_name" TEXT NOT NULL, "email"	TEXT NOT NULL UNIQUE, "username" TEXT NOT NULL UNIQUE, "password" TEXT, "isAdmin" INTEGER DEFAULT 0, "age" INTEGER, "sex" TEXT, "created_time"	datetime, "city" TEXT, "image"	BLOB NOT NULL)`)
	if err != nil {
		log.Println(err)
	}
	voteState, err := db.Prepare(`CREATE TABLE IF NOT EXISTS voteState(id INTEGER PRIMARY KEY AUTOINCREMENT,  user_id INTEGER, post_id INTEGER, comment_id INTEGER,   like_state INTEGER  DEFAULT 0, dislike_state INTEGER  DEFAULT 0, unique(post_id, user_id), FOREIGN KEY(comment_id) REFERENCES comments(id), FOREIGN KEY(post_id) REFERENCES posts(id))`)
	if err != nil {
		log.Println(err)
	}
	notify, err := db.Prepare(`CREATE TABLE IF NOT EXISTS notify(id INTEGER PRIMARY KEY AUTOINCREMENT, post_id INTEGER,  current_user_id INTEGER, voteState INTEGER DEFAULT 0, created_time datetime, to_whom INTEGER, comment_id INTEGER )`)
	if err != nil {
		log.Println(err)
	}
	replyComment, err := db.Prepare(`CREATE TABLE IF NOT EXISTS replyComment(id INTEGER PRIMARY KEY AUTOINCREMENT, content TEXT, post_id INTEGER, comment_id INTEGER, fromWhoId INTEGER, toWhomId INTEGER,  created_time datetime,  FOREIGN KEY(comment_id) REFERENCES comments(id), FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE )`)
	if err != nil {
		log.Println(err)
	}
	commentBridge, err := db.Prepare(`CREATE TABLE IF NOT EXISTS commentBridge(id INTEGER PRIMARY KEY AUTOINCREMENT, post_id INTEGER, reply_comment_id INTEGER, comment_id INTEGER, fromWhoId INTEGER, toWhoId INTEGER, created_time datetime,  FOREIGN KEY(comment_id) REFERENCES comments(id), FOREIGN KEY(reply_comment_id) REFERENCES replyComment(id), FOREIGN KEY(toWhoId) REFERENCES users(id), FOREIGN KEY(fromWhoId) REFERENCES users(id), FOREIGN KEY(post_id) REFERENCES posts(id) )`)
	if err != nil {
		log.Println(err)
	}
	category, err := db.Prepare(`CREATE TABLE IF NOT EXISTS  category(id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT UNIQUE)`)
	if err != nil {
		log.Println(err)
	}

	postCategoryBridge.Exec()
	session.Exec()
	post.Exec()
	comment.Exec()
	user.Exec()
	voteState.Exec()
	notify.Exec()
	commentBridge.Exec()
	replyComment.Exec()
	category.Exec()
	//create ategory first time
	category.Exec()
	putCategoriesInDb()

	controllers.DB = db
	models.DB = db
	utils.DB = db
	fmt.Println("Сукцесс коннект")
}

//first call -> put categories values
func putCategoriesInDb() {
	count := 0
	err = db.QueryRow("SELECT count(*) FROM category").Scan(&count)
	if err != nil {
		log.Println(err)
	}

	if count != 3 {
		categories := []string{"science", "love", "sapid"}
		for i := 0; i < 3; i++ {
			categoryPrepare, err := db.Prepare(`INSERT INTO category(name) VALUES(?)`)
			if err != nil {
				log.Println(err)
			}
			_, err = categoryPrepare.Exec(categories[i])
			if err != nil {
				log.Println(err)
			}
			defer categoryPrepare.Close()
		}
	}
}
