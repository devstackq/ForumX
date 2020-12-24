package controllers

import (
	"ForumX/utils"
	"fmt"
	"log"
	"net/http"
)

// high order function func(func)(callback)
func Middleware(f http.HandlerFunc) http.HandlerFunc {
	//анонимная функция вызывается, и делает логику, смотрит куки, и если надо вызовет хендлер, а отом вернет результат вызова анонимной фукнции
	//Коллбэки же позволяют нам быть уверенными в том, что определенный код не начнет исполнение до того момента, пока другой код не завершит исполнение.
	return func(w http.ResponseWriter, r *http.Request) {
		//valid Input data, and , logger
		c, err := r.Cookie("_cookie")
		if err != nil {
			log.Println(err, "expires timeout")
			utils.IsCookieExpiration(w, r, session)
			return
		}
		// var sid int
		// err = DB.QueryRow("SELECT id FROM session WHERE uuid = ?", c.Value).Scan(&sid)
		// if sid <= 0 {
		// 	fmt.Println("del cookie", sid)
		// 	utils.DeleteCookie(w)
		// }
		cookie := c.Value
		//check cookie, routting, then call handler -> middleware
		isCookie, sessionF := utils.IsCookie(w, r, cookie)
		//update Page call Middleware(getProfile)- > check current cookie(logouted user), == session(newCookie)
		if isCookie {
			//write cookie value & session value - global variable
			session = sessionF
			fmt.Println("ok cookie valid, can do operation", session)
			f(w, r)
		} else {
			//cokkie подменили или новый юзер зашел, isCookie, ref, || logic change
			fmt.Println("another cookie, uuid != session.Uuid Db ")
			//	utils.DeleteCookie(w)
			http.Redirect(w, r, "/signin", 302)
			//utils.IsCookieExpiration(w, r, session)
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
