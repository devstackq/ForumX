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
	"time"

	structure "github.com/devstackq/ForumX/general"
	util "github.com/devstackq/ForumX/utils"
)

//global variable for package models
var (
	err                          error
	DB                           *sql.DB
	rows                         *sql.Rows
	id, creatorID, like, dislike int
	content, title               string
	createdTime                  time.Time
	image                        []byte
	postID                       int
	userID                       int
	post                         Post
	comment                      Comment
	msg                          = structure.API.Message
)

//Posts struct
type Post struct {
	ID            int
	Title         string
	Content       string
	CreatorID     int
	CreatedTime   time.Time
	Endpoint      string
	FullName      string
	CategoryName  string
	Image         []byte
	ImageHTML     string
	PostIDEdit    int
	AuthorForPost int
	Like          int
	Dislike       int
	SVG           bool
	PBGID         int
	PBGPostID     int
	PBGCategory   string
	FileS         multipart.File
	FileI         multipart.File
	Session       structure.Session
	Categories    []string
	Temp          string
}

//PostCategory struct
type PostCategory struct {
	PostID   int64
	Category string
}

//Filter struct
type Filter struct {
	Category string
	Like     string
	Date     string
}

//GetAllPost function
func (f *Filter) GetAllPost(r *http.Request) ([]Post, string, string, error) {

	var post Post
	var leftJoin bool
	var arrPosts []Post
	//check what come client, cats, and filter by like, date and cats
	switch r.URL.Path {
	case "/":
		leftJoin = false
		post.Endpoint = "/"
		if f.Date == "asc" {
			rows, err = DB.Query("SELECT * FROM posts  ORDER BY created_time ASC LIMIT 6")
		} else if f.Date == "desc" {
			rows, err = DB.Query("SELECT * FROM posts  ORDER BY created_time DESC LIMIT 6")
		} else if f.Like == "like" {
			rows, err = DB.Query("SELECT * FROM posts  ORDER BY count_like DESC LIMIT 6")
		} else if f.Like == "dislike" {
			rows, err = DB.Query("SELECT * FROM posts  ORDER BY count_dislike DESC LIMIT 6")
		} else if f.Category != "" {
			leftJoin = true
			rows, err = DB.Query("SELECT  * FROM posts  LEFT JOIN post_cat_bridge  ON post_cat_bridge.post_id = posts.id   WHERE category=? ORDER  BY created_time  DESC LIMIT 6", f.Category)
		} else {
			rows, err = DB.Query("SELECT * FROM posts  ORDER BY created_time DESC LIMIT 6")
		}

	case "/science":
		leftJoin = true
		post.Temp = "Science"
		post.Endpoint = "/science"
		rows, err = DB.Query("SELECT * FROM posts  LEFT JOIN post_cat_bridge  ON post_cat_bridge.post_id = posts.id   WHERE category=?  ORDER  BY created_time  DESC LIMIT 4", "science")
	case "/love":
		leftJoin = true
		post.Temp = "Love"
		post.Endpoint = "/love"
		rows, err = DB.Query("SELECT  * FROM posts  LEFT JOIN post_cat_bridge  ON post_cat_bridge.post_id = posts.id  WHERE category=?   ORDER  BY created_time  DESC LIMIT 4", "love")
	case "/sapid":
		leftJoin = true
		post.Temp = "Sapid"
		post.Endpoint = "/sapid"
		rows, err = DB.Query("SELECT  * FROM posts  LEFT JOIN post_cat_bridge  ON post_cat_bridge.post_id = posts.id  WHERE category=?  ORDER  BY created_time  DESC LIMIT 4", "sapid")
	}

	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	for rows.Next() {
		post := Post{}
		if leftJoin {
			if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatorID, &post.CreatedTime, &post.Image, &post.Like, &post.Dislike, &post.PBGID, &post.PBGPostID, &post.PBGCategory); err != nil {
				fmt.Println(err)
			}
		} else {
			if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatorID, &post.CreatedTime, &post.Image, &post.Like, &post.Dislike); err != nil {
				fmt.Println(err)
			}
		}
		arrPosts = append(arrPosts, post)
	}
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

//DeletePost function
func (p *Post) DeletePost() error {
	_, err := DB.Exec("DELETE FROM  posts  WHERE id =?", p.ID)
	if err != nil {
		return err
	}
	return nil
}

