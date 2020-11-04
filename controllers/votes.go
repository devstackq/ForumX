package controllers

import (
	"fmt"
	"net/http"

	"github.com/devstackq/ForumX/models"
	util "github.com/devstackq/ForumX/utils"
)

//LostVotes func Post
func LostVotes(w http.ResponseWriter, r *http.Request) {

	// table user - like, dislike -> default - 0, 1,-1,
	//к 4 посту - 13 юзер, Like -> update likeColumnt = 1,
	// к 4 посту 13 юзер поставил dislike -> update Disl1ikeColumn = -1
	//check Logic backend Like -> if disCol == 0 && likeCol ==0, { update LikeCount +1} else if { disCol == -1 && likeCol ==0,} -> disCOl = 0, likeCol = 1, update LikeCount + 1,
	// else if {disCol == 0 && likeCol ==1} -> likeCol = 0, update LikeCount - 1,

	// likeBrdige - postid, like_state, dislike_state

	if util.URLChecker(w, r, "/votes") {

		access, s := util.IsCookie(w, r)
		if !access {
			http.Redirect(w, r, "/signin", 200)
			return
		}

		DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).Scan(&s.UserID)

		pid := r.URL.Query().Get("id")
		lukas := r.FormValue("like")
		diskus := r.FormValue("dislike")

		vote := models.Votes{}
		if r.Method == "POST" {

			if lukas == "1" {
				//if not row
				DB.QueryRow("SELECT id FROM voteState where post_id=?", pid).Scan(&vote.ID)

				if vote.ID == 0 {
					fmt.Print("init")
					_, err = DB.Exec("INSERT INTO voteState(post_id, user_id, like_state) VALUES( ?, ?, ?)", pid, s.UserID, 1)
					if err != nil {
						panic(err)
					}
				} else {

					err = DB.QueryRow("SELECT count_like FROM posts WHERE id=?", pid).Scan(&vote.OldLike)

					err = DB.QueryRow("SELECT count_dislike FROM posts WHERE id=?", pid).Scan(&vote.OldDislike)
continue here, 
					fmt.Println("case 2", vote.OldDislike)
					DB.QueryRow("SELECT like_state, dislike_state FROM voteState where post_id=? and user_id=?", pid, s.UserID).Scan(&vote.Like, &vote.Dislike)
					//set like
					if vote.Like == 1 && vote.Dislike == 0 {
						_, err = DB.Exec("UPDATE posts SET count_like = ? WHERE id= ?", vote.OldLike-1, pid)
						_, err = DB.Exec("UPDATE voteState SET like_state = ? WHERE post_id=? and user_id", 0, pid, s.UserID)
					}
					//set dislike -> to like
					if vote.Like == 0 && vote.Dislike == 1 {
						vote.OldDislike--
						vote.OldLike++
						_, err = DB.Exec("UPDATE posts SET count_dislike = ? WHERE id=?", vote.OldDislike, pid)
						_, err = DB.Exec("UPDATE posts SET count_like = ? WHERE id=?", vote.OldLike, pid)
						_, err = DB.Exec("UPDATE voteState SET like_state = ?, dislike_state=? WHERE post_id=? and user_id", 1, 0, pid, s.UserID)
					}

					//	_, err = DB.Exec("UPDATE  posts SET count_like = ? WHERE id= ?", nv, pid)
					if err != nil {
						panic(err)
					}
				}

				//init -> post -> set like or dislike

				// DB.QueryRow("SELECT like_state, dislike_state FROM voteState").Scan(&vote.Like, &vote.Dislike)
				// if vote.Like == 0 && vote.Dislike == 0 {

				// }
				//DB.QueryRow("SELECT like_state, dislike_state FROM voteState where post_id=?", pid).Scan(&vote.Like, &vote.Dislike)
				// err = DB.QueryRow("SELECT like_state, dislike_state FROM voteState where post_id=?", pid).Scan(&vote.Like, &vote.Dislike)
				// if err != nil {
				// 	log.Println(err)
				// }
				//check if not have post and user lost vote this post
				//1 like or 1 dislike 1 user lost 1 post, get previus value and +1

				// var p, u int
				// err = DB.QueryRow("SELECT post_id, user_id FROM likes WHERE post_id=? AND user_id=?", pid, s.UserID).Scan(&p, &u)

				// if p == 0 && u == 0 {
				// 	oldlike := 0
				// 	err = DB.QueryRow("SELECT count_like FROM posts WHERE id=?", pid).Scan(&oldlike)

				// 	nv := oldlike + 1
				// 	_, err = DB.Exec("UPDATE  posts SET count_like = ? WHERE id= ?", nv, pid)
				// 	if err != nil {
				// 		panic(err)
				// 	}

				// 	_, err = DB.Exec("INSERT INTO likes(post_id, user_id, state_id) VALUES( ?, ?, ?)", pid, s.UserID, 1)
				// 	if err != nil {
				// 		panic(err)
				// 	}
				// }
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
					_, err = DB.Exec("INSERT INTO likes(post_id, user_id, state_id) VALUES( ?, ?, ?)", pid, s.UserID, 0)

					if err != nil {
						panic(err)
					}
				}
			}
		}
		http.Redirect(w, r, "post?id="+pid, 302)
	}
}

func LostVotesComment(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/votes/comment") {

		access, s := util.IsCookie(w, r)
		if !access {
			http.Redirect(w, r, "/signin", 200)
			return
		}

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
			http.Redirect(w, r, "/post?id="+pidc, 200)
		}
	}
}

//Likes table, filed posrid, userid, state_id
// 0,1,2 if state ==0, 1 || 2,
// next btn, if 1 == 1, state =0
