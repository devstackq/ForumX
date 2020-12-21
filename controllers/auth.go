package controllers

import (
	"ForumX/general"
	"ForumX/models"
	"ForumX/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/oauth2"
)

var (
	//GoogleConfig *oauth2.Config
	oAuthState = "pseudo-random"
)

//Signup system function
func Signup(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/signup") {

		//	utils.CheckMethod("GET", "signup", auth, w, utils.DisplayTemplate(w, "signup", auth))
		fmt.Println("her -1")
		utils.CheckMethod(r.Method, "signup", auth, w, func(writer http.ResponseWriter) {
			//if r.Method == "POST" {
			fmt.Println("her0")
			intAge, err := strconv.Atoi(r.FormValue("age"))
			if err != nil {
				log.Println(err)
			}
			iB := utils.FileByte(r, "user")
			//checkerEmail & password
			if utils.IsEmailValid(r.FormValue("email")) {

				fullName := r.FormValue("fullname")
				if fullName == "" {
					fullName = "No Name"
				}
				if intAge == 0 {
					intAge = 16
				}
				utils.AuthType = r.FormValue("authType")
				pwd, _ := r.Form["password"]

				if pwd[0] == pwd[1] {
					if utils.IsPasswordValid(r.FormValue("password")) {
						u := models.User{
							FullName: fullName,
							Email:    r.FormValue("email"),
							Username: r.FormValue("username"),
							Age:      intAge,
							Sex:      r.FormValue("sex"),
							City:     r.FormValue("city"),
							Image:    iB,
							Password: r.FormValue("password"),
						}
						fmt.Println("her1")
						u.Signup(writer, r)
						http.Redirect(writer, r, "/signin", 302)
					} else {
						msg := "Incorrect password: must be 8 symbols, 1 big, 1 special character, example: 9Password!"
						utils.DisplayTemplate(writer, "signup", &msg)
					}
				} else {
					msg := "Password fields: not match epta"
					utils.DisplayTemplate(writer, "signup", &msg)
				}
			} else {
				msg := "Incorrect email address: example god@mail.kz"
				utils.DisplayTemplate(writer, "signup", &msg)
			}
			//}
		})

	}
}

//Signin system function
func Signin(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/signin") {

		if r.Method == "GET" {
			utils.DisplayTemplate(w, "signin", &msg)
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
				utils.AuthType = "default"
				u := models.User{
					Email:    person.Email,
					Username: person.Username,
					Password: person.Password,
				}
				u.Signin(w, r)
			}
		}
	}
}

// Logout system function
func Logout(w http.ResponseWriter, r *http.Request) {

	if utils.URLChecker(w, r, "/logout") {
		if r.Method == "GET" {
			models.Logout(w, r)
			http.Redirect(w, r, "/", 302)
		}
	}
}

//GoogleLogin func
func GoogleSignin(w http.ResponseWriter, r *http.Request) {
	url := utils.GoogleConfig.AuthCodeURL(oAuthState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

//GoogleUserData func
func GoogleUserData(w http.ResponseWriter, r *http.Request) {

	utils.AuthType = "google"
	content, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))
	utils.Code = r.FormValue("code")

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

	token, err := utils.GoogleConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	utils.Token = token.AccessToken

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
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func GithubUserData(w http.ResponseWriter, r *http.Request) {

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

	githubSession := general.Session{}
	gitUserData := models.User{}
	json.Unmarshal(responseBody, &githubSession)
	fmt.Println(githubSession, "github token")
	utils.Token = githubSession.AccessToken
	json.Unmarshal(GetGithubData(githubSession.AccessToken), &gitUserData)
	SigninSideService(w, r, gitUserData)
}

func GetGithubData(token string) []byte {

	utils.AuthType = "github"

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
func SigninSideService(w http.ResponseWriter, r *http.Request, u models.User) {

	if utils.IsRegistered(w, r, u.Email) {
		u := models.User{
			Email:    u.Email,
			FullName: u.Name,
		}
		u.Signin(w, r) //login
	} else {
		//if github = location -> else Almaty
		u := models.User{
			Email:    u.Email,
			FullName: u.Name,
			Age:      16,
			Sex:      "Male",
			City:     u.Location,
			Image:    utils.FileByte(r, "user"),
		}
		u.Signup(w, r)
		u.Signin(w, r)
	}
}

//alnikolaevich, Auth
