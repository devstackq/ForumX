package controllers

import (
	"ForumX/utils"
	"fmt"
	"log"
	"net/http"
)

//анонимная функция вызывается, и делает логику, смотрит куки, и если надо вызовет хендлер, а отом вернет результат вызова анонимной фукнции
//Коллбэки же позволяют нам быть уверенными в том, что определенный код не начнет исполнение до того момента, пока другой код не завершит исполнение.
// high order function func(func)(callback)
func Middleware(f http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		//check expires cookie
		c, err := r.Cookie("_cookie")
		if err != nil {
			log.Println(err, "expires timeout")
			utils.Logout(w, r, session)
			return
		}
		// then call handler -> middleware
		if isValidCookie, sessionF := utils.IsCookie(w, r, c.Value); isValidCookie {
			//write session data - global variable
			session = sessionF
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
	mux.HandleFunc("/logout", Middleware(Logout))

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
