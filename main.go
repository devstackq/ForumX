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

//delete notify - when delete comment & post

// add interest func - adv feat -> search, pagination
//fix design
// create post, if category, empty - set -> default category - sapid

//try - event -> add sound & confetti -Login
// save photo, like - source DB refactor
//config, router refactorr
//if cookie = 0, notify message  user, logout etc
