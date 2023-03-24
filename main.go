package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

// type User struct {
// 	Name string `json:"name"`
// 	Age  uint16 `json:"age"`
// 	// 	Money               int16
// 	// 	AvgGrades, Hapiness float64
// 	// 	Hobbies             []string
// }

type Article struct {
	Id                     uint16
	Title, Anons, FullText string
}

var posts = []Article{}
var showPost = Article{}

// func (u User) getAllInfo() string {
// 	return fmt.Sprintf("Name: %s\n Age: %d\n Money: %d\n AvgGrades: %f\n Hapiness: %f", u.Name, u.Age, u.Money, u.AvgGrades, u.Hapiness)
// }

// func (u *User) SetNewName(newName string) {
// 	u.Name = newName
// }

// func home_page(w http.ResponseWriter, r *http.Request) {
// 	// fmt.Fprintf(w, "<b>GO is super easy!</b>")
// 	bob := User{"Bob", 25, -50, 4.2, 0.8, []string{"Football", "Skate", "Dance"}}
// 	tmpl, _ := template.ParseFiles("templates/home_page.html")
// 	tmpl.Execute(w, bob)
// }

// func contacts_page(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Contacts page!")
// }

// func handleRequest() {
// 	http.HandleFunc("/", home_page)
// 	http.HandleFunc("/contacts/", contacts_page)
// 	http.ListenAndServe(":8080", nil)
// }

func index(w http.ResponseWriter, r *http.Request) {
	// 	// fmt.Fprintf(w, "<b>GO is super easy!</b>")
	// 	bob := User{"Bob", 25, -50, 4.2, 0.8, []string{"Football", "Skate", "Dance"}}
	tmpl, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprint(w, err.Error())
	}

	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/Golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Data mining
	res, err := db.Query("SELECT * FROM `articles`")
	if err != nil {
		panic(err)
	}

	posts = []Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}
		// fmt.Printf("Post %s with id %d\n", post.Title, post.Id) //, post.Anons, post.FullText)

		posts = append(posts, post)
	}

	tmpl.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request) {
	// 	// fmt.Fprintf(w, "<b>GO is super easy!</b>")
	// 	bob := User{"Bob", 25, -50, 4.2, 0.8, []string{"Football", "Skate", "Dance"}}
	tmpl, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
	tmpl.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	// fmt.Println(title)
	// fmt.Println(anons)
	// fmt.Println(full_text)

	if title == "" || anons == "" || full_text == "" {
		fmt.Fprint(w, "field is empty!")
	} else {
		db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/Golang")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		// data installing
		insert, err := db.Query(fmt.Sprintf("INSERT INTO `articles` (`title`, `anons`, `full_text`) VALUES('%s', '%s', '%s')", title, anons, full_text))
		if err != nil {
			panic(err)
		}
		defer insert.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther) //301
	}
}

func show_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	tmpl, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprint(w, err.Error())
	}

	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/Golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Data mining
	res, err := db.Query(fmt.Sprintf("SELECT * FROM `articles` WHERE `id` = '%s'", vars["id"]))
	if err != nil {
		panic(err)
	}

	showPost = Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}
		// fmt.Printf("Post %s with id %d\n", post.Title, post.Id) //, post.Anons, post.FullText)
		showPost = post
		posts = append(posts, post)
	}

	tmpl.ExecuteTemplate(w, "show", showPost)

	// w.WriteHeader(http.StatusOK)
	// fmt.Fprintf(w, "ID: %v\n", vars["id"])
}

func handlefunc() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create/", create).Methods("GET")
	rtr.HandleFunc("/save_article/", save_article).Methods("POST")
	// http.HandleFunc("/contacts/", contacts_page)
	rtr.HandleFunc("/post/{id:[0-9]+}/", show_post).Methods("GET")

	http.Handle("/", rtr)
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":8080", nil)
}

func main() {
	// handleRequest()
	handlefunc()

	// db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/Golang")
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()

	// data installing
	// insert, err := db.Query("INSERT INTO `users` (`name`, `age`) VALUES('Bob', 35)")
	// if err != nil {
	// 	panic(err)
	// }
	// defer insert.Close()

	// Data mining
	// res, err := db.Query("SELECT `name`, `age` FROM `users`")
	// if err != nil {
	// 	panic(err)
	// }
	// for res.Next() {
	// 	var user User
	// 	err = res.Scan(&user.Name, &user.Age)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Printf("User: %s with age %d\n", user.Name, user.Age)
	// }
}
