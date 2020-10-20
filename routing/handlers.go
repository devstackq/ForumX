package routing

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"net/http"
	"strconv"

	"github.com/devstackq/ForumX/models"
	util "github.com/devstackq/ForumX/utils"
	"golang.org/x/crypto/bcrypt"
)

var (
	err error
	DB  *sql.DB
	API struct{ Message string }
)

//getAllPost and Posts by cateogry
//receive request, from client, query params, category ID, then query DB, depends catID, get Post this catID
func GetAllPosts(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" && r.URL.Path != "/science" && r.URL.Path != "/love" && r.URL.Path != "/sapid" {
		models.DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}

	posts, endpoint, err := models.GetAllPost(r)
	if err != nil {
		log.Fatal(err)
	}

	models.DisplayTemplate(w, "header", util.IsAuth(r))

	// endpoint -> get post by category
	// profile/ fix, create, get post fix
	if endpoint == "/" {
		models.DisplayTemplate(w, "index", posts)
	} else {
		models.DisplayTemplate(w, "catTemp", posts)
	}
}

//view 1 post by id
func GetPostById(w http.ResponseWriter, r *http.Request) {

	// if r.URL.Path != "/post" {
	// 	models.DisplayTemplate(w, "404page", http.StatusNotFound)
	// 	return
	// }
	//check cookie for  navbar, if not cookie - signin

	comments, post, err := models.GetPostById(r)
	if err != nil {
		panic(err)
	}
	models.DisplayTemplate(w, "header", util.IsAuth(r))
	models.DisplayTemplate(w, "posts", post)
	models.DisplayTemplate(w, "comment", comments)
}

//create post
func CreatePost(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/create/post" {
		models.DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}

	API.Message = ""

	switch r.Method {
	case "GET":
		models.DisplayTemplate(w, "header", util.IsAuth(r))
		models.DisplayTemplate(w, "create", &API.Message)
	case "POST":
		access, session := util.CheckForCookies(w, r)
		log.Println(access, "access status")
		if !access {
			http.Redirect(w, r, "/signin", 302)
			return
		}
		r.ParseForm()
		r.ParseMultipartForm(10 << 20)
		f, _, _ := r.FormFile("uploadfile")
		f2, _, _ := r.FormFile("uploadfile")
		categories, _ := r.Form["input"]
		post := models.Posts{
			Title:      r.FormValue("title"),
			Content:    r.FormValue("content"),
			Categories: categories,
			FileS:      f,
			FileI:      f2,
			Session:    session,
		}
		models.CreatePosts(w, r, post)
		http.Redirect(w, r, "/", http.StatusOK)
	}
}

