package util

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/devstackq/ForumX/model"
)

var (
	DB   *sql.DB
	err  error
	temp = template.Must(template.ParseFiles("templates/header.html", "templates/category_temp.html", "templates/likedpost.html", "templates/404page.html", "templates/postupdate.html", "templates/postuser.html", "templates/commentuser.html", "templates/userupdate.html", "templates/search.html", "templates/user.html", "templates/commentuser.html", "templates/postuser.html", "templates/profile.html", "templates/signin.html", "templates/user.html", "templates/signup.html", "templates/filter.html", "templates/posts.html", "templates/comment.html", "templates/create.html", "templates/footer.html", "templates/index.html"))
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

func CheckForCookies(w http.ResponseWriter, r *http.Request) (bool, model.Session) {

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
		s := model.Session{UUID: cookie.Value}
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
	s := model.Session{}
	if flag {
		c, _ := r.Cookie("_cookie")
		s.UUID = c.Value
		return flag, s
	}

	return flag, s
}

func CheckLetter(value string) bool {
	for _, v := range value {
		if v >= 97 && v <= 122 || v >= 65 && v <= 90 && v >= 32 && v <= 64 || v > 128 {
			return true
		}
	}
	return false
}

//DisplayTemplate comment
func DisplayTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err = temp.ExecuteTemplate(w, tmpl, data)
	fmt.Println(err, "exec ERR")
	if err != nil {
		http.Error(w, err.Error(),
			http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
	}
}
