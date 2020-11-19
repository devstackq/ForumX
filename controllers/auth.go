package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/devstackq/ForumX/models"
	util "github.com/devstackq/ForumX/utils"
	"golang.org/x/oauth2"

	structure "github.com/devstackq/ForumX/general"
)

var (
	//GoogleConfig *oauth2.Config
	oAuthState = "pseudo-random"
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
					u.Signup(w, r)
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
				util.AuthType = "default"
				u := models.User{
					Email:    person.Email,
					Password: person.Password,
				}
				u.Signin(w, r)

			}
		}
	}
}

// Logout system function
func Logout(w http.ResponseWriter, r *http.Request) {

	if util.URLChecker(w, r, "/logout") {
		if r.Method == "GET" {
			models.Logout(w, r)
			http.Redirect(w, r, "/", 302)
		}
	}
}

//GoogleLogin func
func GoogleSignin(w http.ResponseWriter, r *http.Request) {

	url := util.GoogleConfig.AuthCodeURL(oAuthState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

//GoogleUserData func
func GoogleUserData(w http.ResponseWriter, r *http.Request) {
	util.AuthType = "google"
	content, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))
	util.Code = r.FormValue("code")

	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	googleData := models.User{}
	json.Unmarshal(content, &googleData)

	SigninSideService(w, r, googleData)

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

func GithubSignin(w http.ResponseWriter, r *http.Request) {

	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&scope=user:email&redirect_uri=%s", "b8f04afed4e89468b1cf", "http://localhost:6969/githubUserInfo")
	http.Redirect(w, r, redirectURL, 301)
}

func GithubUserData(w http.ResponseWriter, r *http.Request) {
	util.AuthType = "github"
	reqBody := map[string]string{"client_id": "b8f04afed4e89468b1cf", "client_secret": "6ab9cf0c812fbf5ed4e44aea599c418bd3d8cf08", "code": r.URL.Query().Get("code")}
	reqJSON, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(reqJSON))
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	responseBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
	}

	githubData := structure.Session{}

	json.Unmarshal(responseBody, &githubData)
	fmt.Println(githubData)

	gitData := models.User{}
	util.Token = githubData.AccessToken
	json.Unmarshal(GetGithubData(githubData.AccessToken), &gitData)

	SigninSideService(w, r, gitData)

}
func SigninSideService(w http.ResponseWriter, r *http.Request, u models.User) {

	if util.IsRegistered(w, r, u.Email) {
		u := models.User{
			Email:    u.Email,
			FullName: u.Name,
		}
		u.Signin(w, r)
	} else {
		//if github = location -> else Almaty
		u := models.User{
			Email:    u.Email,
			FullName: u.Name,
			Age:      16,
			Sex:      "Male",
			City:     u.Location,
			Image:    util.FileByte(r, "user"),
		}
		u.Signup(w, r)
		u.Signin(w, r)
	}
}
func GetGithubData(token string) []byte {

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println(err)
	}
	responseBody, _ := ioutil.ReadAll(resp.Body)

	return responseBody
}
