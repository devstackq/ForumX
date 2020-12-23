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
func (u User) Signup(w http.ResponseWriter, r *http.Request) {

	var hashPwd []byte
	if utils.AuthType == "default" {
		hashPwd, err = bcrypt.GenerateFromPassword([]byte(u.Password), 8)
		if err != nil {
			log.Println(err)
		}
	}
	emailCheck := utils.IsRegistered(w, r, u.Email)
	userCheck := utils.IsRegistered(w, r, u.Username)
	if emailCheck == userCheck {
		userPrepare, err := DB.Prepare(`INSERT INTO users(full_name, email, username, password, age, sex, created_time, city, image) VALUES(?,?,?,?,?,?,?,?,?)`)
		if err != nil {
			log.Println(err)
		}
		_, err = userPrepare.Exec(u.FullName, u.Email, u.Username, hashPwd, u.Age, u.Sex, time.Now(), u.City, u.Image)
		if err != nil {
			log.Println(err)
		}
		defer userPrepare.Close()
	} else {
		if emailCheck == false {
			msg = "Not unique email"
		} else if userCheck == false {
			msg = "Not unique username"
		} else {
			msg = "Not unique email && username"
		}
		if utils.AuthType == "default" {
			utils.DisplayTemplate(w, "signup", &msg)
		}
	}
}

//Signin function dsds
func (uStr *User) Signin(w http.ResponseWriter, r *http.Request) {

	var isUserOrEmail bool

	if uStr.Username != "" {
		isUserOrEmail = true
	} else if uStr.Email != "" {
		isUserOrEmail = false
	}
	var user User

	err = DB.QueryRow("SELECT id FROM users WHERE email=?", uStr.Email).Scan(&user.ID)
	if err != nil {
		log.Println(err)
	}
	if utils.AuthType == "default" {
		if !isUserOrEmail {
			err = DB.QueryRow("SELECT id, password FROM users WHERE email=?", uStr.Email).Scan(&user.ID, &user.Password)
			log.Println("err email")
			if err != nil {
				utils.AuthError(w, r, err, "user by Email not found", utils.AuthType)
				return
			}
		} else if isUserOrEmail {
			err = DB.QueryRow("SELECT id, password FROM users WHERE username=?", uStr.Username).Scan(&user.ID, &user.Password)
			log.Println("errr username")
			if err != nil {
				utils.AuthError(w, r, err, "user by Username not found", utils.AuthType)
				return
			}
		}
		//check pwd, if not correct, error
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

	userPrepare, err := DB.Prepare(`INSERT INTO session(uuid, user_id) VALUES (?, ?)`)
	if err != nil {
		log.Println(err)
	}
	_, err = userPrepare.Exec(uuid, s.UserID)
	defer userPrepare.Close()

	if err != nil {
		utils.AuthError(w, r, err, "the user is already in the system", utils.AuthType)
		//get ssesion id, by local struct uuid
		log.Println(err)
		return
	}
	// get user in info by session Id
	err = DB.QueryRow("SELECT id, uuid FROM session WHERE user_id = ?", s.UserID).Scan(&s.ID, &s.UUID)
	if err != nil {
		utils.AuthError(w, r, err, "not find user from session", utils.AuthType)
		log.Println(err)
		return
	}

	//if cookie not, dele
	// this should give you time in location
	//set cookie 9128ueq9widjaisdh238yrhdeiuwandijsan
	t := time.Now().UTC()
	fmt.Println(time.Now())
	cookie := http.Cookie{
		Name:    "_cookie",
		Value:   s.UUID,
		Path:    "/",
		Expires: t.Add(6 * time.Second),
		//MaxAge:   15,
		HttpOnly: false,
	}
	http.SetCookie(w, &cookie)
	utils.AuthError(w, r, nil, "success", utils.AuthType)
	fmt.Println(utils.AuthType, "auth type")
}

// query -> if Db not have Cookie -> delete client side cookie

//Logout function
func Logout(w http.ResponseWriter, r *http.Request, s general.Session) {
	var idSession int
	err = DB.QueryRow("SELECT id FROM session WHERE uuid = ?", s.UUID).Scan(&idSession)
	if err != nil {
		log.Println(err)
	}
	//delete session by id session
	_, err = DB.Exec("DELETE FROM session WHERE user_id = ?", idSession)

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
