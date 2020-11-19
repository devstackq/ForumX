package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/devstackq/ForumX/models"
	util "github.com/devstackq/ForumX/utils"
	"golang.org/x/oauth2"
)

var (
	//GoogleConfig *oauth2.Config
	oAuthState  = "pseudo-random"
	defaultAuth string
	Code        string
)

//Signup system function
func Signup(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/signup") {

		if r.Method == "GET" {
			util.DisplayTemplate(w, "signup", &auth)
		}

		if r.Method == "POST" {

			intAge, err := strconv.Atoi(r.FormValue("age"))
			if err != nil {
				log.Println(err)
			}
			iB := util.FileByte(r, "user")
			//checkerEmail & password
			if util.IsEmailValid(r.FormValue("email")) {

				fullName := r.FormValue("fullname")
				if fullName == "" {
					fullName = "Noname"
				}
				if intAge == 0 {
					intAge = 16
				}
				if util.IsPasswordValid(r.FormValue("password")) {

					u := models.User{
						FullName: fullName,
						Email:    r.FormValue("email"),
						Age:      intAge,
						Sex:      r.FormValue("sex"),
						City:     r.FormValue("city"),
						Image:    iB,
						Password: r.FormValue("password"),
					}
					u.Signup(w, r, "default")
					http.Redirect(w, r, "/signin", 302)
				} else {
					msg := "Password must be 8 symbols, 1 big, 1 special character, example: 9Password!"
					util.DisplayTemplate(w, "signup", &msg)
				}
			} else {
				msg := "Incorrect email address, example god@mail.com"
				util.DisplayTemplate(w, "signup", &msg)
			}
		}
	}
}

//Signin system function
func Signin(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/signin") {

		if r.Method == "GET" {
			util.DisplayTemplate(w, "signin", &msg)
		}

		if r.Method == "POST" {

			var person models.User
			err := json.NewDecoder(r.Body).Decode(&person)
			//badrequest
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if person.Type == "default" {
				defaultAuth = "default"
				u := models.User{
					Email:    person.Email,
					Password: person.Password,
				}
				u.Signin(w, r, "default")

			} else if person.Type == "google" {
				//	GoogleLogin(w, r)
				//check if not exist user Db -> Signup else ->  Signin
				fmt.Println("google auth", person.Type)
			} else if person.Type == "github" {
				fmt.Println("todo github auth")
			}
		}
	}
}

// Logout system function
func Logout(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/logout") {
		if r.Method == "GET" {
			models.Logout(w, r, defaultAuth)
			http.Redirect(w, r, "/", 302)
		}
	}
}

//GoogleLogin func
func GoogleLogin(w http.ResponseWriter, r *http.Request) {

	defaultAuth = "google"
	url := util.GoogleConfig.AuthCodeURL(oAuthState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

//GoogleUserData func
func GoogleUserData(w http.ResponseWriter, r *http.Request) {

	content, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))

	util.Code = r.FormValue("code")

	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	gData := models.User{}
	json.Unmarshal(content, &gData)

	if util.IsRegistered(w, r, gData.Email) {
		u := models.User{
			Email:    gData.Email,
			FullName: gData.Name,
		}
		u.Signin(w, r, "google")
	} else {
		u := models.User{
			Email:    gData.Email,
			FullName: gData.Name,
			Age:      16,
			Sex:      "male",
			City:     "Almaty",
			Image:    util.FileByte(r, "user"),
		}
		u.Signup(w, r, "google")
		u.Signin(w, r, "google")
	}
}

func getUserInfo(state, code string) ([]byte, error) {
	//state random string todo
	if state != oAuthState {
		return nil, fmt.Errorf("invalid oauth state")
	}

	token, err := util.GoogleConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	util.Token = token.AccessToken

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	return contents, nil
}
