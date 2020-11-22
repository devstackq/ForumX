package models

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	structure "github.com/devstackq/ForumX/general"
	util "github.com/devstackq/ForumX/utils"
)

//Users struct
type User struct {
	ID          int
	FullName    string
	Email       string `json:"email"`
	Password    string
	IsAdmin     bool
	Age         int
	Sex         string
	CreatedTime time.Time
	City        string
	Image       []byte
	ImageHTML   string
	Role        string
	SVG         bool
	Type        string
	Temp        string
	Name        string `json:"name"`
	Location    string `json:"location"`
}

type Notify struct {
	ID          int
	PostID      int
	CommentID   int
	UserLostID  int
	voteState   int
	CreatedTime string
	ToWhom      int
	PostTitle string
	UserLost string
	CommentTitle string
}

//GetUserProfile function
func GetUserProfile(r *http.Request, w http.ResponseWriter, cookie *http.Cookie) ([]Post, []Post, []Comment, User, error) {

	//time.AfterFunc(10, checkCookieLife(cookie, w, r)) try check every 30 min cookie
	s := structure.Session{UUID: cookie.Value}
	u := User{}
	DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).Scan(&s.UserID)
	likedPostArr := []Votes{}

	likedpost, err := DB.Query("select post_id from voteState where user_id =? and like_state =?", s.UserID, 1)
	defer likedpost.Close()

	for likedpost.Next() {
		post := Votes{}
		var pid int
		err = likedpost.Scan(&pid)
		post.PostID = pid
		likedPostArr = append(likedPostArr, post)
	}

	err = DB.QueryRow("SELECT id, full_name, email, isAdmin, age, sex, created_time, city, image  FROM users WHERE id = ?", s.UserID).Scan(&u.ID, &u.FullName, &u.Email, &u.IsAdmin, &u.Age, &u.Sex, &u.CreatedTime, &u.City, &u.Image)
	//Age, sex, picture, city, date ?
	if err != nil {
		log.Println(err)
	}
	if u.Image[0] == 60 {
		u.SVG = true
	}
	u.Temp = u.CreatedTime.Format("2006 Jan _2 15:04:05")
	u.ImageHTML = base64.StdEncoding.EncodeToString(u.Image)

	var smtp *sql.Rows
	postsL := []Post{}

	var arrIDLiked []int

	for _, v := range likedPostArr {
		arrIDLiked = append(arrIDLiked, v.PostID)
	}

	//unique liked post by user
	fin := util.IsUnique(arrIDLiked)

	for _, v := range fin {
		//get each only  liked post by ID, then likedpost, put array post
		smtp, err = DB.Query("SELECT * FROM posts WHERE id=? and count_like > 0", v)
		if err != nil {
			log.Println(err)
		}
		for smtp.Next() {
			err = smtp.Scan(&id, &title, &content, &creatorID, &createdTime, &image, &like, &dislike)
			if err != nil {
				log.Println(err.Error())
			}

			post = AppendPost(id, title, content, creatorID, image, like, dislike, s.UserID, createdTime)
			postsL = append(postsL, post)
		}
	}
	//create post current user
	pStmp, err := DB.Query("SELECT * FROM posts WHERE creator_id=?", s.UserID)
	//defer psu.Close()
	var postCr Post
	postsCreated := []Post{}

	//todo get uniq post - created post
	for pStmp.Next() {
		err = pStmp.Scan(&id, &title, &content, &creatorID, &createdTime, &image, &like, &dislike)
		if err != nil {
			log.Println(err.Error())
		}
		//post.AuthorForPost = s.UserID
		postCr = AppendPost(id, title, content, creatorID, image, like, dislike, s.UserID, createdTime)
		postsCreated = append(postsCreated, postCr)
	}

	commentQuery, err := DB.Query("SELECT * FROM comments WHERE creator_id=?", s.UserID)

	var comments []Comment
	var cmt Comment
	defer commentQuery.Close()

	for commentQuery.Next() {

		err = commentQuery.Scan(&id, &content, &postID, &userID, &createdTime, &like, &dislike)
		err = DB.QueryRow("SELECT title FROM posts WHERE id = ?", postID).Scan(&title)
		if err != nil {
			log.Println(err.Error())
		}

		err = commentQuery.Scan(&cmt.ID, &cmt.Content, &cmt.PostID, &cmt.UserID, &cmt.CreatedTime, &cmt.Like, &cmt.Dislike)
		comments = append(comments, comment)
	}
	//------------------
	var notifies []Notify
	nQuery, err := DB.Query("SELECT * FROM notify WHERE to_whom=?", s.UserID)

