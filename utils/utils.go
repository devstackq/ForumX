package util

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	structure "github.com/devstackq/ForumX/general"
)

var (
	DB   *sql.DB
	err  error
	temp = template.Must(template.ParseFiles("view/header.html", "view/category_temp.html", "view/likedpost.html", "view/404page.html", "view/postupdate.html", "view/postuser.html", "view/commentuser.html", "view/userupdate.html", "view/search.html", "view/user.html", "view/commentuser.html", "view/postuser.html", "view/profile.html", "view/signin.html", "view/user.html", "view/signup.html", "view/filter.html", "view/posts.html", "view/comment.html", "view/create.html", "view/footer.html", "view/index.html"))
)

type API struct {
	Authenticated bool
}

func IsAuth(r *http.Request) API {
	var auth API
	for _, cookie := range r.Cookies() {
		if cookie.Name == "_cookie" {
			auth.Authenticated = true
		}
	}
	return auth
}

func CheckForCookies(w http.ResponseWriter, r *http.Request) (bool, structure.Session) {

	var flag, cookieHave bool

	if IsAuth(r).Authenticated {
		cookieHave = true
	}
	if !cookieHave {
		http.Redirect(w, r, "/signin", 302)
	} else {
		//get client cookie
		//set local struct -> cookie value
		cookie, _ := r.Cookie("_cookie")
		s := structure.Session{UUID: cookie.Value}
		var tmp string
		// get userid by Client sessionId
		err = DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).
			Scan(&s.UserID)
		//get uuid by userid, and write UUID data
		if err != nil {
			log.Println(err)
		}
		err = DB.QueryRow("SELECT uuid FROM session WHERE user_id = ?", s.UserID).Scan(&tmp)
		if err != nil {
			log.Println(err)
		}
		//check local and DB session
		if cookie.Value == tmp {
			flag = true
		}
	}
	s := structure.Session{}
	if flag {
		c, _ := r.Cookie("_cookie")
		s.UUID = c.Value
		return flag, s
	}

	return flag, s
}

//check correct letter
func CheckLetter(value string) bool {
	for _, v := range value {
		if v >= 97 && v <= 122 || v >= 65 && v <= 90 && v >= 32 && v <= 64 || v > 128 {
			return true
		}
	}
	return false
}

func DisplayTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err = temp.ExecuteTemplate(w, tmpl, data)
	fmt.Println(err, "exec ERR")
	if err != nil {
		http.Error(w, err.Error(),
			http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
	}
}

func СheckCookieLife(t time.Time, cookie *http.Cookie, w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == "_cookie" {
			//Logout(w, r)
			s := structure.Session{UUID: cookie.Value}
			//get ssesion id, by local struct uuid
			DB.QueryRow("SELECT id FROM session WHERE uuid = ?", s.UUID).
				Scan(&s.ID)

			_, err = DB.Exec("DELETE FROM session WHERE id = ?", s.ID)

			// then delete cookie from client
			cDel := http.Cookie{
				Name:     "_cookie",
				Value:    "",
				Path:     "/",
				Expires:  time.Unix(0, 0),
				HttpOnly: false,
			}
			http.SetCookie(w, &cDel)
			http.Redirect(w, r, "/", http.StatusOK)
			return
		}
	}
}

//find unique liked post
func IsUnique(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

//func for convert receive file - to fileByte
func FileByte(r *http.Request) []byte {

	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("uploadfile")

	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	imgBytes, err := ioutil.ReadAll(file)

	if err != nil {
		log.Println(err)
	}

	return imgBytes
}

//show auth error
func AuthError(w http.ResponseWriter, err error, text string) {
	fmt.Println(text, "errka")
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		m, _ := json.Marshal(text)
		w.Write(m)
		return
	} else {
		w.WriteHeader(200)
		m, _ := json.Marshal(text)
		w.Write(m)
	}
}
