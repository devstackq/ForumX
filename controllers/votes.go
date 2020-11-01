package controllers

import (
	"net/http"

	util "github.com/devstackq/ForumX/utils"
)

//like dislike post
func LostVotes(w http.ResponseWriter, r *http.Request) {

	//flag if like, get llike count + 1, click like still -> get count like, and minus 1,  - toggler,
	//flag dislike, click1 - get count dislike, +1, click2 count dislike, minus 1,

	//flag 3 click like, then click dislike, likeFlag === dislikeFlag, (true), 5 2, 6 2, 5 , 3/ 6 ,2

	//likeToggle, 1,0, true, false
	// L, gL - cV + 1, if likeToggle == true -> +1, else -> get currLike - 1
	// DislikeToggle, true , -> +1, else -> get currDislike - 1

	if util.URLChecker(w, r, "/votes") {

		access, s := util.IsCookie(w, r)
		if !access {
			http.Redirect(w, r, "/signin", 302)
			return
		}

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

					_, err = DB.Exec("INSERT INTO likes(post_id, user_id, state_id) VALUES( ?, ?, ?)", pid, s.UserID, 1)
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
					_, err = DB.Exec("INSERT INTO likes(post_id, user_id, state_id) VALUES( ?, ?, ?)", pid, s.UserID, 0)

					if err != nil {
						panic(err)
					}
				}
			}
		}
		http.Redirect(w, r, "post?id="+pid, 301)
	}
}

func LostVotesComment(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/votes/comment") {

		access, s := util.IsCookie(w, r)
		if !access {
			http.Redirect(w, r, "/signin", 302)
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
			http.Redirect(w, r, "/post?id="+pidc, 301)
		}
	}
}

//Likes table, filed posrid, userid, state_id
// 0,1,2 if state ==0, 1 || 2,
// next btn, if 1 == 1, state =0