//update post
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	var pid int
	if r.Method == "GET" {

		pid, _ = strconv.Atoi(r.URL.Query().Get("id"))
		p := models.Posts{}
		p.PostIDEdit = pid
		models.DisplayTemplate(w, "updatepost", p)

	}
	if r.Method == "POST" {
		access, _ := util.CheckForCookies(w, r)
		if !access {
			http.Redirect(w, r, "/signin", 302)
			return
		}

		r.ParseForm()
		r.ParseMultipartForm(10 << 20)
		file, _, err := r.FormFile("uploadfile")

		if err != nil {
			panic(err)

		}
		defer file.Close()

		fileBytes, err := ioutil.ReadAll(file)

		if err != nil {
			panic(err)
		}

		pid := r.FormValue("pid")
		pidnum, _ := strconv.Atoi(pid)

		p := models.Posts{
			Title:      r.FormValue("title"),
			Content:    r.FormValue("content"),
			Image:      fileBytes,
			PostIDEdit: pidnum,
		}

		err = p.UpdatePost()

		if err != nil {
			panic(err.Error())
		}
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

//delete post
func DeletePost(w http.ResponseWriter, r *http.Request) {

	var pid int

	pid, _ = strconv.Atoi(r.URL.Query().Get("id"))
	p := models.Posts{}
	p.PostIDEdit = pid

	access, _ := util.CheckForCookies(w, r)
	if !access {
		http.Redirect(w, r, "/signin", 302)
		return
	}

	err = p.DeletePost()

	if err != nil {
		panic(err.Error())
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

//create comment
func CreateComment(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/comment" {
		models.DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}

	if r.Method == "POST" {

		access, s := util.CheckForCookies(w, r)
		if !access {
			http.Redirect(w, r, "/signin", 302)
			return
		}

		r.ParseForm()
		//c, _ := r.Cookie("_cookie")
		//s := models.Session{UUID: c.Value}
		DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).Scan(&s.UserID)

		pid, _ := strconv.Atoi(r.FormValue("curr"))
		comment := r.FormValue("comment-text")

		checkLetter := false
		for _, v := range comment {
			if v >= 97 && v <= 122 || v >= 65 && v <= 90 && v >= 32 && v <= 64 || v > 128 {
				checkLetter = true
			}
		}

		if checkLetter {
			com := models.Comments{
				Commentik: comment,
				PostID:    pid,
				UserID:    s.UserID,
			}

			err = com.LostComment()

			if err != nil {
				panic(err.Error())

			}
		}
		http.Redirect(w, r, "post?id="+r.FormValue("curr"), 301)
	}
}

//profile current -> user page
func GetProfileById(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/profile" {
		models.DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}
	if r.Method == "GET" {

		//if userId now, createdPost uid equal -> show
		likedpost, posts, comments, user, err := models.GetUserProfile(r, w)
		if err != nil {
			panic(err)
		}

		models.DisplayTemplate(w, "header", util.IsAuth(r))
		models.DisplayTemplate(w, "profile", user)
		models.DisplayTemplate(w, "likedpost", likedpost)
		models.DisplayTemplate(w, "postuser", posts)
		models.DisplayTemplate(w, "commentuser", comments)

		cookie, _ := r.Cookie("_cookie")
		//delete coookie db
		go func() {
			for now := range time.Tick(30 * time.Minute) {
				checkCookieLife(now, cookie, w, r)
				//next logout each 10 min
				time.Sleep(30 * time.Minute)
			}
		}()
	}
}

func checkCookieLife(t time.Time, cookie *http.Cookie, w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == "_cookie" {
			models.Logout(w, r)
			return
		}
	}
}

//user page, other anyone
func GetUserById(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		posts, user, err := models.GetOtherUser(r)
		if err != nil {
			panic(err)
		}

		models.DisplayTemplate(w, "header", util.IsAuth(r))
		models.DisplayTemplate(w, "user", user)
		models.DisplayTemplate(w, "postuser", posts)
	}
}

//update profile
func UpdateProfile(w http.ResponseWriter, r *http.Request) {

	//check cookie for  navbar

	if r.Method == "GET" {
		models.DisplayTemplate(w, "header", util.IsAuth(r))
		models.DisplayTemplate(w, "updateuser", "")
	}

	if r.Method == "POST" {
		access, s := util.CheckForCookies(w, r)
		if !access {
			http.Redirect(w, r, "/signin", 302)
			return
		}

		r.ParseForm()
		r.ParseMultipartForm(10 << 20)
		file, _, err := r.FormFile("uploadfile")

		if err != nil {
			panic(err)
		}

		defer file.Close()

		fileBytes, err := ioutil.ReadAll(file)

		if err != nil {
			panic(err)
		}

		//		c, _ := r.Cookie("_cookie")
		//s := models.Session{UUID: c.Value}

		DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).
			Scan(&s.UserID)

		is, _ := strconv.Atoi(r.FormValue("age"))

		p := models.Users{
			FullName: r.FormValue("fullname"),
			Age:      is,
			Sex:      r.FormValue("sex"),
			City:     r.FormValue("city"),
			Image:    fileBytes,
			ID:       s.UserID,
		}

		err = p.UpdateProfile()

		if err != nil {
			panic(err.Error())
		}
	}

	http.Redirect(w, r, "/profile", http.StatusFound)
}

