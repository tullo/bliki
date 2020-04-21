package main

import (
	"database/sql"
	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Entry struct {
	Id          int
	Title       string
	Body        template.HTML
	Tags        string
	Created     time.Time
	CreatedText string
}

type IndexData struct {
	Entries []Entry
}

var database *sql.DB

func initDatabase() {
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS entries (id INTEGER PRIMARY KEY, title TEXT, body TEXT, tags TEXT, created DATETIME)")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
}

func main() {
	database, _ = sql.Open("sqlite3", "./justblog.db")
	initDatabase()
	mux := mux.NewRouter()
	mux.HandleFunc("/new", newHandler)
	mux.HandleFunc("/create", createHandler)
	mux.HandleFunc("/delete/{id}", deleteHandler)
	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":3000", mux)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	statement, _ := database.Prepare("DELETE FROM entries WHERE id = ?")
	statement.Exec(id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	body := r.FormValue("body")
	tags := r.FormValue("tags")
	created := time.Now()
	statement, err := database.Prepare("INSERT INTO entries (title, body, tags, created) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec(title, body, tags, created)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func newHandler(w http.ResponseWriter, request *http.Request) {
	box := packr.NewBox("./templates")
	s, err := box.FindString("new.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl, _ := template.New("new").Parse(s)
	tmpl.Execute(w, nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Query("SELECT id, title, body, tags, created FROM entries order by created desc ")
	if err != nil {
		log.Fatal(err)
	}
	data := IndexData{}
	for rows.Next() {
		e := Entry{}
		var body string
		rows.Scan(&e.Id, &e.Title, &body, &e.Tags, &e.Created)
		e.CreatedText = e.Created.Format(time.RFC1123)
		e.Body = template.HTML(strings.Replace(body, "\r\n", "<br>", -1))
		data.Entries = append(data.Entries, e)
	}

	box := packr.NewBox("./templates")
	s, _ := box.FindString("index.html")
	tmpl, _ := template.New("index").Parse(s)
	tmpl.Execute(w, data)
}
