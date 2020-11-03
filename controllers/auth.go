package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/devstackq/ForumX/models"
	util "github.com/devstackq/ForumX/utils"
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

				if util.IsPasswordValid(r.FormValue("password")) {

					u := models.User{
						FullName: r.FormValue("fullname"),
						Email:    r.FormValue("email"),
						Age:      intAge,
						Sex:      r.FormValue("sex"),
						City:     r.FormValue("city"),
						Image:    iB,
						Password: r.FormValue("password"),
					}

					err = u.Signup(w, r)
					if err != nil {
						log.Println(err)
					}
					http.Redirect(w, r, "/signin", 301)
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

			fmt.Println(person, "person data")

			if person.Type == "default" {

				u := models.User{
					Email:    person.Email,
					Password: person.Password,
				}

				u.Signin(w, r)
				http.Redirect(w, r, "/profile", 200)

			} else if person.Type == "google" {
				fmt.Println("todo google auth")
				http.Redirect(w, r, "/profile", 301)
			} else if person.Type == "github" {
				fmt.Println("todo github auth")
				http.Redirect(w, r, "/profile", 301)
			}
			//w.Header().Set("Access-Control-Allow-Origin", "*")
		}
	}
}

// Logout system function
func Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logout")
	if util.URLChecker(w, r, "/logout") {
		if r.Method == "GET" {
			models.Logout(w, r)
			http.Redirect(w, r, "/", 301)
		}
	}
}
