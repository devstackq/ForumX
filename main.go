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

	//Not unique username - msg not correct -> if 2 time signup
	//google auth, github auth -> signup -Chekck  if exist user -> show message
	//by nickename, google auth, github -> logout other browser todo,

	//try another logic -> when user login, but cookie another -> logout, with Table field - resseion(true), -> Note -phone see

	//optimisation - beatu code, -> fix, signin server auth utils MEthods refactor
	//no rows in result set  -> fix

	//IsCookieExpiration -> refactor - logour function
	//add writeHeader()
	//redirect - signin not work,

	// try errors -> with gorutine
	//|| controllers -> send UserID -> to Model
	//func args, refactor,(cookie delete)

	//add Middleware, logger & checkData

	//add valid Input data, and logger -> Middleware
	//create post button -> вынести  В main page ...
	// edit comment/.post -> add Edited message in Post


	//not require, 
	//save copokie - local - Map[string]string
	//save image -> local folder, no Db
	// todo another Func add	CheckMethod
	//mod Name -> change github/devstackq/...
	//Start - reply system
	
	// перегрузку методов
	// use constructor
	// use anonim func
	// use gorutine
	// try -> func use with Interface
	// try architect like - Zhassymov Gt Search
	//500 status - check
	// docker check
	
	// done: category - table, nickname auth, url post/id=?, add 2 password form, createbutton in main page, show user name activity,  optimisation code, signin system, -> another session delete
	//,cookie expired -> logout

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
