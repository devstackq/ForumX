package main

import (
	"ForumX/config"
	"ForumX/controllers"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type person struct {
	name string
	age  int
}

func (p *person) updatePerson(args ...interface{}) {

	fmt.Println(len(args))
	//for _, v := range data {
	//s := reflect.TypeOf(data)
	if len(args) > 1 {

		for i, arg := range args {
			// if arg.(string) != "" {
			// 	p.name = arg.(string)
			// }
			switch i {
			case 0: // name
				//check inside Struct index, 1 : name, 2 age..., then cast to stirng, then change value strucnt field Name -  by Pointer
				name, ok := arg.(string)
				if !ok {
					log.Println("error")
				} else {
					p.name = name
				}
			case 1:
				age, ok := arg.(int)
				if !ok {
					log.Println("error")
				}
				p.age = age
			default:
				log.Println("Wrong parametes passed")
			}
		}
	} else {
		//comapre type, then set value - variable

	// 	for _, arg := range args {
	// 		if reflect.DeepEqual(arg.(string), arg) {
	// 			if arg.(string) != "" {
	// 				p.name = arg.(string)
	// 			}
	// 		} else if reflect.DeepEqual(arg.(int), arg) {{
	// 			if arg.(int) > 0 {
	// 				p.age = arg.(int)
	// 			}
	// 		}
	// 	}
	// }
}

	// if data.(string){
	// //if  (string) == reflect.TypeOf(data) {
	// 	fmt.Print("ds")
	// 	p.name = fmt.Sprintf("%s" ,data)
	// 	fmt.Print(p)
	// }else if data == reflect.Int {
	// 	a, _  := fmt.Printf("%v \n" ,data)
	// 	print(a)
	// 		p.age = a
	// 	}
}

// func (p *person) updatePersons(arr ...interface{}){

// 	for _, interface := range arr {

// for _, v := range interface {
// 	if  v == reflect.String {
// 		s, _ := fmt.Printf("%v" ,v)
// 		p.name = s
// 	}else if v == reflect.Int {
// 			p.age = fmt.Printf("%v \n" ,v)
// 		}

// }
// 	}
// }
// func (p *person) updatePerson(name string){
// 	p.name = name

// }
func main() {

	//	strunct Method, with *, allows -> Change  specific Struct by Address, change fields
	//init structure
	var Lester = person{name: "Guru", age: 21}
	//update fields by Pointer address struct -> fields

	//перегрузка методов, через struct, и interface
	//Lester.updatePerson( "Soma", 32)
	Lester.updatePerson(32)

	Lester.updatePerson("Jonny")

	fmt.Println(Lester)

	config.Init()
	controllers.Init()
	time.Sleep(10 * time.Second)
	fmt.Println("timer 10 sec")

	перегрузку методов
	use constructor
	use anonim func
	use gorutine
	try -> func use with Interface

	Andrei - узнать конкаренси, паралелизм, как работает в каих случаях использовать
}

// eaxmple reply system https://codewithawa.com/posts/creating-a-comment-and-reply-system-php-and-mysql

//comment system step 3.1
// 1 table create RepliesComment, FK(reply_id) References comments(id) -> Comment -> []ReplyComments
// form inside Client(answer comment )
// Client - form Comment, form each Comments inside comment -> ReplyForm todo

//----------------------
// comment table - comment noraml & comment under reply comment,
// reply table, uid, comment id, content, , comment id,
// insert into - 43 com -  setParentID, 12,
// client - show List comment, if have ParentId-> append Array,
//else show only COmment

//CLient -  answer -> 44com -> Form(setParentId) -> answerId : 14, parentID 44
//------------

//show/hidden by ID -> comment Field textarea
//global variable
// 	DB.QueryRow("SELECT user_id FROM session WHERE uuid = ?", s.UUID).Scan(&s.UserID)
// var toWhom int
// DB.QueryRow("SELECT creator_id FROM comments WHERE id = ?", cid).Scan(&toWhom)

//toggle - windows under comment JS
//answer - COmments -> by userNickname -> ?

//each comment By Id-> show comments
//query - out -> models
//try todo  answer -> to another comment
// interest func - adv feat -> search, pagination

//try - event -> add sound & confetti -Login
// save photo, like - source DB refactor
//config, router refactorr
//if cookie = 0, notify message  user, logout etc
