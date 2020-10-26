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

//signup system
func Signup(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/signup" {
		util.DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}

	if r.Method == "GET" {
		util.DisplayTemplate(w, "signup", &auth)
	}

	if r.Method == "POST" {

		iA, _ := strconv.Atoi(r.FormValue("age"))
		iB := util.FileByte(r)

		u := models.Users{
			FullName: r.FormValue("fullname"),
			Email:    r.FormValue("email"),
			Age:      iA,
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
	}
}

//signin system
func Signin(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/signin" {
		util.DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}

	r.Header.Add("Accept", "text/html")
	r.Header.Add("User-Agent", "MSIE/15.0")

	if r.Method == "GET" {
		util.DisplayTemplate(w, "signin", &msg)
	}

	if r.Method == "POST" {
		var person models.Users
		//b, _ := ioutil.ReadAll(r.Body)
		err := json.NewDecoder(r.Body).Decode(&person)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println(person, "person value")

		if person.Type == "default" {

			u := models.Users{
				Email:    person.Email,
				Password: person.Password,
			}

			u.Signin(w, r)
			http.Redirect(w, r, "/profile", 200)

		} else if person.Type == "google" {
			fmt.Println("todo google auth")
			http.Redirect(w, r, "/profile", http.StatusFound)
		} else if person.Type == "github" {
			fmt.Println("todo github auth")
			http.Redirect(w, r, "/profile", http.StatusFound)
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
}

// Logout
func Logout(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/logout" {
		util.DisplayTemplate(w, "404page", http.StatusNotFound)
		return
	}
	if r.Method == "GET" {
		models.Logout(w, r)
		http.Redirect(w, r, "/signin", 302)
	}
}
