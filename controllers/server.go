package controllers

import (
	"log"
	"net/http"
)

//routes -> handlers
//mux own server,  route init  - google, config FileServer

//Init func handlers
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
	http.HandleFunc("/activity", GetUserActivities)

	http.HandleFunc("/comment", LeaveComment)
	http.HandleFunc("/edit/comment", UpdateComment)
	http.HandleFunc("/delete/comment", DeleteComment)
	http.HandleFunc("/answer/comment", AnswerComment)
	http.HandleFunc("/reply/answer/comment", ReplyAnswer)

	http.HandleFunc("/votes", VotesPost)
	http.HandleFunc("/votes/comment", VotesComment)
	http.HandleFunc("/search", Search)

	http.HandleFunc("/profile", GetUserProfile)
	http.HandleFunc("/user/id/", GetAnotherProfile)
	http.HandleFunc("/edit/user", UpdateProfile)
	http.HandleFunc("/delete/account", DeleteAccount)

	http.HandleFunc("/signup", Signup)
	http.HandleFunc("/signin", Signin)
	http.HandleFunc("/logout", Logout)

	http.HandleFunc("/googleSignin", GoogleSignin)
	http.HandleFunc("/googleUserInfo", GoogleUserData)

	http.HandleFunc("/githubSignin", GithubSignin)
	http.HandleFunc("/githubUserInfo", GithubUserData)
	// http.HandleFunc("/chat", routing.StartChat)
	log.Fatal(http.ListenAndServe(":6969", nil))
}
