package models

//Votes struct
type Votes struct {
	ID         int
	Like       int
	Dislike    int
	PostID     int
	UserID     int
	Voted      bool
	CommentID  int
	OldLike    int
	OldDislike int
}

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

// ----------- dilike
// var p, u int
// err = DB.QueryRow("SELECT post_id, user_id FROM likes WHERE post_id=? AND user_id=?", pid, s.UserID).Scan(&p, &u)

// if p == 0 && u == 0 {

// 	oldlike := 0
// 	err = DB.QueryRow("select count_dislike from posts where id=?", pid).Scan(&oldlike)
// 	nv := oldlike + 1
// 	_, err = DB.Exec("UPDATE  posts SET count_dislike = ? WHERE id= ?", nv, pid)
// 	if err != nil {
// 		panic(err)
// 	}
// 	_, err = DB.Exec("INSERT INTO likes(post_id, user_id, state_id) VALUES( ?, ?, ?)", pid, s.UserID, 0)

// 	if err != nil {
// 		panic(err)
// 	}
// }
