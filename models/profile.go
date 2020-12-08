package models

import (
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

//Notify struct
type Notify struct {
	UID          int
	PID          int
	CID          int
	CLID         int
	ID           int
	CIDPID       int
	PostID       int
	CommentID    int
	UserLostID   int
	VoteState    int
	CreatedTime  string
	ToWhom       int
	PostTitle    string
	UserLost     string
	CommentTitle string
}

//GetUserProfile function
func GetUserProfile(r *http.Request, w http.ResponseWriter, cookie *http.Cookie) ([]Post, []Post, []Post, []Comment, User, error) {

	//time.AfterFunc(10, checkCookieLife(cookie, w, r)) try check every 30 min cookie
	s := structure.Session{UUID: cookie.Value}
	u := User{}
	DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).Scan(&s.UserID)

	liked := VotedPosts("like_state", s.UserID)
	disliked := VotedPosts("dislike_state", s.UserID)

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

	//get posts current user
	pStmp, err := DB.Query("SELECT * FROM posts WHERE creator_id=?", s.UserID)
	postsCreated := []Post{}

	for pStmp.Next() {
		err = pStmp.Scan(&post.ID, &post.Title, &post.Content, &post.CreatorID, &post.CreatedTime, &post.Image, &post.Like, &post.Dislike)
		
		if err != nil {
			log.Println(err.Error())
		}
		post.AuthorForPost = s.UserID
		post.Time = post.CreatedTime.Format("2006 Jan _2 15:04:05")
		postsCreated = append(postsCreated, post)
	}

	commentQuery, err := DB.Query("SELECT * FROM comments WHERE creator_id=?", s.UserID)

	var comments []Comment
	var cmt Comment
	defer commentQuery.Close()

	for commentQuery.Next() {

		err = commentQuery.Scan(&cmt.ID, &cmt.Content, &cmt.PostID, &cmt.UserID, &cmt.Time, &cmt.Like, &cmt.Dislike)
		if err != nil {
			log.Println(err.Error())
		}
		err = DB.QueryRow("SELECT post_id FROM comments WHERE id = ?", cmt.ID).Scan(&postID)
		if err != nil {
			log.Println(err.Error())
		}
		err = DB.QueryRow("SELECT title FROM posts WHERE id = ?", postID).Scan(&cmt.TitlePost)
		if err != nil {
			log.Println(err.Error())
		}

		cmt.CreatedTime = cmt.Time.Format("2006 Jan _2 15:04:05")
		comments = append(comments, cmt)
	}
	// if err != nil {
	// 	return nil, nil, nil, nil, u, err
	// }
	return disliked, liked, postsCreated, comments, u, nil
}

//GetUserActivities func
func GetUserActivities(w http.ResponseWriter, r *http.Request) (result []Notify) {

	cookie, _ := r.Cookie("_cookie")
	s := structure.Session{UUID: cookie.Value}
	DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).Scan(&s.UserID)

	var notifies []Notify
	nQuery, err := DB.Query("SELECT * FROM notify WHERE to_whom=?", s.UserID)

	for nQuery.Next() {
		n := Notify{}
		err = nQuery.Scan(&n.ID, &n.PostID, &n.UserLostID, &n.VoteState, &n.CreatedTime, &n.ToWhom, &n.CommentID)
		if err != nil {
			log.Println(err)
		}
		notifies = append(notifies, n)
	}

	for _, v := range notifies {
		//get postTitle, by postID, / get userLost Name, - uid /
		n := Notify{}
		//like/dislike case
		// commneId - delete, but notify - Have row
		err = DB.QueryRow("SELECT title FROM posts WHERE id = ?", v.PostID).Scan(&n.PostTitle)
		if err != nil {
			log.Println(err)
		}
		err = DB.QueryRow("SELECT post_id FROM comments WHERE id = ?", v.CommentID).Scan(&n.CIDPID)
		if err != nil {
			log.Println(err)
		}

		err = DB.QueryRow("SELECT content FROM comments WHERE id = ?", v.CommentID).Scan(&n.CommentTitle)
		if err != nil {
			log.Println(err)
		}
		err = DB.QueryRow("SELECT full_name FROM users WHERE id = ?", v.UserLostID).Scan(&n.UserLost)
		if err != nil {
			log.Println(err)
		}

		n.VoteState = v.VoteState
		n.UID = v.UserLostID

		if v.VoteState == 1 && v.PostID != 0 {
			n.PID = v.PostID
			fmt.Println("user: ", n.UserLost, " lost liked your post : ", n.PostTitle, " in ", v.CreatedTime, "")
		}
		if v.VoteState == 2 && v.PostID != 0 {
			n.PID = v.PostID
			fmt.Println("user: ", n.UserLost, " lost Dislike your post : ", n.PostTitle, " in ", v.CreatedTime, "")
		}
		if v.VoteState == 1 && v.CommentID != 0 {
			// if n.CommentTitle == "" {
			// 	return
			// }
			n.CID = v.CommentID
			n.PostTitle = n.CommentTitle
			fmt.Println("user: ", n.UserLost, " lost liked your Comment : ", n.CommentTitle, " in ", v.CreatedTime, "")
		}
		if v.VoteState == 2 && v.CommentID != 0 {
			n.CID = v.CommentID
			n.PostTitle = n.CommentTitle
			fmt.Println("user: ", n.UserLost, " lost Dislike your Comment!!!: ", n.CommentTitle, " in ", v.CreatedTime, "", n.CID, n.CIDPID)
		}
		//comment lost case
		if v.VoteState == 0 && v.CommentID != 0 {
			fmt.Println("user: ", n.UserLost, " lost Comment u Post: ", n.CommentTitle, " in ", v.CreatedTime)
			n.CLID = v.PostID
			n.PostTitle = n.CommentTitle
		}
		result = append(result, n)
	}
	return result
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

		err = psu.Scan(&post.ID, &post.Title, &post.Content, &post.CreatorID, &post.CreatedTime, &post.Image, &post.Like, &post.Dislike)
		if err != nil {
			log.Println(err.Error())
		}
		//AuthorForPost
		post.Time = post.CreatedTime.Format("2006 Jan _2 15:04:05")
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

func VotedPosts(voteType string, uid int) (result []Post) {

	postArr := []Votes{}
	arrIDVote := []int{}

	votedPost, err := DB.Query("select post_id from voteState where user_id=? and  "+voteType+" and comment_id is null", uid, 1)
	if err != nil {
		log.Println(err)
	}
	for votedPost.Next() {
		voteLiked := Votes{}
		err = votedPost.Scan(&voteLiked.PostID)
		postArr = append(postArr, voteLiked)
	}
	defer votedPost.Close()

	for _, v := range postArr {
		arrIDVote = append(arrIDVote, v.PostID)
	}

	for _, v := range arrIDVote {
		smtp, err := DB.Query("SELECT * FROM posts WHERE id=?", v)
		if err != nil {
			log.Println(err)
		}
		p := Post{}
		for smtp.Next() {
			err = smtp.Scan(&p.ID, &p.Title, &p.Content, &p.CreatorID, &p.CreatedTime, &p.Image, &p.Like, &p.Dislike)
			if err != nil {
				log.Println(err.Error())
			}
			p.Time = p.CreatedTime.Format("2006 Jan _2 15:04:05")
			result = append(result, p)
		}
	}
	return result
}