//search
func Search(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/search" {
		models.DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}

	if r.Method == "GET" {
		models.DisplayTemplate(w, "search", http.StatusFound)
	}

	if r.Method == "POST" {

		findPosts, err := models.Search(w, r)

		if err != nil {
			panic(err)
		}

		models.DisplayTemplate(w, "header", util.IsAuth(r))
		models.DisplayTemplate(w, "index", findPosts)
	}
}

//signup system
func Signup(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/signup" {
		models.DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}
	msg := models.API

	if r.Method == "GET" {
		models.DisplayTemplate(w, "signup", &msg)
	}

	if r.Method == "POST" {

		fn := r.FormValue("fullname")
		e := r.FormValue("email")
		p := r.FormValue("password")
		a := r.FormValue("age")
		s := r.FormValue("sex")
		c := r.FormValue("city")

		r.ParseMultipartForm(10 << 20)
		file, _, err := r.FormFile("uploadfile")

		if err != nil {
			panic(err)
		}

		defer file.Close()

		fB, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(p), 8)
		if err != nil {
			panic(err)
		}

		//check email by unique, if have same email
		checkEmail, err := DB.Query("SELECT email FROM users")
		if err != nil {
			panic(err)
		}

		all := []models.Users{}

		for checkEmail.Next() {
			user := models.Users{}
			var email string
			err = checkEmail.Scan(&email)
			if err != nil {
				panic(err.Error)
			}

			user.Email = email
			all = append(all, user)
		}

		for _, v := range all {
			if v.Email == e {
				API.Message = "Not unique email lel"
				models.DisplayTemplate(w, "signup", &API.Message)
				return
			}
		}

		_, err = DB.Exec("INSERT INTO users( full_name, email, password, age, sex, city, image) VALUES (?, ?, ?, ?, ?, ?, ?)",
			fn, e, hash, a, s, c, fB)

		if err != nil {
			panic(err.Error())
		}

		http.Redirect(w, r, "/signin", 301)
	}
}

//signin system
func Signin(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/signin" {
		models.DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}
	r.Header.Add("Accept", "text/html")
	r.Header.Add("User-Agent", "MSIE/15.0")

	API.Message = ""

	if r.Method == "GET" {
		models.DisplayTemplate(w, "signin", &API.Message)
	}

	if r.Method == "POST" {
		var person models.Users
		//b, _ := ioutil.ReadAll(r.Body)
		err := json.NewDecoder(r.Body).Decode(&person)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(person)

		if person.Type == "default" {
			fmt.Println(" default auth")

			email := person.Email
			pwd := person.Password

			models.Signin(w, r, email, pwd)
			http.Redirect(w, r, "/profile", 200)
			//http.Redirect(w, r, "/profile", 200)
			//http.Redirect(w, r, "/profile", 200)
			//citiesArtist := FindCityArtist(w, r, strings.ToLower(string(body)))
			//w.Header().Set("Content-Type", "application/json")
			//json.NewEncoder(w).Encode(citiesArtist)

			// json.NewEncoder(w).Encode(msg)
			// ok := "okay"
			// b := []byte(ok)

			// msg.Msg = "okay"
			// w.Header().Set("Content-Type", "application/json")
			// m, _ := json.Marshal(msg)
			// w.Write(m)
		} else if person.Type == "google" {
			fmt.Println("todo google auth")
			http.Redirect(w, r, "/profile", http.StatusFound)
		} else if person.Type == "github" {
			fmt.Println("todo github")
			http.Redirect(w, r, "/profile", http.StatusFound)
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")

	}
}

// Logout
func Logout(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/logout" {
		models.DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}
	if r.Method == "GET" {
		models.Logout(w, r)
	}
}

