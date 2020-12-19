package main

import (
	"ForumX/config"
	"ForumX/controllers"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config.Init()
	controllers.Init()
	time.Sleep( 10 * time.Second)
	fmt.Println("timer 10 sec")
	
	try -> func use with Interface
	difference ? default func  utils.UpdateVoteNotify
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