//CreatePost function
func (p *Post) CreatePost(w http.ResponseWriter, r *http.Request) {

	//try default photo user or post
	// fImg, err := os.Open("./1553259670.jpg")

	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// defer fImg.Close()

	// imgInfo, err := fImg.Stat()
	// if err != nil {
	// 	fmt.Println(err, "stats")
	// 	os.Exit(1)
	// }

	// var size int64 = imgInfo.Size()
	// fmt.Println(size, "size")
	// byteArr := make([]byte, size)

	// read file into bytes
	// buffer := bufio.NewReader(fImg)
	// _, err = buffer.Read(byteArr)
	//defer fImg.Close()

	var fileBytes []byte
	var buff bytes.Buffer

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

	DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", p.Session.UUID).Scan(&p.Session.UserID)

	//check empty values
	if util.CheckLetter(p.Title) && util.CheckLetter(p.Content) {

		db, err := DB.Exec("INSERT INTO posts (title, content, creator_id,  image) VALUES ( ?,?, ?, ?)",
			p.Title, p.Content, p.Session.UserID, fileBytes)
		if err != nil {
			log.Println(err)
		}
		last, err := db.LastInsertId()
		if err != nil {
			log.Println(err)
		}

		if len(p.Categories) == 1 {
			pcb := PostCategory{
				PostID:   last,
				Category: p.Categories[0],
			}
			err = pcb.CreateBridge()
			if err != nil {
				log.Println(err)
			}
		} else if len(p.Categories) > 1 {
			//loop add > 1 category post
			for _, v := range p.Categories {
				pcb := PostCategory{
					PostID:   last,
					Category: v,
				}
				err = pcb.CreateBridge()
				if err != nil {
					log.Println(err)
				}
			}
		}
		w.WriteHeader(http.StatusCreated)
		http.Redirect(w, r, "/", http.StatusOK)
		util.DisplayTemplate(w, "index", "")
	} else {
		util.DisplayTemplate(w, "header", util.IsAuth(r))
		util.DisplayTemplate(w, "create", "Empty title or content")
	}
}

//GetPostById function take from all post, only post by id, then write p struct Post
func (post *Post) GetPostByID(r *http.Request) ([]Comment, Post, error) {

	p := Post{}
	DB.QueryRow("SELECT * FROM posts WHERE id = ?", post.ID).Scan(&p.ID, &p.Title, &p.Content, &p.CreatorID, &p.CreatedTime, &p.Image, &p.Like, &p.Dislike)
	p.CreatedTime.Format(time.RFC1123)
	//[]byte -> encode string, client render img base64
	//check svg || jpg,png
	if len(p.Image) > 0 {
		if p.Image[0] == 60 {
			p.SVG = true
		}
	}

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

		comment := Comment{}
		err = stmp.Scan(&id, &content, &postID, &userID, &createdTime, &like, &dislike)
		if err != nil {
			panic(err.Error)
		}

		comment = AppendComment(id, content, postID, userID, createdTime, like, dislike, "")
		DB.QueryRow("SELECT full_name FROM users WHERE id = ?", userID).Scan(&comment.Author)
		comments = append(comments, comment)
	}

	if err != nil {
		return nil, p, err
	}
	return comments, p, nil
}

// /CreateBridge create post  -> post_id + category
func (pcb *PostCategory) CreateBridge() error {

	_, err := DB.Exec("INSERT INTO post_cat_bridge (post_id, category) VALUES (?, ?)",
		pcb.PostID, pcb.Category)
	if err != nil {
		return err
	}
	return nil
}

//Search post by contain title
func Search(w http.ResponseWriter, r *http.Request) ([]Post, error) {

	var posts []Post
	psu, err := DB.Query("SELECT * FROM posts WHERE title LIKE ?", "%"+r.FormValue("search")+"%")
	defer psu.Close()

	for psu.Next() {

		err = psu.Scan(&id, &title, &content, &creatorID, &createdTime, &image, &like, &dislike)
		if err != nil {
			log.Println(err.Error())
		}
		post = AppendPost(id, title, content, creatorID, image, like, dislike, 0, createdTime)
		posts = append(posts, post)
	}

	if err != nil {
		return nil, err
	}
	return posts, nil
}

//appendPost each post put value from Db
func AppendPost(id int, title, content string, creatorID int, image []byte, like, dislike, authorID int, createdTime time.Time) Post {

	post = Post{
		ID:            id,
		Title:         title,
		Content:       content,
		CreatorID:     creatorID,
		Image:         image,
		Like:          like,
		Dislike:       dislike,
		AuthorForPost: authorID,
		CreatedTime:   createdTime,
	}
	return post
}