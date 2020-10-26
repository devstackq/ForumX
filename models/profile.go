package models

import (
	"database/sql"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	structure "github.com/devstackq/ForumX/general"
	util "github.com/devstackq/ForumX/utils"
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
	Temp        string
}

//get profile by id
func GetUserProfile(r *http.Request, w http.ResponseWriter, cookie *http.Cookie) ([]Posts, []Posts, []Comment, Users, error) {

	//time.AfterFunc(10, checkCookieLife(cookie, w, r))
	s := structure.Session{UUID: cookie.Value}
	u := Users{}
	DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).Scan(&s.UserID)
	lps := []Likes{}

	//count dislike equal 0 - add query
	lp, err := DB.Query("select post_id from likes where user_id =? and state_id =?", s.UserID, 1)
	defer lp.Close()
	for lp.Next() {
		l := Likes{}
		var lpid int
		err = lp.Scan(&lpid)
		l.PostID = lpid
		lps = append(lps, l)
	}

	err = DB.QueryRow("SELECT * FROM users WHERE id = ?", s.UserID).Scan(&u.ID, &u.FullName, &u.Email, &u.Password, &u.IsAdmin, &u.Age, &u.Sex, &u.CreatedTime, &u.City, &u.Image)
	if u.Image[0] == 60 {
		u.SVG = true
	}

	encStr := base64.StdEncoding.EncodeToString(u.Image)
	u.ImageHtml = encStr

	var likedpost *sql.Rows
	postsL := []Posts{}

	var can []int

	for _, v := range lps {
		can = append(can, v.PostID)
	}

	//unique liked post by user
	fin := util.IsUnique(can)
	//accum liked post

	for _, v := range fin {
		//get each liked post by ID, then likedpost, put array post

		//count_dislike не
		likedpost, err = DB.Query("SELECT * FROM posts WHERE id=? and count_like > 0", v)
		if err != nil {
			log.Println(err)
		}
		for likedpost.Next() {
			err = likedpost.Scan(&id, &title, &content, &creatorID, &createdTime, &image, &like, &dislike)
			if err != nil {
				panic(err.Error)
			}
			post = appendPost(id, title, content, creatorID, image, like, dislike, s.UserID, createdTime)
			postsL = append(postsL, post)
		}
	}
	//create post current user
	psu, err := DB.Query("SELECT * FROM posts WHERE creator_id=?", s.UserID)
	//defer psu.Close()
	var postCr Posts
	postsX := []Posts{}

	//todo get uniq post - created post
	for psu.Next() {
		//here
		err = psu.Scan(&id, &title, &content, &creatorID, &createdTime, &image, &like, &dislike)
		if err != nil {
			log.Println(err.Error())
		}
		//post.AuthorForPost = s.UserID

		postCr = appendPost(id, title, content, creatorID, image, like, dislike, s.UserID, createdTime)
		postsX = append(postsX, postCr)
	}

	csu, err := DB.Query("SELECT * FROM comments WHERE user_idx=?", s.UserID)
	var comments []Comment
	defer csu.Close()

	for csu.Next() {

		err = csu.Scan(&id, &content, &postID, &userID, &createdTime, &like, &dislike)
		err = DB.QueryRow("SELECT title FROM posts WHERE id = ?", postID).Scan(&title)
		if err != nil {
			panic(err.Error)
		}

		comment = appendComment(id, content, postID, userID, createdTime, like, dislike, title)
		comments = append(comments, comment)
	}

	if err != nil {
		return nil, nil, nil, u, err
	}

	return postsL, postsX, comments, u, nil
}

//get other user, posts
func (user *Users) GetAnotherProfile(r *http.Request) ([]Posts, Users, error) {

	userQR := DB.QueryRow("SELECT * FROM users WHERE id = ?", user.Temp)

	u := Users{}
	postsU := []Posts{}

	err = userQR.Scan(&u.ID, &u.FullName, &u.Email, &u.Password, &u.IsAdmin, &u.Age, &u.Sex, &u.CreatedTime, &u.City, &u.Image)
	if u.Image[0] == 60 {
		u.SVG = true
	}
	encStr := base64.StdEncoding.EncodeToString(u.Image)
	u.ImageHtml = encStr
	psu, err := DB.Query("SELECT * FROM posts WHERE creator_id=?", u.ID)

	defer psu.Close()

	var image []byte

	for psu.Next() {
		err = psu.Scan(&id, &title, &content, &creatorID, &createdTime, &image, &like, &dislike)

		if err != nil {
			panic(err.Error)
		}
		post = appendPost(id, title, content, creatorID, image, like, dislike, 0, createdTime)
		postsU = append(postsU, post)
	}
	if err != nil {
		return nil, u, err
	}
	return postsU, u, nil
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
