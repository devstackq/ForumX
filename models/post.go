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
	post                         Posts
	comment                      Comment
	msg                          = structure.API.Message
)

type Posts struct {
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
	CountLike     int
	CountDislike  int
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

type PostCategory struct {
	PostID   int64
	Category string
}

type Filter struct {
	Category string
	Like     string
	Date     string
}

//func GetAllPost(r *http.Request, like, date, category string) ([]Posts, string, string, error) {
func (f *Filter) GetAllPost(r *http.Request) ([]Posts, string, string, error) {

	var post Posts
	//send from controlle, then check-> then send model

	var leftJoin bool
	var arrPosts []Posts

	switch r.URL.Path {
	//check what come client, cats, and filter by like, date and cats
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
		post.Endpoint = "/science"
		post.Temp = "Science"
		rows, err = DB.Query("SELECT * FROM posts  LEFT JOIN post_cat_bridge  ON post_cat_bridge.post_id = posts.id   WHERE category=?  ORDER  BY created_time  DESC LIMIT 4", "science")
	case "/love":
		post.Temp = "Love"
		leftJoin = true
		post.Endpoint = "/love"
		rows, err = DB.Query("SELECT  * FROM posts  LEFT JOIN post_cat_bridge  ON post_cat_bridge.post_id = posts.id  WHERE category=?   ORDER  BY created_time  DESC LIMIT 4", "love")
	case "/sapid":
		post.Temp = "Sapid"
		leftJoin = true
		post.Endpoint = "/sapid"
		rows, err = DB.Query("SELECT  * FROM posts  LEFT JOIN post_cat_bridge  ON post_cat_bridge.post_id = posts.id  WHERE category=?  ORDER  BY created_time  DESC LIMIT 4", "sapid")
	}

	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	for rows.Next() {
		post := Posts{}
		if leftJoin {
			if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatorID, &post.CreatedTime, &post.Image, &post.CountLike, &post.CountDislike, &post.PBGID, &post.PBGPostID, &post.PBGCategory); err != nil {
				fmt.Println(err)
			}
		} else {
			if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatorID, &post.CreatedTime, &post.Image, &post.CountLike, &post.CountDislike); err != nil {
				fmt.Println(err)
			}
		}

		arrPosts = append(arrPosts, post)
	}
	//	fmt.Println(arrayPosts, "osts all")
	return arrPosts, post.Endpoint, post.Temp, nil
}

//update post
func (p *Posts) UpdatePost() error {

	_, err := DB.Exec("UPDATE  posts SET title=?, content=?, image=? WHERE id =?",
		p.Title, p.Content, p.Image, p.ID)

	if err != nil {
		return err
	}
	return nil
}

//delete post
func (p *Posts) DeletePost() error {
	_, err := DB.Exec("DELETE FROM  posts  WHERE id =?", p.ID)
	if err != nil {
		return err
	}
	return nil
}

func (p *Posts) CreatePost(w http.ResponseWriter, r *http.Request) {

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
		//file2, _, err := r.FormFile("uploadfile")
		if err != nil {
			log.Fatal(err)
		}
		fileBytes, err = ioutil.ReadAll(p.FileI)
	} else {
		fmt.Print("file more 20mb")
		//message  client send
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
		//DB.QueryRow("SELECT id FROM posts").Scan(&p.La)
		last, err := db.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
		//return last, nil
		//insert cat_post_bridge value

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
			//loop
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
	} else {
		util.DisplayTemplate(w, "header", util.IsAuth(r))
		util.DisplayTemplate(w, "create", "Empty title or content")
	}
}

//link to COmments struct, then call func(r), return arr comments, post, err
func (post *Posts) GetPostById(r *http.Request) ([]Comment, Posts, error) {

	p := Posts{}
	//take from all post, only post by id, then write p struct Post
	DB.QueryRow("SELECT * FROM posts WHERE id = ?", post.ID).Scan(&p.ID, &p.Title, &p.Content, &p.CreatorID, &p.CreatedTime, &p.Image, &p.CountLike, &p.CountDislike)
	p.CreatedTime.Format(time.RFC1123)
	//write values from tables Likes, and write data table Post fileds like, dislikes
	//[]byte -> encode string, client render img base64
	if len(p.Image) > 0 {
		if p.Image[0] == 60 {
			p.SVG = true
		}
	}

	encodedString := base64.StdEncoding.EncodeToString(p.Image)
	p.ImageHTML = encodedString

	//creator post
	DB.QueryRow("SELECT full_name FROM users WHERE id = ?", p.CreatorID).Scan(&p.FullName)

	//get all comments from post1
	stmp, err := DB.Query("SELECT * FROM comments WHERE  post_id =?", p.ID)
	if err != nil {
		log.Fatal(err)
	}
	defer stmp.Close()
	//write each fileds inside Comment struct -> then  append Array Comments
	var comments []Comment

	for stmp.Next() {

		comment := Comment{}

		err = stmp.Scan(&id, &content, &postID, &userID, &createdTime, &like, &dislike)
		if err != nil {
			panic(err.Error)
		}

		comment = Comment{
			ID:          id,
			Content:     content,
			PostID:      postID,
			UserID:      userID,
			CreatedTime: createdTime,
			Like:        like,
			Dislike:     dislike,
		}
		//comment = util.AppendComment(id, content, postID, userID, createdTime, like, dislike)
		DB.QueryRow("SELECT full_name FROM users WHERE id = ?", userID).Scan(&comment.Author)
		comments = append(comments, comment)
	}

	if err != nil {
		return nil, p, err
	}
	return comments, p, nil
}

func (pcb *PostCategory) CreateBridge() error {
	_, err := DB.Exec("INSERT INTO post_cat_bridge (post_id, category) VALUES (?, ?)",
		pcb.PostID, pcb.Category)
	if err != nil {
		return err
	}
	return nil
}

//search post by contain title
func Search(w http.ResponseWriter, r *http.Request) ([]Posts, error) {

	psu, err := DB.Query("SELECT * FROM posts WHERE title LIKE ?", "%"+r.FormValue("search")+"%")
	defer psu.Close()
	var posts []Posts

	for psu.Next() {

		err = psu.Scan(&id, &title, &content, &creatorID, &createdTime, &image, &like, &dislike)
		if err != nil {
			log.Println(err.Error())
		}
		post = appendPost(id, title, content, creatorID, image, like, dislike, 0, createdTime)
		posts = append(posts, post)
	}

	if err != nil {
		return nil, err
	}
	return posts, nil
}

//appendPost each post put value from Db
func appendPost(id int, title, content string, creatorID int, image []byte, like, dislike, authorID int, createdTime time.Time) Posts {

	post = Posts{
		ID:            id,
		Title:         title,
		Content:       content,
		CreatorID:     creatorID,
		Image:         image,
		CountLike:     like,
		CountDislike:  dislike,
		AuthorForPost: authorID,
		CreatedTime:   createdTime,
	}
	return post
}
