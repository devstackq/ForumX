package main

import (
	"ForumX/config"
	"ForumX/controllers"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	run(world)
	run(hello)
	config.Init()
	controllers.Init()

	// func main start, UserID, global Variable Вынести при запуске сервера и заполнить Uid, maybe middleware,
	//|| controllers -> send UserID -> to Model
	// google logour not work -> timeout and logour - not delete Session
	//func args, refactor,(cookie delete)

	//session delete in Db fix
	//func -> route, globalVariable -> userID

	// add valid Input data, and logger -> Middleware
	//save copokie - local - Map[string]string

	// todo another Func add	CheckMethod
	//category - table, nickname auth, url post/id=?, add 2 password form, createbutton in main page

	//logout system, when login another browser, create new Token
	//mod Name -> change github/devstackq/...

	//create post button -> вынести  В main page ...
	//author id  activity page - show
	// edit comment/.post -> add Edited message in Post
	//delete cookie & when time Expires 0, and redirect, cook life 20min

	//login current User (logged), if have session in Db, drop Session and cookie, and create new cookie and save Db & cookie

	// перегрузку методов
	// use constructor
	// use anonim func
	// use gorutine
	// try -> func use with Interface
	// architect like - Zhassymov Gt Search
	//500 status - check
	// docker check
}

//example anonim func
func hello() {
	fmt.Println("Hello")
}

func world() {
	fmt.Println("World!")
}

func run(f func()) {
	f()
}

// eaxmple reply system https://codewithawa.com/posts/creating-a-comment-and-reply-system-php-and-mysql

//comment system step 3.1
// 1 table create RepliesComment, FK(reply_id) References comments(id) -> Comment -> []ReplyComments
// form inside Client(answer comment )
// Client - form Comment, form each Comments inside comment -> ReplyForm todo

//----------------------
// comment table - comment noraml & comment under reply comment,
// reply table, uid, comment id, content, , comment id,
// insert into - 43 com -  setParentID, 12,
// client - show List comment, if have ParentId-> append Array,
//else show only COmment

//CLient -  answer -> 44com -> Form(setParentId) -> answerId : 14, parentID 44
//------------

//show/hidden by ID -> comment Field textarea
//global variable
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
