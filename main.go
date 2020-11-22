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

//ref 	nQuery.Scan(&n.ID, &n.PostID, &n.CommentID, &n.UserLostID, &n.voteState, &n.CreatedTime, &n.ToWhom) -> like this -> REFactor
// add interest func
//fix create user -> signup  after -> github signout?- password null
//fix -> like =0, dislike =1, like & dislike = 0 etc
//todo comment - for post -> notify

// 1 show notify  like/dislike post, comment, Lost comment by post
// 2 activity page -> show user created post?comment, liked, disliked post/comment
// 3 add - func, update/delete -> comment/post
//3.1 link another post -> show

//try - event -> add sound & confetti -Login
// save photo, like - source DB refactor
//config, router refactor
//if cookie = 0, notify message  user, logout etc

// 1 request, 910 additional, 0904 - 101202 ->
// 2 request -7575
// 3 request 910 additional, 090410 - 101202 ->Otegen batyr etc
