package controllers

import (
	"fmt"
	"log"
	"net/http"

	structure "github.com/devstackq/ForumX/general"
	"github.com/devstackq/ForumX/models"
	util "github.com/devstackq/ForumX/utils"
)

//LostVotes func Post
func VotesPost(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/votes") {

		access, s := util.IsCookie(w, r)
		if !access {
			http.Redirect(w, r, "/signin", 200)
			return
		}

		pid := r.URL.Query().Get("id")
		lukas := r.FormValue("like")
		dislike := r.FormValue("dislike")

		if r.Method == "POST" {

			if lukas == "1" {
				voteLike(w, r, pid, "post", s)
			}

			if dislike == "1" {
				voteDislike(w, r, pid, "post", s)
			}
		}
		http.Redirect(w, r, "post?id="+pid, 302)
	}
}

// /VotesComment function
func VotesComment(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/votes/comment") {

		access, s := util.IsCookie(w, r)
		if !access {
			http.Redirect(w, r, "/signin", 200)
			return
		}

		commentID := r.URL.Query().Get("cid")
		commentDis := r.FormValue("comdis")
		commentLike := r.FormValue("comlike")

		pidc := r.FormValue("pidc")

		if r.Method == "POST" {

			if commentLike == "1" {
				voteLike(w, r, commentID, "comment", s)
			}

			if commentDis == "1" {
				voteDislike(w, r, commentID, "comment", s)
			}
			http.Redirect(w, r, "/post?id="+pidc, 302)
		}
	}
}
fix
func voteDislike(w http.ResponseWriter, r *http.Request, id, any string, s structure.Session) {

	vote := models.Votes{}
	field := any + "_id"
	table := any + "s"

	DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).Scan(&s.UserID)

	DB.QueryRow("SELECT id FROM voteState WHERE "+field+"=? AND user_id=?", id, s.UserID).Scan(&vote.ID)

	if vote.ID == 0 {
		fmt.Println("init dislike")

		_, err = DB.Exec("INSERT INTO voteState("+field+", user_id, dislike_state, like_state) VALUES( ?, ?, ?,?)", id, s.UserID, 1, 0)
		err = DB.QueryRow("SELECT count_dislike FROM "+table+" WHERE id=?", id).Scan(&vote.OldDislike)

		_, err = DB.Exec("UPDATE "+table+" SET count_dislike=? WHERE id=?", vote.OldDislike+1, id)

		//_, err = DB.Exec("UPDATE voteState SET dislike_state = ? WHERE "+field+"=? and user_id", 1, id, s.UserID)
		if err != nil {
			panic(err)
		}

	} else {
		err = DB.QueryRow("SELECT count_like FROM "+table+" WHERE id=?", id).Scan(&vote.OldLike)
		err = DB.QueryRow("SELECT count_dislike FROM "+table+" WHERE id=?", id).Scan(&vote.OldDislike)

		DB.QueryRow("SELECT like_state, dislike_state FROM voteState where "+field+"=? and user_id=?", id, s.UserID).Scan(&vote.LikeState, &vote.DislikeState)

		//set dislike
		if vote.LikeState == 0 && vote.DislikeState == 1 {

			vote.OldDislike--
			_, err = DB.Exec("UPDATE "+table+" SET count_dislike = ? WHERE id=?", vote.OldDislike, id)
			_, err = DB.Exec("UPDATE voteState SET  dislike_state=? WHERE "+field+"=? and user_id=?", 0, id, s.UserID)

			fmt.Println("case 2 like 0, dis 1")
		}

		//set dislike -> to like
		if vote.LikeState == 1 && vote.DislikeState == 0 {
			fmt.Println("case 3 like 1, dis 0")

			vote.OldDislike++
			vote.OldLike--
			_, err = DB.Exec("UPDATE "+table+" SET count_dislike = ? WHERE id=?", vote.OldDislike, id)
			_, err = DB.Exec("UPDATE "+table+" SET count_like = ? WHERE id=?", vote.OldLike, id)
			_, err = DB.Exec("UPDATE voteState SET like_state = ?, dislike_state=? WHERE "+field+"=? and user_id=?", 0, 1, id, s.UserID)
		}

		if vote.LikeState == 0 && vote.DislikeState == 0 {
			fmt.Println("case 4 like 0, dis 0 LS 0, DS 1")
			vote.OldDislike++
			_, err = DB.Exec("UPDATE "+table+" SET count_dislike=? WHERE id=?", vote.OldDislike, id)
			_, err = DB.Exec("UPDATE voteState SET dislike_state = ?, like_state=? WHERE "+field+"=? and user_id=?", 1, 0, id, s.UserID)
		}

		if err != nil {
			panic(err)
		}
	}
}
func voteLike(w http.ResponseWriter, r *http.Request, id, any string, s structure.Session) {

	vote := models.Votes{}
	field := any + "_id"
	table := any + "s"
	//get current UserId by uuid
	DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).Scan(&s.UserID)
	//get by post_id and user_id -> row -> in voteState, if not -> create new row, set chenge likeState = 1, add post by ID  - like_count + 1
	DB.QueryRow("SELECT id FROM voteState where "+field+"=? and user_id=?", id, s.UserID).Scan(&vote.ID)
	fmt.Println(vote.ID, s.UserID, "start", id, table, "table", "init Like field", field)

	if vote.ID == 0 {

		_, err = DB.Exec("INSERT INTO voteState( "+field+"user_id, like_state, dislike_state) VALUES(?, ?, ?, ?)", id, s.UserID, 1, 0)
		err = DB.QueryRow("SELECT count_like FROM "+table+" WHERE id=?", id).Scan(&vote.OldLike)
		_, err = DB.Exec("UPDATE "+table+" SET count_like=? WHERE id=?", vote.OldLike+1, id)

		if err != nil {
			log.Println(err)
		}

	} else {
		//if post -> liked or Disliked -> get CountLike & Dislike current Post, and get LikeState & dislike State
		err = DB.QueryRow("SELECT count_like FROM "+table+"  WHERE id=?", id).Scan(&vote.OldLike)
		err = DB.QueryRow("SELECT count_dislike FROM "+table+"  WHERE id=?", id).Scan(&vote.OldDislike)
		if err != nil {
			log.Println(err)
		}
		DB.QueryRow("SELECT like_state, dislike_state FROM voteState where "+field+"=?  and user_id=?", id, s.UserID).Scan(&vote.LikeState, &vote.DislikeState)

		fmt.Println(" old Dislike & like", vote.OldDislike, vote.OldLike)
		//set like
		if vote.LikeState == 1 && vote.DislikeState == 0 {
			fmt.Println("case 2 like 1, dis 0")

			vote.OldLike--
			_, err = DB.Exec("UPDATE "+table+"  SET count_like = ? WHERE id= ?", vote.OldLike, id)
			_, err = DB.Exec("UPDATE voteState SET like_state = ? WHERE "+field+"=?  and user_id=?", 0, id, s.UserID)
		}
		//set dislike -> to like
		if vote.LikeState == 0 && vote.DislikeState == 1 {
			fmt.Println("case 3 like 0, dis 1")

			vote.OldDislike--
			vote.OldLike++
			_, err = DB.Exec("UPDATE "+table+"  SET count_dislike = ?, count_like=? WHERE id=?", vote.OldDislike, vote.OldLike, id)
			_, err = DB.Exec("UPDATE voteState SET like_state = ?, dislike_state=? WHERE "+field+"=?  and user_id=?", 1, 0, id, s.UserID)
		}
		//set like,
		if vote.LikeState == 0 && vote.DislikeState == 0 {
			fmt.Println("case 4 like 0, dis 0, Ls 1, DS 0")

			vote.OldLike++
			_, err = DB.Exec("UPDATE "+table+" SET count_like=? WHERE id=?", vote.OldLike, id)
			_, err = DB.Exec("UPDATE voteState SET like_state = ?, dislike_state=? WHERE "+field+"= ?  and user_id=?", 1, 0, id, s.UserID)
		}

	}
}

//Likes table, filed posrid, userid, state_id
// 0,1,2 if state ==0, 1 || 2,
// next btn, if 1 == 1, state =0
