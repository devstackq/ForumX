package models

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/devstackq/ForumX/model"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	err  error
	DB   *sql.DB
	rows *sql.Rows
)

type Users struct {
	ID          int
	FullName    string
	Email       string
	Password    string
	IsAdmin     bool
	Age         int
	Sex         string
	CreatedTime time.Time
	City        string
	Image       []byte
	ImageHtml   string
	Role        string
	SVG         bool
	Type        string
}

type Category struct {
	ID     int
	Name   string
	UserID int
}

type Filter struct {
	Category string
	Like     string
	Date     string
}

type PostCategory struct {
	PostID   int64
	Category string
}

//comment ID -> foreign key -> postID
type Comments struct {
	ID             int
	Commentik      string
	PostID         int
	UserID         int
	CreatedTime    time.Time
	AuthorComment  string
	CommentLike    int
	CommentDislike int
}

var API struct {
	Authenticated bool
}

type Likes struct {
	ID      int
	Like    int
	Dislike int
	PostID  int
	UserID  int
	Voted   bool
}
type Notify struct {
	Message string
}

//get data from client, put data in Handler, then models -> query db
func (c *Comments) LostComment() error {

	_, err := DB.Exec("INSERT INTO comments( content, post_id, user_idx) VALUES(?,?,?)",
		c.Commentik, c.PostID, c.UserID)
	if err != nil {
		panic(err.Error())
	}
	return nil
}

func (pcb *PostCategory) CreateBridge() error {
	_, err := DB.Exec("INSERT INTO post_cat_bridge (post_id, category) VALUES (?, ?)",
		pcb.PostID, pcb.Category)
	if err != nil {
		return err
	}
	return nil
}

//update profile
func (u *Users) UpdateProfile() error {

	_, err := DB.Exec("UPDATE  users SET full_name=?, age=?, sex=?, city=?, image=? WHERE id =?",
		u.FullName, u.Age, u.Sex, u.City, u.Image, u.ID)
	if err != nil {
		return err
	}
	return nil
}

//siginin
func Signin(w http.ResponseWriter, r *http.Request, email, password string) {

	u := DB.QueryRow("SELECT id, password FROM users WHERE email=?", email)

	var user Users
	var err error
	//check pwd, if not correct, error
	err = u.Scan(&user.ID, &user.Password)
	if err != nil {
		authError(w, err, "user not found")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		authError(w, err, "password incorrect")
		return
	}
	//get user by Id, and write session struct
	s := model.Session{
		UserID: user.ID,
	}
	uuid := uuid.Must(uuid.NewV4(), err).String()
	if err != nil {
		authError(w, err, "uuid trouble")
		return
	}

	//create uuid and set uid DB table session by userid,
	_, err = DB.Exec("INSERT INTO session(uuid, user_id) VALUES (?, ?)", uuid, s.UserID)
	if err != nil {
		authError(w, err, "the user is already in the system")
		return
	}

	// get user in info by session Id
	err = DB.QueryRow("SELECT id, uuid FROM session WHERE user_id = ?", s.UserID).Scan(&s.ID, &s.UUID)
	if err != nil {
		authError(w, err, "not find user from session")
		return
	}
	//set cookie 9128ueq9widjaisdh238yrhdeiuwandijsan
	// Crete post -> Cleint cookie == session, Userd
	cookie := http.Cookie{
		Name:     "_cookie",
		Value:    s.UUID,
		Path:     "/",
		Expires:  time.Now().Add(30 * time.Minute),
		HttpOnly: false,
	}
	http.SetCookie(w, &cookie)
	authError(w, nil, "success")
}

