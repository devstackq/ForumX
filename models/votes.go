package models

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	structure "github.com/devstackq/ForumX/general"
	util "github.com/devstackq/ForumX/utils"
)

//Votes struct
type Votes struct {
	ID           int
	LikeState    int
	DislikeState int
	PostID       int
	UserID       int
	CommentID    int
	OldLike      int
	OldDislike   int
	CreatorID    int
}

//VoteDislike func
func VoteDislike(w http.ResponseWriter, r *http.Request, id, any string, s structure.Session) {

	vote := Votes{}
	field := any + "_id"
	table := any + "s"

	DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).Scan(&s.UserID)
	DB.QueryRow("SELECT id FROM voteState WHERE "+field+"=? AND user_id=?", id, s.UserID).Scan(&vote.ID)

	err = DB.QueryRow("SELECT creator_id FROM "+table+"  WHERE id=?", id).Scan(&vote.CreatorID)
	objID, _ := strconv.Atoi(id)

	if vote.ID == 0 {
		fmt.Println(vote.ID, s.UserID, "start", objID, table, "table", "init Dislike field", field)

		_, err = DB.Exec("INSERT INTO voteState("+field+", user_id, dislike_state, like_state) VALUES( ?, ?, ?,?)", id, s.UserID, 1, 0)
		if err != nil {
			log.Println(err)
		}
		err = DB.QueryRow("SELECT count_dislike FROM "+table+" WHERE id=?", id).Scan(&vote.OldDislike)
		if err != nil {
			log.Println(err)
		}
		_, err = DB.Exec("UPDATE "+table+" SET count_dislike=? WHERE id=?", vote.OldDislike+1, id)

		if err != nil {
			log.Println(err)
		}
		util.SetVoteNotify(any, vote.CreatorID, s.UserID, objID, false)

	} else {
		err = DB.QueryRow("SELECT count_like FROM "+table+" WHERE id=?", id).Scan(&vote.OldLike)

		err = DB.QueryRow("SELECT count_dislike FROM "+table+" WHERE id=?", id).Scan(&vote.OldDislike)

		DB.QueryRow("SELECT like_state, dislike_state FROM voteState where "+field+"=? and user_id=?", id, s.UserID).Scan(&vote.LikeState, &vote.DislikeState)

		//set dislike default
		if vote.LikeState == 0 && vote.DislikeState == 1 {

			vote.OldDislike--
			_, err = DB.Exec("UPDATE "+table+" SET count_dislike = ? WHERE id=?", vote.OldDislike, id)
			_, err = DB.Exec("UPDATE voteState SET  dislike_state=? WHERE "+field+"=? and user_id=?", 0, id, s.UserID)
			//remove notify table
			util.UpdateVoteNotify(any, vote.CreatorID, s.UserID, objID, 0)
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

			util.UpdateVoteNotify(any, vote.CreatorID, s.UserID, objID, 2)
		}

		if vote.LikeState == 0 && vote.DislikeState == 0 {
			fmt.Println("case 4 like 0, dis 0 LS 0, DS 1")
			vote.OldDislike++
			_, err = DB.Exec("UPDATE "+table+" SET count_dislike=? WHERE id=?", vote.OldDislike, id)
			_, err = DB.Exec("UPDATE voteState SET dislike_state = ?, like_state=? WHERE "+field+"=? and user_id=?", 1, 0, id, s.UserID)

			util.UpdateVoteNotify(any, vote.CreatorID, s.UserID, objID, 2)
		}

		if err != nil {
			log.Println(err)
		}
	}
}

