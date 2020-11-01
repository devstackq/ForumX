package controllers

import (
	"log"
	"net/http"
)

// type Server struct {
// 	hadler  http.Handler,
// }

// func newServer() {
// 	serverMux
// }

// func (s *Server) Run() {

// }

//handlers
//mux own server,  route init  - google, config FileServer
func Init() {
	http.Handle("/statics/", http.StripPrefix("/statics/", http.FileServer(http.Dir("./statics/"))))

	http.HandleFunc("/", GetAllPosts)
	http.HandleFunc("/sapid", GetAllPosts)
	http.HandleFunc("/love", GetAllPosts)
	http.HandleFunc("/science", GetAllPosts)

	http.HandleFunc("/post", GetPostByID)
	http.HandleFunc("/create/post", CreatePost)
	http.HandleFunc("/edit/post", UpdatePost)
	http.HandleFunc("/delete/post", DeletePost)

	http.HandleFunc("/comment", LeaveComment)

	http.HandleFunc("/votes", LostVotes)
	http.HandleFunc("/votes/comment", LostVotesComment)

	http.HandleFunc("/search", Search)

	http.HandleFunc("/profile", GetUserProfile)
	http.HandleFunc("/user/id/", GetAnotherProfile)
	http.HandleFunc("/edit/user", UpdateProfile)

	http.HandleFunc("/signup", Signup)
	http.HandleFunc("/signin", Signin)
	http.HandleFunc("/logout", Logout)

	// http.HandleFunc("/chat", routing.StartChat)
	log.Fatal(http.ListenAndServe(":6969", nil))
}
