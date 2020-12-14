package main

import (
	"github.com/devstackq/ForumX/config"
	"github.com/devstackq/ForumX/controllers"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config.Init()
	controllers.Init()
}

//show/hidden by ID -> comment Field textarea
//redirect -> lost answer comment || lost answer answer
//global variable
// 	DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).Scan(&s.UserID)
// var toWhom int
// DB.QueryRow("SELECT creator_id FROM comments WHERE id = ?", cid).Scan(&toWhom)

//toggle - windows under comment JS
//answer - COmments -> by userNickname -> ?

//each comment By Id-> show comments
//query - out -> models
//try todo  answer -> to another comment
// interest func - adv feat -> search, pagination

//try - event -> add sound & confetti -Login
// save photo, like - source DB refactor
//config, router refactorr
//if cookie = 0, notify message  user, logout etc
