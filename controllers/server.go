package controllers

import (
	"ForumX/utils"
	"fmt"
	"log"
	"net/http"
)

//Middleware func wrapper
//cokie, check post, get, print log, print IN data from client
func Middleware(f http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		//valid Input data, and , logger
		c, _ := r.Cookie("_cookie")
		cookieBrowser := ""
		if c != nil {
			cookieBrowser = c.Value
		}
		isCookie, sessionF := utils.IsCookie(w, r, cookieBrowser)
		if isCookie {
			//write cookie value & session value - global variable
			CookieBrowser = c.Value
			session = sessionF
			fmt.Println("ok cookie have", session)
			f(w, r)
		}
	}
}

//Init func handlers
func Init() {
	const PORT = ":6969"
	//create multiplexer
	mux := http.NewServeMux()
	//file server
	mux.Handle("/statics/", http.StripPrefix("/statics/", http.FileServer(http.Dir("./statics/"))))

	mux.HandleFunc("/", GetAllPosts)
	mux.HandleFunc("/sapid", GetAllPosts)
	mux.HandleFunc("/love", GetAllPosts)
	mux.HandleFunc("/science", GetAllPosts)

	mux.HandleFunc("/post", GetPostByID)
	mux.HandleFunc("/create/post", Middleware(CreatePost))
	mux.HandleFunc("/edit/post", Middleware(UpdatePost))
	mux.HandleFunc("/delete/post", Middleware(DeletePost))

	mux.HandleFunc("/comment", Middleware(LeaveComment))
	mux.HandleFunc("/edit/comment", Middleware(UpdateComment))
	mux.HandleFunc("/delete/comment", Middleware(DeleteComment))
	mux.HandleFunc("/answer/comment", Middleware(AnswerComment))

	mux.HandleFunc("/votes/post", Middleware(VotesPost))
	mux.HandleFunc("/votes/comment", Middleware(VotesComment))

	mux.HandleFunc("/signin", Signin)
	mux.HandleFunc("/signup", Signup)
	mux.HandleFunc("/googleSignin", GoogleSignin)
	mux.HandleFunc("/googleUserInfo", GoogleUserData)

	mux.HandleFunc("/githubSignin", GithubSignin)
	mux.HandleFunc("/githubUserInfo", GithubUserData)
	mux.HandleFunc("/logout", Logout)

	mux.HandleFunc("/profile", Middleware(GetUserProfile))
	mux.HandleFunc("/user/id", Middleware(GetAnotherProfile))
	mux.HandleFunc("/edit/user", Middleware(UpdateProfile))
	mux.HandleFunc("/delete/account", Middleware(DeleteAccount))

	mux.HandleFunc("/activity", Middleware(GetUserActivities))
	mux.HandleFunc("/search", Search)
	// http.HandleFunc("/chat", routing.StartChat)
	log.Fatal(http.ListenAndServe(PORT, mux))
	fmt.Println("Listening port:", PORT)
}

//chaining
// func RequireAuthentication(next http.Handler) http.Handler {
// 	return http.HandlerFunc(
// 		func(w http.ResponseWriter, r *http.Request) {
// 			b, _ := utils.IsCookie(w, r)
// 			fmt.Print(b, "ccokie")
// 			if !b {
// 				http.Redirect(w, r, "/signin", 302)
// 				return
// 			}
// 			next.ServeHTTP(w, r)
// 		})
// }
