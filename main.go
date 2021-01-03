package main

import (
	"ForumX/config"
	"ForumX/controllers"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config.Init()
	controllers.Init()

	//Comments table, parent field, 2 form, 1 form udner post, 2 form under comment
	// 1 form - parent nil, - relation only post, 2 form -> parent - commentId, relation comment
	//design -> another site, copy colors, etc
	// Heroku deploy

	//not delete rows in table- add field - visible, if Client delete post/comment-> filed visible false
	//try - create div - content editable
	//edit - user -> send  data inputs
	//optimisation, and prettier code, -> 10 princip projectio programmin - post models

	//create uniq Func -> queryDb(table, ...fields string, db)
	// todo another Func add CheckMethod
	//Start - reply system, ParentId - comment, 2 form,- 1 form not parent, second form, have parent - commentId

	//valid Input data, and , logger - add midlaweare
	//add last seen in System- when logout, save time

	//no rows in result set  -> fix

	//http: superfluous response.WriteHeader  fix

	// try errors -> with gorutine
	//func args, refactor,(cookie delete)

	//add valid Input data, and logger -> Middleware

	//not require,
	//save copokie - local - Map[string]string
	//save image -> local folder, no Db
	//mod Name -> change github/devstackq/...

	//Heroku deploy

	// перегрузку методов
	// use constructor
	// use anonim func
	// use gorutine
	// try -> func use with Interface
	// try architect like - Zhassymov Gt Search
	//500 status - check
	// docker check


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
