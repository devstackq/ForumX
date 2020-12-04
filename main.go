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

// add interest func - adv feat -> search, pagination
//fix create user -> signup  after -> github signout?- password null
//fix design
delete notify - when delete comment & post
// create post, if category, empty - set -> default category - sapid

//try - event -> add sound & confetti -Login
// save photo, like - source DB refactor
//config, router refactorr
//if cookie = 0, notify message  user, logout etc

// 1 request, 910 additional, 0904 - 101202 ->
// 2 request -7575
// 3 request 910 additional, 090410 - 101202 ->Otegen batyr etc
