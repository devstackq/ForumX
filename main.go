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

// another browser - signin ->  logout  current browser delete session -> and signin new browser

//if not photo -> show without photo - client - Post
//redierect - user profile -> comment - fix
//domen check - org, kz ru, etc
// save photo, like - source DB refactor
//config, router refactor
//sql query - optimize, вынести в global variable
//like, dislike - refactor
//redirect logout not work,  &create post

//if cookie = 0, notify message  user, logout etc
//Представлен ли в проекте скрипт для создания образов и контейнеров? (используя скрипт для упрощения сборки)

//обработать ошикбки, log & http errors check http etc

//photo not required || set defauklt photo
//refactor function  -> Single responsibility, DRY

//design style refactor
//pagination for posts

//google acc signin -> -> back signin ? what??
//start Auth
//google token, client id, event signin Google, -> get data User,
//Name. email, photo, -> then save Db. -> authorized Forum
// Logout event, logout system, delete cookie, logout Google
//272819090705-qu6arlmkvs66hc5fuvalv6liuf2n9fj8.apps.googleusercontent.com   || W42c6sfYqhPc4O5wXMobY3av