//like dislike post
func LostVotes(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/votes" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	access, s := util.CheckForCookies(w, r)
	if !access {
		http.Redirect(w, r, "/signin", 302)
		return
	}
	//c, _ := r.Cookie("_cookie")
	//s := models.Session{UUID: c.Value}

	DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).
		Scan(&s.UserID)

	pid := r.URL.Query().Get("id")
	lukas := r.FormValue("lukas")
	diskus := r.FormValue("diskus")

	if r.Method == "POST" {

		if lukas == "1" {
			//check if not have post and user lost vote this post
			//1 like or 1 dislike 1 user lost 1 post, get previus value and +1
			var p, u int
			err = DB.QueryRow("SELECT post_id, user_id FROM likes WHERE post_id=? AND user_id=?", pid, s.UserID).Scan(&p, &u)

			if p == 0 && u == 0 {

				oldlike := 0
				err = DB.QueryRow("SELECT count_like FROM posts WHERE id=?", pid).Scan(&oldlike)
				nv := oldlike + 1
				_, err = DB.Exec("UPDATE  posts SET count_like = ? WHERE id= ?", nv, pid)
				if err != nil {
					panic(err)
				}

				_, err = DB.Exec("INSERT INTO likes(post_id, user_id) VALUES( ?, ?)", pid, s.UserID)
				if err != nil {
					panic(err)
				}
			}
		}

		if diskus == "1" {

			var p, u int
			err = DB.QueryRow("SELECT post_id, user_id FROM likes WHERE post_id=? AND user_id=?", pid, s.UserID).Scan(&p, &u)

			if p == 0 && u == 0 {

				oldlike := 0
				err = DB.QueryRow("select count_dislike from posts where id=?", pid).Scan(&oldlike)
				nv := oldlike + 1
				_, err = DB.Exec("UPDATE  posts SET count_dislike = ? WHERE id= ?", nv, pid)
				if err != nil {
					panic(err)
				}
				_, err = DB.Exec("INSERT INTO likes(post_id, user_id) VALUES( ?, ?)", pid, s.UserID)

				if err != nil {
					panic(err)
				}
			}
		}
	}
	http.Redirect(w, r, "post?id="+pid, 301)
}

func LostVotesComment(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/votes/comment" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	access, s := util.CheckForCookies(w, r)
	if !access {
		http.Redirect(w, r, "/signin", 302)
		return
	}
	//c, _ := r.Cookie("_cookie")
	//s := models.Session{UUID: c.Value}
	DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).
		Scan(&s.UserID)

	cid := r.URL.Query().Get("cid")
	comdis := r.FormValue("comdis")
	comlike := r.FormValue("comlike")

	pidc := r.FormValue("pidc")

	if r.Method == "POST" {

		if comlike == "1" {

			var c, u int
			err = DB.QueryRow("SELECT comment_id, user_id FROM likes WHERE comment_id=? AND user_id=?", cid, s.UserID).Scan(&c, &u)

			if c == 0 && u == 0 {

				oldlike := 0
				err = DB.QueryRow("SELECT com_like FROM comments WHERE id=?", cid).Scan(&oldlike)
				nv := oldlike + 1

				_, err = DB.Exec("UPDATE  comments SET com_like = ? WHERE id= ?", nv, cid)

				if err != nil {
					panic(err)
				}

				_, err = DB.Exec("INSERT INTO likes(comment_id, user_id) VALUES( ?, ?)", cid, s.UserID)
				if err != nil {
					panic(err)
				}
			}
		}

		if comdis == "1" {

			var c, u int
			err = DB.QueryRow("SELECT comment_id, user_id FROM likes WHERE comment_id=? AND user_id=?", cid, s.UserID).Scan(&c, &u)

			if c == 0 && u == 0 {

				oldlike := 0
				err = DB.QueryRow("SELECT com_dislike FROM comments WHERE id=?", cid).Scan(&oldlike)
				nv := oldlike + 1

				_, err = DB.Exec("UPDATE  comments SET com_dislike = ? WHERE id= ?", nv, cid)

				if err != nil {
					panic(err)
				}

				_, err = DB.Exec("INSERT INTO likes(comment_id, user_id) VALUES( ?, ?)", cid, s.UserID)
				if err != nil {
					panic(err)
				}
			}
		}
		http.Redirect(w, r, "/post?id="+pidc, 301)
	}
}

//Likes table, filed posrid, userid, state_id
// 0,1,2 if state ==0, 1 || 2,
// next btn, if 1 == 1, state =0