//	write if commentId || voteState == 0 -> from Db
	for nQuery.Next() {
		n := Notify{}
		err = nQuery.Scan(&n.ID, &n.PostID,  &n.UserLostID, &n.voteState, &n.CreatedTime, &n.ToWhom, &n.CommentID)
		if err !=nil {
			log.Println(err)
		}
		notifies = append(notifies, n)
	}
	fmt.Println(notifies, "notidy arr")
	//send client history(list) Likes/Dislikes
	for _, v := range notifies {
		n := Notify{}
		//like/dislike case
		err = DB.QueryRow("SELECT title FROM posts WHERE id = ?", v.PostID).Scan(&n.PostTitle)

		err = DB.QueryRow("SELECT title FROM comments WHERE id = ?", v.CommentID).Scan(&n.CommentTitle)
			//get postTitle, by postID, / get userLost Name, - uid /
			err = DB.QueryRow("SELECT full_name FROM users WHERE id = ?", v.UserLostID).Scan(&n.UserLost)

		if  v.voteState == 1 && v.PostID !=0{
			fmt.Println("user: ", n.UserLost,  " lost liked your post : ",  n.PostTitle, " in ", v.CreatedTime, "")
		}
		if  v.voteState == 2  && v.PostID != 0{
			fmt.Println("user: ", n.UserLost,  " lost Dislike your post : ",  n.PostTitle, " in ", v.CreatedTime, "")
		}

		if  v.voteState == 1 && v.CommentID !=0{
			fmt.Println("user: ", n.UserLost,  " lost liked your COmment : ",  n.CommentTitle, " in ", v.CreatedTime, "")
		}
		if  v.voteState == 2  && v.CommentID != 0{
			fmt.Println("user: ", n.UserLost,  " lost Dislike your Comment : ",  n.CommentTitle, " in ", v.CreatedTime, "")
		}
		// if v.CommentID > 0 {
		// }
	}
	//lops@mail.com
	// if to_whom == currUser_id ? -> notify_table -> send Client, liked/dislieked post || comment -> UserID
	//check if CommentID == nil && voteState != nil 2 case if CommentID !=nil && voteState == nil -> show Comment Post User
	//--------------------
	if err != nil {
		return nil, nil, nil, u, err
	}

	return postsL, postsCreated, comments, u, nil
}

//GetAnotherProfile other user data
func (user *User) GetAnotherProfile(r *http.Request) ([]Post, User, error) {

	//userQR := DB.QueryRow("SELECT * FROM users WHERE id = ?", user.Temp)

	u := User{}
	postsU := []Post{}

	//err = userQR.Scan(&u.ID, &u.FullName, &u.Email, &u.Password, &u.IsAdmin, &u.Age, &u.Sex, &u.CreatedTime, &u.City, &u.Image)
	err = DB.QueryRow("SELECT id, full_name, email, isAdmin, age, sex, created_time, city, image  FROM users WHERE id = ?", user.Temp).Scan(&u.ID, &u.FullName, &u.Email, &u.IsAdmin, &u.Age, &u.Sex, &u.CreatedTime, &u.City, &u.Image)
	if u.Image[0] == 60 {
		u.SVG = true
	}
	u.ImageHTML = base64.StdEncoding.EncodeToString(u.Image)
	psu, err := DB.Query("SELECT * FROM posts WHERE creator_id=?", u.ID)

	defer psu.Close()

	for psu.Next() {
		err = psu.Scan(&id, &title, &content, &creatorID, &createdTime, &image, &like, &dislike)

		if err != nil {
			log.Println(err.Error())
		}
		post = AppendPost(id, title, content, creatorID, image, like, dislike, 0, createdTime)
		postsU = append(postsU, post)
	}
	if err != nil {
		return nil, u, err
	}
	return postsU, u, nil
}

//UpdateProfile function
func (u *User) UpdateProfile() error {

	_, err := DB.Exec("UPDATE  users SET full_name=?, age=?, sex=?, city=?, image=? WHERE id =?",
		u.FullName, u.Age, u.Sex, u.City, u.Image, u.ID)
	if err != nil {
		return err
	}
	return nil
}

//DeleteAccount then dlogut - delete cookie, delete lsot comment, session Db, voteState
func (u *User) DeleteAccount(w http.ResponseWriter, r *http.Request) {

	_, err = DB.Exec("DELETE FROM  session  WHERE user_id=?", u.ID)
	_, err = DB.Exec("DELETE FROM  voteState  WHERE user_id=?", u.ID)
	_, err = DB.Exec("DELETE FROM  comments  WHERE creator_id=?", u.ID)
	_, err = DB.Exec("DELETE FROM  users  WHERE id=?", u.ID)

	if err != nil {
		log.Println(err)
		return
	}

	util.DeleteCookie(w)
}
