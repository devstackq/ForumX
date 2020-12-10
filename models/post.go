package models

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	structure "github.com/devstackq/ForumX/general"
	util "github.com/devstackq/ForumX/utils"
)

//global variable for package models
var (
	err                          error
	DB                           *sql.DB
	rows                         *sql.Rows
	post                         Post
	comment                      Comment
	msg                          = structure.API.Message
	pageNum                      = 1
)

//Posts struct
type Post struct {
	ID            int `json:"id"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	CreatorID     int `json:"creatorId"`
	CreatedTime   time.Time `json:"createdTime"`
	Endpoint      string `json:"endpoint"`
	FullName      string `json:"fullName"`
	Image         []byte `json:"image"`
	ImageHTML     string `json:"imageHtml"`
	PostIDEdit    int `json:"postEditId"`
	AuthorForPost int `json:"authorPost"`
	Like          int `json:"like"`
	Dislike       int `json:"dislike"`
	SVG           bool `json:"svg"`
	PBGID         int `json:"pbgId"`
	PBGPostID     int `json:"pbgPostId"`
	PBGCategory   string `json:"pbgCategory"`
	FileS         multipart.File `json:"fileS"`
	FileI         multipart.File `json:"fileB"`
	Session       structure.Session `json:"session"`
	Categories    []string `json:"categories"`
	Temp          string `json:"temp"`
	IsPhoto       bool `json:"isPhoto"`
	Time          string `json:"time"`
	CountPost     int `json:"countPost"`
}

//PostCategory struct
type PostCategory struct {
	PostID   int64
	Category string
}

//Filter struct
type Filter struct {
	Category string `json:"cateogry"`
	Like     string `json:"like"`
	Date     string `json:"date"`
}

//GetAllPost function
func (f *Filter) GetAllPost(r *http.Request, next, prev string) ([]Post, string, string, error) {
	//pageNum = 1
	var post Post
	var leftJoin bool
	var arrPosts []Post

	//each call +1
	if next == "next" {
		pageNum++
	}
	if prev == "prev" {
		pageNum--
	}
	//count pageNum, fix

	limit := 4
	offset := limit * (pageNum - 1)

	switch r.URL.Path {
	case "/":
		leftJoin = false
		post.Endpoint = "/"
		if f.Date == "asc" {
			rows, err = DB.Query("SELECT * FROM posts  ORDER BY created_time ASC LIMIT 8 ")
		} else if f.Date == "desc" {
			rows, err = DB.Query("SELECT * FROM posts  ORDER BY created_time DESC LIMIT 8")
		} else if f.Like == "like" {
			rows, err = DB.Query("SELECT * FROM posts  ORDER BY count_like DESC LIMIT 8")
		} else if f.Like == "dislike" {
			rows, err = DB.Query("SELECT * FROM posts  ORDER BY count_dislike DESC LIMIT 8")
		} else if f.Category != "" {
			leftJoin = true
			rows, err = DB.Query("SELECT  * FROM posts  LEFT JOIN post_cat_bridge  ON post_cat_bridge.post_id = posts.id   WHERE category=? ORDER  BY created_time  DESC LIMIT 8", f.Category)
		} else {
			rows, err = DB.Query("SELECT * FROM posts ORDER BY created_time DESC LIMIT ? OFFSET ?", limit, offset)
		}

	case "/science":
		leftJoin = true
		post.Temp = "Science"
		post.Endpoint = "/science"
		rows, err = DB.Query("SELECT * FROM posts  LEFT JOIN post_cat_bridge  ON post_cat_bridge.post_id = posts.id   WHERE category=?  ORDER  BY created_time  DESC LIMIT 5", "science")
	case "/love":
		leftJoin = true
		post.Temp = "Love"
		post.Endpoint = "/love"
		rows, err = DB.Query("SELECT  * FROM posts  LEFT JOIN post_cat_bridge  ON post_cat_bridge.post_id = posts.id  WHERE category=?   ORDER  BY created_time  DESC LIMIT 5", "love")
	case "/sapid":
		leftJoin = true
		post.Temp = "Sapid"
		post.Endpoint = "/sapid"
		rows, err = DB.Query("SELECT  * FROM posts  LEFT JOIN post_cat_bridge  ON post_cat_bridge.post_id = posts.id  WHERE category=?  ORDER  BY created_time  DESC LIMIT 5", "sapid")
	}

	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	for rows.Next() {
		if leftJoin {
			if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatorID, &post.CreatedTime, &post.Image, &post.Like, &post.Dislike, &post.PBGID, &post.PBGPostID, &post.PBGCategory); err != nil {
				fmt.Println(err)
			}
		} else {
			if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatorID, &post.CreatedTime, &post.Image, &post.Like, &post.Dislike); err != nil {
				fmt.Println(err)
			}
			//fmt.Print(post.ID)
		}

		if err != nil {
			log.Println(err)
		}
		//send countr +1
		err = DB.QueryRow("SELECT COUNT(id) FROM posts").Scan(&post.CountPost)
		post.Time = post.CreatedTime.Format("2006 Jan _2 15:04:05")
		arrPosts = append(arrPosts, post)
	}
	//err = DB.QueryRow("SELECT COUNT(id) FROM posts").Scan(&post.CountPost)
	return arrPosts, post.Endpoint, post.Temp, nil
}

//UpdatePost fucntion
func (p *Post) UpdatePost() error {

	_, err := DB.Exec("UPDATE  posts SET title=?, content=?, image=? WHERE id =?",
		p.Title, p.Content, p.Image, p.ID)

	if err != nil {
		return err
	}
	return nil
}

//DeletePost function, delete rows, notify, voteState, comment, by postId
func (p *Post) DeletePost() error {

	_, err = DB.Exec("DELETE FROM posts  WHERE id =?", p.ID)
	if err != nil {
		return err
	}
	_, err = DB.Exec("DELETE FROM comments  WHERE post_id =?", p.ID)
	if err != nil {
		return err
	}
	_, err = DB.Exec("DELETE FROM notify  WHERE post_id =?", p.ID)
	if err != nil {
		return err
	}
	_, err = DB.Exec("DELETE FROM voteState  WHERE post_id =?", p.ID)
	if err != nil {
		return err
	}
	_, err = DB.Exec("DELETE FROM post_cat_bridge  WHERE post_id =?", p.ID)
	if err != nil {
		return err
	}

	return nil
}

//CreatePost function
func (p *Post) CreatePost(w http.ResponseWriter, r *http.Request) {

	var fileBytes []byte
	var buff bytes.Buffer

	if p.IsPhoto {
		fileSize, _ := buff.ReadFrom(p.FileS)
		defer p.FileS.Close()

		if fileSize < 20000000 {
			if err != nil {
				log.Fatal(err)
			}
			fileBytes, err = ioutil.ReadAll(p.FileI)
		} else {
			util.DisplayTemplate(w, "header", util.IsAuth(r))
			util.DisplayTemplate(w, "create", "Large file, more than 20mb")
		}
	} else {
		//set empty photo post
		fileBytes = []byte{0, 0}
	}

	DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", p.Session.UUID).Scan(&p.Session.UserID)

	//check empty values
	if util.CheckLetter(p.Title) && util.CheckLetter(p.Content) {

		createPostPrepare, err := DB.Prepare(`INSERT INTO posts (title, content, creator_id, created_time, image) VALUES ( ?,?, ?, ?, ?)` )
		if err != nil {
			log.Println(err)
		}
		defer createPostPrepare.Close()
		createPostExec, err := createPostPrepare.Exec(p.Title, p.Content, p.Session.UserID, time.Now(), fileBytes) 
		if err != nil {
			log.Println(err)
		}

		last, err := createPostExec.LastInsertId()

		if err != nil {
			log.Println(err)
		}
		pcb := PostCategory{}

		if len(p.Categories) == 0 {

			pcb = PostCategory{
				PostID:   last,
				Category: "sapid",
			}
			pcb.CreateBridge()
		} else if len(p.Categories) == 1 {

			pcb = PostCategory{
				PostID:   last,
				Category: p.Categories[0],
			}
			pcb.CreateBridge()

		} else if len(p.Categories) > 1 {
			//loop add > 1 category post
			for _, v := range p.Categories {
				pcb = PostCategory{
					PostID:   last,
					Category: v,
				}
				pcb.CreateBridge()
			}
		}

		http.Redirect(w, r, "/post?id="+strconv.Itoa(int(last)), 302)

	} else {
		msg = "Empty title or content"
		util.DisplayTemplate(w, "header", util.IsAuth(r))
		util.DisplayTemplate(w, "create_post", &msg)
	}
}

//GetPostById function take from all post, only post by id, then write p struct Post
func (post *Post) GetPostByID(r *http.Request) ([]Comment, Post, error) {

	p := Post{}
	DB.QueryRow("SELECT * FROM posts WHERE id = ?", post.ID).Scan(&p.ID, &p.Title, &p.Content, &p.CreatorID, &p.CreatedTime, &p.Image, &p.Like, &p.Dislike)

	//[]byte -> encode string, client render img base64
	//check svg || jpg,png
	if len(p.Image) > 0 {
		if p.Image[0] == 60 {
			p.SVG = true
		}
	}
	p.Time = p.CreatedTime.Format("2006 Jan _2 15:04:05")

	p.ImageHTML = base64.StdEncoding.EncodeToString(p.Image)

	DB.QueryRow("SELECT full_name FROM users WHERE id = ?", p.CreatorID).Scan(&p.FullName)

	stmp, err := DB.Query("SELECT * FROM comments WHERE  post_id =?", p.ID)
	if err != nil {
		log.Fatal(err)
	}
	defer stmp.Close()
	//write each fields inside Comment struct -> then  append Array Comments
	var comments []Comment

	for stmp.Next() {

		c := Comment{}
		err = stmp.Scan(&c.ID, &c.Content, &c.PostID, &c.UserID, &c.Time, &c.Like, &c.Dislike)
		if err != nil {
			log.Println(err.Error())
		}
		c.CreatedTime = c.Time.Format("2006 Jan _2 15:04:05")

		DB.QueryRow("SELECT full_name FROM users WHERE id = ?", c.UserID).Scan(&c.Author)
		comments = append(comments, c)
	}

	if err != nil {
		return nil, p, err
	}
	return comments, p, nil
}

//CreateBridge create post  -> post_id + category
func (pcb *PostCategory) CreateBridge() {

		createBridgePrepare, err := DB.Prepare(`INSERsT INTO post_cat_bridge (post_id, category) VALUES (?, ?)` )
		if err != nil {
			log.Println(err)
		}
		defer createBridgePrepare.Close()
		_, err = createBridgePrepare.Exec(pcb.PostID, pcb.Category) 
		if err != nil {
			log.Println(err)
		}
}

//Search post by contain title
func Search(w http.ResponseWriter, r *http.Request) ([]Post, error) {

	var posts []Post
	psu, err := DB.Query("SELECT * FROM posts WHERE title LIKE ?", "%"+r.FormValue("search")+"%")
	defer psu.Close()

	for psu.Next() {

		err = psu.Scan(&post.ID, &post.Title, &post.Content, &post.CreatorID, &post.CreatedTime, &post.Image, &post.Like, &post.Dislike)
		post.Time = post.CreatedTime.Format("2006 Jan _2 15:04:05")
		posts = append(posts, post)
	}

	if err != nil {
		return nil, err
	}
	return posts, nil
}
