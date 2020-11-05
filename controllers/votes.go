package controllers

import (
	"fmt"
	"net/http"

	"github.com/devstackq/ForumX/models"
	util "github.com/devstackq/ForumX/utils"
)

//LostVotes func Post
func LostVotes(w http.ResponseWriter, r *http.Request) {

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
				//2 case - uniq row - inuq user - 1 post
				DB.QueryRow("SELECT id FROM voteState where post_id=? and user_id=?", pid).Scan(&vote.ID, s.UserID)
				fmt.Println(vote.ID, s.UserID, "her")
				if vote.ID == 0 {
					fmt.Println("init like")

					_, err = DB.Exec("INSERT INTO voteState(post_id, user_id, like_state, dislike_state) VALUES( ?, ?, ?,?)", pid, s.UserID, 1, 0)
					_, err = DB.Exec("UPDATE posts SET count_like=? WHERE id=?", 1, pid)

					if err != nil {
						panic(err)
					}
				} else {

					err = DB.QueryRow("SELECT count_like FROM posts WHERE id=?", pid).Scan(&vote.OldLike)
					err = DB.QueryRow("SELECT count_dislike FROM posts WHERE id=?", pid).Scan(&vote.OldDislike)

					DB.QueryRow("SELECT like_state, dislike_state FROM voteState where post_id=? and user_id=?", pid, s.UserID).Scan(&vote.LikeState, &vote.DislikeState)

					//continue here,
					fmt.Println(" old Dislike & like", vote.OldDislike, vote.OldLike)
					//set like
					if vote.LikeState == 1 && vote.DislikeState == 0 {
						fmt.Println("case 2 like 1, dis 0")

						vote.OldLike--
						_, err = DB.Exec("UPDATE posts SET count_like = ? WHERE id= ?", vote.OldLike, pid)
						_, err = DB.Exec("UPDATE voteState SET like_state = ? WHERE post_id=? and user_id", 0, pid, s.UserID)
					}
					//set dislike -> to like
					if vote.LikeState == 0 && vote.DislikeState == 1 {
						fmt.Println("case 3 like 0, dis 1")

						vote.OldDislike--
						vote.OldLike++
						_, err = DB.Exec("UPDATE posts SET count_dislike = ?, count_like=? WHERE id=?", vote.OldDislike, vote.OldLike, pid)
						_, err = DB.Exec("UPDATE voteState SET like_state = ?, dislike_state=? WHERE post_id=? and user_id", 1, 0, pid, s.UserID)
					}

					if vote.LikeState == 0 && vote.DislikeState == 0 {
						fmt.Println("case 4 like 0, dis 0, Ls 1, DS 0")

						vote.OldLike++
						_, err = DB.Exec("UPDATE posts SET count_like=? WHERE id=?", vote.OldLike, pid)
						_, err = DB.Exec("UPDATE voteState SET like_state = ?, dislike_state=? WHERE post_id=? and user_id", 1, 0, pid, s.UserID)
					}

					if err != nil {
						panic(err)
					}
				}
			}

			if diskus == "1" {

				DB.QueryRow("SELECT id FROM voteState WHERE post_id=? AND user_id=?", pid).Scan(&vote.ID, s.UserID)
				fmt.Println(vote.ID, s.UserID, "her1")
				if vote.ID == 0 {
					fmt.Println("init dislike")
					_, err = DB.Exec("INSERT INTO voteState(post_id, user_id, dislike_state, like_state) VALUES( ?, ?, ?,?)", pid, s.UserID, 1, 0)
					_, err = DB.Exec("UPDATE posts SET count_dislike=? WHERE id=?", 1, pid)

					//_, err = DB.Exec("UPDATE voteState SET dislike_state = ? WHERE post_id=? and user_id", 1, pid, s.UserID)
					if err != nil {
						panic(err)
					}

				} else {
					err = DB.QueryRow("SELECT count_like FROM posts WHERE id=?", pid).Scan(&vote.OldLike)
					err = DB.QueryRow("SELECT count_dislike FROM posts WHERE id=?", pid).Scan(&vote.OldDislike)

					DB.QueryRow("SELECT like_state, dislike_state FROM voteState where post_id=? and user_id=?", pid, s.UserID).Scan(&vote.LikeState, &vote.DislikeState)

					//set dislike
					if vote.LikeState == 0 && vote.DislikeState == 1 {

						vote.OldDislike--
						_, err = DB.Exec("UPDATE posts SET count_dislike = ? WHERE id=?", vote.OldDislike, pid)
						_, err = DB.Exec("UPDATE voteState SET  dislike_state=? WHERE post_id=? and user_id", 0, pid, s.UserID)

						fmt.Println("case 2 like 0, dis 1")
					}

					//set dislike -> to like
					if vote.LikeState == 1 && vote.DislikeState == 0 {
						fmt.Println("case 3 like 1, dis 0")

						vote.OldDislike++
						vote.OldLike--
						_, err = DB.Exec("UPDATE posts SET count_dislike = ? WHERE id=?", vote.OldDislike, pid)
						_, err = DB.Exec("UPDATE posts SET count_like = ? WHERE id=?", vote.OldLike, pid)
						_, err = DB.Exec("UPDATE voteState SET like_state = ?, dislike_state=? WHERE post_id=? and user_id", 0, 1, pid, s.UserID)
					}

					if vote.LikeState == 0 && vote.DislikeState == 0 {
						fmt.Println("case 4 like 0, dis 0 LS 0, DS 1")
						vote.OldDislike++
						_, err = DB.Exec("UPDATE posts SET count_dislike=? WHERE id=?", vote.OldDislike, pid)
						_, err = DB.Exec("UPDATE voteState SET dislike_state = ?, like_state=? WHERE post_id=? and user_id", 1, 0, pid, s.UserID)
					}

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
