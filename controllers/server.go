package controllers

import (
	"ForumX/utils"
	"log"
	"net/http"
)

//анонимная функция вызывается, и делает логику, смотрит куки, и если надо вызовет хендлер, а отом вернет результат вызова анонимной фукнции
//Коллбэки же позволяют нам быть уверенными в том, что определенный код не начнет исполнение до того момента, пока другой код не завершит исполнение.
// high order function func(func)(callback)
func IsValidCookie(f http.HandlerFunc) http.HandlerFunc {
	
	return func(w http.ResponseWriter, r *http.Request) {
		//check expires cookie
		c, err := r.Cookie("_cookie")
		if err != nil {
			log.Println(err, "expires timeout || cookie deleted")
			utils.Logout(w, r, session)
			return
		}
		// then call handler -> middleware
		if isValidCookie, sessionF := utils.IsCookie(w, r, c.Value); isValidCookie {
			//write session data - global variable
			session = &sessionF
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
	mux.HandleFunc("/create/post", IsValidCookie(CreatePost))
	mux.HandleFunc("/edit/post", IsValidCookie(UpdatePost))
	mux.HandleFunc("/delete/post", IsValidCookie(DeletePost))

	mux.HandleFunc("/comment", IsValidCookie(LeaveComment))
	mux.HandleFunc("/edit/comment", IsValidCookie(UpdateComment))
	mux.HandleFunc("/delete/comment", IsValidCookie(DeleteComment))
	mux.HandleFunc("/reply/comment/replyId", IsValidCookie(ReplyComment))

	mux.HandleFunc("/votes/post", IsValidCookie(VotesPost))
	mux.HandleFunc("/votes/comment", IsValidCookie(VotesComment))

	mux.HandleFunc("/signin", Signin)
	mux.HandleFunc("/signup", Signup)
	mux.HandleFunc("/googleSignin", GoogleSignin)
	mux.HandleFunc("/googleUserInfo", GoogleUserData)

	mux.HandleFunc("/githubSignin", GithubSignin)
	mux.HandleFunc("/githubUserInfo", GithubUserData)
	mux.HandleFunc("/logout", IsValidCookie(Logout))

	mux.HandleFunc("/profile", IsValidCookie(GetUserProfile))
	mux.HandleFunc("/user/id", IsValidCookie(GetAnotherProfile))
	mux.HandleFunc("/edit/user", IsValidCookie(UpdateProfile))
	mux.HandleFunc("/delete/account", IsValidCookie(DeleteAccount))


	mux.HandleFunc("/activity", IsValidCookie(GetUserActivities))
	mux.HandleFunc("/search", Search)
	// http.HandleFunc("/chat", routing.StartChat)
	log.Println("Listening port:", PORT)
	log.Fatal(http.ListenAndServe(PORT, mux))
}
