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

// add interest func - adv feat
//fix create user -> signup  after -> github signout?- password null

//fix normal show - liede/disleked post
// 3 add - func, update/delete -> comment
// like/dislike -> comment -> uid, pid, uid

//try - event -> add sound & confetti -Login
// save photo, like - source DB refactor
//config, router refactor
//if cookie = 0, notify message  user, logout etc

// 1 request, 910 additional, 0904 - 101202 ->
// 2 request -7575
// 3 request 910 additional, 090410 - 101202 ->Otegen batyr etc

// {{if gt .CID 0 }}

//         <form action="/user/id/?{{.UID}}" method="post" class="author-post">
//             <input type="hidden" name="uid" value="{{.UID}}">
//             <input type="submit" value="{{.UserLost}}">
//         </form>

//         {{if eq .VoteState 1 }}  <a href="/post?id={{.CID}}"> Liked comment</a>in Coment  "{{.CommentTitle}} "{{end}}
//         {{if eq .VoteState 2 }}  <a href="/post?id={{.CID}}"> Disliked comment</a>in Coment  "{{.CommentTitle}} "{{end}}

//          {{end}}