func authError(w http.ResponseWriter, err error, text string) {
	fmt.Println(text, "errka")
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		m, _ := json.Marshal(text)
		w.Write(m)
		return
	} else {
		w.WriteHeader(200)
		m, _ := json.Marshal(text)
		w.Write(m)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("_cookie")
	if err != nil {
		fmt.Println(err, "cookie err")
	}
	//add cookie -> fields uuid
	s := model.Session{UUID: cookie.Value}
	//get ssesion id, by local struct uuid
	DB.QueryRow("SELECT id FROM session WHERE uuid = ?", s.UUID).
		Scan(&s.ID)
	fmt.Println(s.ID, "id del session")
	//delete session by id session
	_, err = DB.Exec("DELETE FROM session WHERE id = ?", s.ID)

	if err != nil {
		panic(err)
	}

	// then delete cookie from client
	cDel := http.Cookie{
		Name:     "_cookie",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: false,
	}
	http.SetCookie(w, &cDel)
	http.Redirect(w, r, "/", http.StatusOK)

}

//get profile by id
func GetUserProfile(r *http.Request, w http.ResponseWriter) ([]Posts, []Posts, []Comments, Users, error) {

	cookie, _ := r.Cookie("_cookie")

	//time.AfterFunc(10, checkCookieLife(cookie, w, r))
	s := model.Session{UUID: cookie.Value}
	u := Users{}
	DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).Scan(&s.UserID)
	lps := []Likes{}
	lp, err := DB.Query("select post_id from likes where user_id =?", s.UserID)
	fmt.Print("her1")
	for lp.Next() {
		l := Likes{}
		var lpid int

		err = lp.Scan(&lpid)
		fmt.Print("her2")
		l.PostID = lpid
		lps = append(lps, l)
	}

	err = DB.QueryRow("SELECT * FROM users WHERE id = ?", s.UserID).Scan(&u.ID, &u.FullName, &u.Email, &u.Password, &u.IsAdmin, &u.Age, &u.Sex, &u.CreatedTime, &u.City, &u.Image)
	if u.Image[0] == 60 {
		u.SVG = true
	}
	fmt.Print("her3", err)

	encStr := base64.StdEncoding.EncodeToString(u.Image)
	u.ImageHtml = encStr

	var likedpost *sql.Rows
	LikedPosts := []Posts{}
	var can []int

	for _, v := range lps {
		can = append(can, v.PostID)
	}
	fmt.Print("her4")

	//unique liked post by user
	fin := isUnique(can)
	fmt.Print("her5")
	//accum liked post
	for _, v := range fin {
		//get each liked post by ID, then likedpost, puth array post
		likedpost, err = DB.Query("SELECT * FROM posts WHERE id=?", v)

		for likedpost.Next() {

			post := Posts{}

			var id, creatorid, countlike, countdislike int
			var content, title string
			var creationtime time.Time
			var image []byte

			err = likedpost.Scan(&id, &title, &content, &creatorid, &creationtime, &image, &countlike, &countdislike)
			if err != nil {
				panic(err.Error)
			}

			post.ID = id
			post.Title = title
			post.Content = content
			post.CreatorID = creatorid
			post.CreatedTime = creationtime
			post.Image = image
			post.CountLike = countlike
			post.CountDislike = countdislike

			LikedPosts = append(LikedPosts, post)
		}
	}
	fmt.Print("her6")
	psu, err := DB.Query("SELECT * FROM posts WHERE creator_id=?", s.UserID)

	PostsCreatedUser := []Posts{}

	for psu.Next() {

		post := Posts{}

		var id, creatorid, countlike, countdislike int
		var content, title string
		var creationtime time.Time
		var image []byte

		err = psu.Scan(&id, &title, &content, &creatorid, &creationtime, &image, &countlike, &countdislike)
		if err != nil {
			panic(err.Error)
		}

		post.AuthorForPost = s.UserID

		post.ID = id
		post.Title = title
		post.Content = content
		post.CreatorID = creatorid
		post.CreatedTime = creationtime
		post.Image = image
		post.CountLike = countlike
		post.CountDislike = countdislike

		PostsCreatedUser = append(PostsCreatedUser, post)
	}
	fmt.Print("her7")
	csu, err := DB.Query("SELECT * FROM comments WHERE user_idx=?", s.UserID)

	CommentsLostUser := []Comments{}

	defer csu.Close()

	for csu.Next() {

		comment := Comments{}

		var id, postid, useridx, comLike, comDislike int
		var content string
		var createdtime time.Time

		err = csu.Scan(&id, &content, &postid, &useridx, &createdtime, &comLike, &comDislike)
		if err != nil {
			panic(err.Error)
		}

		comment.ID = id
		comment.PostID = postid
		comment.UserID = useridx
		comment.Commentik = content
		comment.CreatedTime = createdtime
		comment.CommentLike = comLike
		comment.CommentDislike = comDislike

		CommentsLostUser = append(CommentsLostUser, comment)
	}
	fmt.Print("her8")
	if err != nil {
		return nil, nil, nil, u, err
	}

	// /2020-10-18T17:33:35.013Z
	fmt.Print(cookie.Value, "cook")
	return LikedPosts, PostsCreatedUser, CommentsLostUser, u, nil
}

