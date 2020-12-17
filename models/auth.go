package models

import (
	"ForumX/general"
	"ForumX/utils"
	"fmt"
	"log"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

//Signup func
func (u *User) Signup(w http.ResponseWriter, r *http.Request) {

	users := []User{}
	var hashPwd []byte
	if utils.AuthType == "default" {
		hashPwd, err = bcrypt.GenerateFromPassword([]byte(u.Password), 8)
		if err != nil {
			log.Println(err)
		}
	}
	//check email by unique, if have same email
	checkEmail, err := DB.Query("SELECT email FROM users")
	if err != nil {
		log.Println(err)
	}

	for checkEmail.Next() {
		user := User{}
		var email string
		err = checkEmail.Scan(&email)
		if err != nil {
			log.Println(err.Error())
		}

		user.Email = email
		users = append(users, user)
	}

	for _, v := range users {
		if v.Email == u.Email {
			msg = "Not unique email lel"
			utils.DisplayTemplate(w, "signup", &msg)
			log.Println(err)
		}
	}
	userPrepare, err := DB.Prepare(`INSERT INTO users(full_name, email, password, age, sex, created_time, city, image) VALUES(?,?,?,?,?,?,?,?)` )
	if err != nil {
		log.Println(err)
	}
	_, err = userPrepare.Exec(u.FullName, u.Email, hashPwd, u.Age, u.Sex, time.Now(), u.City, u.Image) 
	if err != nil {
		log.Println(err)
	}
	defer userPrepare.Close()
}

//Signin function
func (uStr *User) Signin(w http.ResponseWriter, r *http.Request) {

	u := DB.QueryRow("SELECT id FROM users WHERE email=?", uStr.Email)

	var user User
	var err error
	err = u.Scan(&user.ID)
	
	if utils.AuthType == "default" {
		u := DB.QueryRow("SELECT id, password FROM users WHERE email=?", uStr.Email)
		//check pwd, if not correct, error
		err = u.Scan(&user.ID, &user.Password)

		if err != nil {
			utils.AuthError(w, r, err, "user not found", utils.AuthType)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(uStr.Password))
		if err != nil {
			utils.AuthError(w, r, err, "password incorrect", utils.AuthType)
			return
		}
	}
	//get user by Id, and write session struct
	s := general.Session{
		UserID: user.ID,
	}
	uuid := uuid.Must(uuid.NewV4(), err).String()
	if err != nil {
		utils.AuthError(w, r, err, "uuid trouble", utils.AuthType)
		return
	}
	//create uuid and set uid DB table session by userid,

	userPrepare, err := DB.Prepare(`INSERT INTO session(uuid, user_id) VALUES (?, ?)` )
	if err != nil {
		log.Println(err)
	}
	_, err = userPrepare.Exec(uuid, s.UserID) 
	defer userPrepare.Close()

	if err != nil {
		utils.AuthError(w, r, err, "the user is already in the system", utils.AuthType)
		//get ssesion id, by local struct uuid
		return
	}
	// get user in info by session Id
	err = DB.QueryRow("SELECT id, uuid FROM session WHERE user_id = ?", s.UserID).Scan(&s.ID, &s.UUID)
	if err != nil {
		utils.AuthError(w, r, err, "not find user from session", utils.AuthType)
		return
	}

	//set cookie 9128ueq9widjaisdh238yrhdeiuwandijsan
	cookie := http.Cookie{
		Name:     "_cookie",
		Value:    s.UUID,
		Path:     "/",
		Expires:  time.Now().Add(300 * time.Minute),
		HttpOnly: false,
	}
	http.SetCookie(w, &cookie)
	utils.AuthError(w, r, nil, "success", utils.AuthType)
	fmt.Println(utils.AuthType, "auth type")
}

//Logout function
func Logout(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("_cookie")
	if err != nil {
		fmt.Println(err, "cookie err")
	}
	//add cookie -> fields uuid
	s := general.Session{UUID: cookie.Value}
	//get ssesion id, by local struct uuid
	DB.QueryRow("SELECT id FROM session WHERE uuid = ?", s.UUID).
		Scan(&s.ID)
	fmt.Println(s.ID, "user id deleted session")
	//delete session by id session
	_, err = DB.Exec("DELETE FROM session WHERE id = ?", s.ID)

	if err != nil {
		log.Println(err)
	}
	// then delete cookie from client
	utils.DeleteCookie(w)

	if utils.AuthType == "google" {
		_, err = http.Get("https://accounts.google.com/o/oauth2/revoke?token=" + utils.Token)
		if err != nil {
			log.Println(err)
		}
	}
}
