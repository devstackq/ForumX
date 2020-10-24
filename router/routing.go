package routers

import (
	"log"
	"net/http"

	"github.com/devstackq/ForumX/routing"
)

func init() {

	http.Handle("/statics/", http.StripPrefix("/statics/", http.FileServer(http.Dir("./statics/"))))

	http.HandleFunc("/", routing.GetAllPosts)
	http.HandleFunc("/sapid", routing.GetAllPosts)
	http.HandleFunc("/love", routing.GetAllPosts)
	http.HandleFunc("/science", routing.GetAllPosts)

	http.HandleFunc("/post", routing.GetPostById)
	http.HandleFunc("/profile", routing.GetProfileById)
	http.HandleFunc("/user/id/", routing.GetUserById)

	http.HandleFunc("/comment", routing.CreateComment)
	http.HandleFunc("/create/post", routing.CreatePost)

	http.HandleFunc("/edit/post", routing.UpdatePost)
	http.HandleFunc("/delete/post", routing.DeletePost)
	http.HandleFunc("/edit/user", routing.UpdateProfile)

	http.HandleFunc("/signup", routing.Signup)
	http.HandleFunc("/signin", routing.Signin)
	http.HandleFunc("/logout", routing.Logout)

	http.HandleFunc("/votes", routing.LostVotes)
	http.HandleFunc("/votes/comment", routing.LostVotesComment)
	http.HandleFunc("/search", routing.Search)
	// http.HandleFunc("/chat", routing.StartChat)
	log.Fatal(http.ListenAndServe(":6969", nil))
}
