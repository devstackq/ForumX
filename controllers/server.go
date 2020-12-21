package controllers

import (
	"ForumX/utils"
	"fmt"
	"log"
	"net/http"
)

//Init func handlers
//cokie, check post, get, print log, print IN data from client
func Init() {
	mux := http.NewServeMux()
	
	//mux.Handle("/", http.HandleFunc(utils.IsCookie(GetAllPosts))

	//mux.Handle("/", utils.IsCookie(GetAllPosts))
	mux.HandleFunc("/", GetAllPosts)
	mux.HandleFunc("/sapid",GetAllPosts)
	mux.HandleFunc("/love",GetAllPosts)
	mux.HandleFunc("/science", GetAllPosts)
	mux.HandleFunc("/signin",Signin)
	// fs := http.StripPrefix("/statics/", http.FileServer(http.Dir("./statics/")))
	// mux.Handle("/statics/", fs)

	routers := RequireAuthentication(mux)

	//mux.Handle("/", http.HandlerFunc( GetAllPosts))

	// mux.Handle("/post", http.HandlerFunc(GetPostByID))
	// mux.Handle("/create/post", http.HandlerFunc(CreatePost))
	// mux.Handle("/edit/post", http.HandlerFunc(UpdatePost))

	// mux.Handle("/delete/post", http.HandlerFunc(DeletePost))
	// mux.Handle("/activity", http.HandlerFunc(GetUserActivities))
	// mux.Handle("/comment", http.HandlerFunc(LeaveComment))

	// mux.Handle("/edit/comment", http.HandlerFunc(UpdateComment))
	// mux.Handle("/delete/comment", http.HandlerFunc(DeleteComment))
	// mux.Handle("/answer/comment", http.HandlerFunc(AnswerComment))

	// mux.Handle("/votes/post", http.HandlerFunc(VotesPost))
	// mux.Handle("/votes/comment", http.HandlerFunc(VotesComment))
	// mux.Handle("/search", http.HandlerFunc(Search))

	// mux.Handle("/profile", http.HandlerFunc(GetUserProfile))
	// mux.Handle("/user/id", http.HandlerFunc(GetAnotherProfile))
	// mux.Handle("/edit/user", http.HandlerFunc(UpdateProfile))
	// mux.Handle("/delete/account", http.HandlerFunc(DeleteAccount))

	// mux.Handle("/signup", http.HandlerFunc(Signup))
	// mux.Handle("/signin", http.HandlerFunc(Signin))
	// mux.Handle("/logout", http.HandlerFunc(Logout))

	// mux.Handle("/googleSignin", http.HandlerFunc(GoogleSignin))
	// mux.Handle("/googleUserInfo", http.HandlerFunc(GoogleUserData))

	// mux.Handle("/githubSignin", http.HandlerFunc(GithubSignin))
	// mux.Handle("/githubUserInfo", http.HandlerFunc(GithubUserData))

	// http.HandleFunc("/chat", routing.StartChat)
	fmt.Println("Listening port: 6969")
	log.Fatal(http.ListenAndServe(":6969", routers))
}

// func Middleware(handleW http.Handler) {
//     utils.IsCookie(http.ResponseWriter, r *http.Request)
// }

func RequireAuthentication(next http.Handler) http.Handler {
	// We wrap our anonymous function, and cast it to a http.HandlerFunc
	// Because our function signature matches ServeHTTP(w, r), this allows
	// our function (type) to implicitly satisify the http.Handler interface.
	return http.HandlerFunc(
	  func(w http.ResponseWriter, r *http.Request) {
		// Logic before - reading request values, putting things into the
		// request context, performing authentication
  
		// Important that we call the 'next' handler in the chain. If we don't,
		// then request handling will stop here.
		b, _ := utils.IsCookie(w,r)
		fmt.Print(b, "ccokie")
		if !b {
            http.Redirect(w, r, "/signin", 302)
			return
        }

		next.ServeHTTP(w, r)
		// Logic after - useful for logging, metrics, etc.
		//
		// It's important that we don't use the ResponseWriter after we've called the
		// next handler: we may cause conflicts when trying to write the response
	  })
  }

