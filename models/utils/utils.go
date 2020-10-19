package util

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/devstackq/ForumX/models"
)

var (
	DB  *sql.DB
	err error
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

func CheckForCookies(w http.ResponseWriter, r *http.Request) bool {

	flag := false
	cookieHave := false

	if IsAuth(r).Authenticated {
		cookieHave = true
	}
	log.Println(cookieHave, "cook her")

	// for _, cookie := range r.Cookies() {
	// 	if cookie.Name == "_cookie" {
	// 		cookieHave = true
	// 		break
	// 	}
	// }
	if !cookieHave {
		http.Redirect(w, r, "/signin", 302)
	} else {
		//get client cookie
		//set local struct -> cookie value
		cookie, _ := r.Cookie("_cookie")
		s := models.Session{UUID: cookie.Value}
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
		fmt.Println(flag, "falg")
	}
	return flag
}