//VoteLike func
func VoteLike(w http.ResponseWriter, r *http.Request, id, any string, s structure.Session) {

	vote := Votes{}
	field := any + "_id"
	table := any + "s"
	//get current UserId by uuid
	DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).Scan(&s.UserID)
	//get by post_id and user_id -> row -> in voteState, if not -> create new row, set chenge likeState = 1, add post by ID  - like_count + 1
	DB.QueryRow("SELECT id FROM voteState where "+field+"=? and user_id=?", id, s.UserID).Scan(&vote.ID)

	err = DB.QueryRow("SELECT creator_id FROM "+table+"  WHERE id=?", id).Scan(&vote.CreatorID)
	pid, _ := strconv.Atoi(id)

	if vote.ID == 0 {
		fmt.Println(vote.ID, s.UserID, "start", id, table, "table", "init Like field", field)

		_, err = DB.Exec("INSERT INTO voteState( "+field+", user_id, like_state, dislike_state) VALUES(?, ?, ?, ?)", id, s.UserID, 1, 0)
		err = DB.QueryRow("SELECT count_like FROM "+table+" WHERE id=?", id).Scan(&vote.OldLike)
		_, err = DB.Exec("UPDATE "+table+" SET count_like=? WHERE id=?", vote.OldLike+1, id)

		if err != nil {
			log.Println(err)
		}
		util.SetVoteNotify(any, vote.CreatorID, s.UserID, pid, true)

	} else {
		//if post -> liked or Disliked -> get CountLike & Dislike current Post, and get LikeState & dislike State
		err = DB.QueryRow("SELECT count_like FROM "+table+"  WHERE id=?", id).Scan(&vote.OldLike)
		err = DB.QueryRow("SELECT count_dislike FROM "+table+"  WHERE id=?", id).Scan(&vote.OldDislike)
		if err != nil {
			log.Println(err)
		}
		DB.QueryRow("SELECT like_state, dislike_state FROM voteState where "+field+"=?  and user_id=?", id, s.UserID).Scan(&vote.LikeState, &vote.DislikeState)

		fmt.Println(" old Dislike & like", vote.OldDislike, vote.OldLike)

		//set like 0, default
		if vote.LikeState == 1 && vote.DislikeState == 0 {

			fmt.Println("l-1, d-0 -> l-0,d-0")
			vote.OldLike--
			_, err = DB.Exec("UPDATE "+table+"  SET count_like = ? WHERE id= ?", vote.OldLike, id)
			_, err = DB.Exec("UPDATE voteState SET like_state = ? WHERE "+field+"=?  and user_id=?", 0, id, s.UserID)

			//util.RemoveVoteNotify(any, vote.CreatorID, s.UserID, pid)
			util.UpdateVoteNotify(any, vote.CreatorID, s.UserID, pid, 0)

		}
		//set dislike -> to like
		if vote.LikeState == 0 && vote.DislikeState == 1 {

			fmt.Println("l-0,d-1, -> l-1, d-0")
			vote.OldDislike--
			vote.OldLike++
			_, err = DB.Exec("UPDATE "+table+"  SET count_dislike = ?, count_like=? WHERE id=?", vote.OldDislike, vote.OldLike, id)
			_, err = DB.Exec("UPDATE voteState SET like_state = ?, dislike_state=? WHERE "+field+"=?  and user_id=?", 1, 0, id, s.UserID)

			//add like notify &  remove DislikeNotify
			util.UpdateVoteNotify(any, vote.CreatorID, s.UserID, pid, 1)

			//util.RemoveVoteNotify(any, vote.CreatorID, s.UserID, pid)
			//util.SetVoteNotify(any, vote.CreatorID, s.UserID, pid, true)
		}
		//set like,
		if vote.LikeState == 0 && vote.DislikeState == 0 {

			fmt.Println("l-0, d-0 -> l-1, d-0")
			vote.OldLike++
			_, err = DB.Exec("UPDATE "+table+" SET count_like=? WHERE id=?", vote.OldLike, id)
			_, err = DB.Exec("UPDATE voteState SET like_state = ?, dislike_state=? WHERE "+field+"= ?  and user_id=?", 1, 0, id, s.UserID)

			//util.SetVoteNotify(any, vote.CreatorID, s.UserID, pid, true)
			util.UpdateVoteNotify(any, vote.CreatorID, s.UserID, pid, 1)

		}
	}
}
