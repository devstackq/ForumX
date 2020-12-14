package config

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/devstackq/ForumX/controllers"
	"github.com/devstackq/ForumX/models"
	util "github.com/devstackq/ForumX/utils"
)

//Init Db
func Init() {
	// create DB and columns
	db, err := sql.Open("sqlite3", "forumx.db")
	if err != nil {
		log.Println(err)
	}
	db.Exec("PRAGMA foreign_keys=ON")

	postCategoryBridge, err := db.Prepare(`CREATE TABLE IF NOT EXISTS post_cat_bridge(id INTEGER PRIMARY KEY AUTOINCREMENT, post_id INTEGER, category TEXT,  FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE )`)
	comment, err := db.Prepare(`CREATE TABLE IF NOT EXISTS comments(id INTEGER PRIMARY KEY AUTOINCREMENT, content TEXT, post_id	INTEGER, creator_id	INTEGER, created_time	datetime,  count_like INTEGER DEFAULT 0, count_dislike  INTEGER DEFAULT 0, CONSTRAINT fk_key_post_comment FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE )`)
	post, err := db.Prepare(`CREATE TABLE IF NOT EXISTS posts(id INTEGER PRIMARY KEY AUTOINCREMENT, title	TEXT, content	TEXT, creator_id	INTEGER,  created_time	datetime, image	BLOB NOT NULL, count_like	INTEGER DEFAULT 0, count_dislike INTEGER DEFAULT 0, FOREIGN KEY(creator_id) REFERENCES users(id) ON DELETE CASCADE ) `)
	session, err := db.Prepare(`CREATE TABLE IF NOT EXISTS "session"("id" INTEGER PRIMARY KEY AUTOINCREMENT, "uuid"	TEXT, "user_id"	INTEGER UNIQUE )`)
	user, err := db.Prepare(`CREATE TABLE IF NOT EXISTS "users"("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "full_name"	TEXT NOT NULL, "email"	TEXT NOT NULL UNIQUE, "password" TEXT, "isAdmin" INTEGER DEFAULT 0, "age" INTEGER, "sex" TEXT, "created_time"	datetime, "city" TEXT, "image"	BLOB NOT NULL)`)
	voteState, err := db.Prepare(`CREATE TABLE IF NOT EXISTS voteState(id INTEGER PRIMARY KEY AUTOINCREMENT,  user_id INTEGER , post_id INTEGER, comment_id INTEGER,   like_state INTEGER  DEFAULT 0, dislike_state INTEGER  DEFAULT 0, unique(post_id, user_id), FOREIGN KEY(comment_id) REFERENCES comments(id), FOREIGN KEY(post_id) REFERENCES posts(id)  )`)
	notify, err := db.Prepare(`CREATE TABLE IF NOT EXISTS notify(id INTEGER PRIMARY KEY AUTOINCREMENT, post_id INTEGER,  current_user_id INTEGER, voteState INTEGER DEFAULT 0, created_time datetime, to_whom INTEGER, comment_id	INTEGER, FOREIGN KEY(comment_id) REFERENCES comments(id), FOREIGN KEY(post_id) REFERENCES posts(id) )`)
	replyComment, err := db.Prepare(`CREATE TABLE IF NOT EXISTS replyComment(id INTEGER PRIMARY KEY AUTOINCREMENT, content TEXT, post_id INTEGER, comment_id INTEGER, fromWho INTEGER, toWhom INTEGER,  created_time	datetime, FOREIGN KEY(comment_id) REFERENCES comments(id),  FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE )`)
	replyAnswer, err := db.Prepare(`CREATE TABLE IF NOT EXISTS replyAnswer (id INTEGER PRIMARY KEY AUTOINCREMENT, content TEXT, post_id INTEGER, reply_comment_id INTEGER, fromWho INTEGER, toWhom INTEGER,  created_time datetime,  FOREIGN KEY(reply_comment_id) REFERENCES replyComment(id),  FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE )`)
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
	replyComment.Exec()
	replyAnswer.Exec()

	//add connection - controllers/models & utils
	controllers.DB = db
	models.DB = db
	util.DB = db
	fmt.Println("Сукцесс коннект")
}