//find unique liked post
func isUnique(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

//get other user
func GetOtherUser(r *http.Request) ([]Posts, Users, error) {

	r.ParseForm()
	uid := r.FormValue("uid")

	user := DB.QueryRow("SELECT * FROM users WHERE id = ?", uid)
	u := Users{}
	err = user.Scan(&u.ID, &u.FullName, &u.Email, &u.Password, &u.IsAdmin, &u.Age, &u.Sex, &u.CreatedTime, &u.City, &u.Image)
	if u.Image[0] == 60 {
		u.SVG = true
	}
	encStr := base64.StdEncoding.EncodeToString(u.Image)
	u.ImageHtml = encStr
	psu, err := DB.Query("SELECT * FROM posts WHERE creator_id=?", u.ID)

	PostsOtherUser := []Posts{}

	defer psu.Close()

	for psu.Next() {

		post := Posts{}
		var id, creatorid, countlike, countdislike int
		var content, title string
		var creationtime time.Time
		var image []byte

		err = psu.Scan(&id, &title, &content, &creatorid, &creationtime, &image, &countlike, &countdislike)
		if err != nil {
			panic(err.Error)
		}

		post.ID = id
		post.Title = title
		post.Content = content
		post.CreatorID = creatorid
		post.CreatedTime = creationtime
		post.CountLike = countlike
		post.CountDislike = countdislike
		PostsOtherUser = append(PostsOtherUser, post)
	}
	if err != nil {
		return nil, u, err
	}
	return PostsOtherUser, u, nil
}

//search
func Search(w http.ResponseWriter, r *http.Request) ([]Posts, error) {

	keyword := r.FormValue("search")
	psu, err := DB.Query("SELECT * FROM posts WHERE title LIKE ?", "%"+keyword+"%")
	found := []Posts{}
	defer psu.Close()

	for psu.Next() {

		var id, creatorID, like, dislike int
		var content, title string
		var createdTime time.Time
		var image []byte

		err = psu.Scan(&id, &title, &content, &creatorID, &createdTime, &image, &like, &dislike)

		if err != nil {
			panic(err.Error)
		}
		found = putPost(id, title, content, creatorID, createdTime, image, like, dislike)

		// postX := Posts{
		// 	ID:      id,
		// 	Title:   title,
		// 	Content: &content,
		// }

		//PostsUser = append(PostsUser, postX)
	}

	if err != nil {
		return nil, err
	}
	return found, nil
}

func putPost(id int, title, content string, creatorID int, createdTime time.Time, image []byte, like, dislike int) []Posts {

	var post Posts
	var posts []Posts
	post.ID = id
	post.Title = title
	post.Content = content
	post.CreatorID = creatorID
	post.CreatedTime = createdTime
	post.Image = image
	post.CountLike = like
	post.CountDislike = dislike
	posts = append(posts, post)
	return posts
}

putPost use func