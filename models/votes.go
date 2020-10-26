package models

//Votes struct
type Votes struct {
	ID      int
	Like    int
	Dislike int
	PostID  int
	UserID  int
	Voted   bool
}
