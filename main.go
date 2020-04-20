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
	"time"
)

type Entry struct {
	Id          int
	Title       string
	Body        string
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
	statement, err = database.Prepare("INSERT INTO entries (title, body, tags, created) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec(e1().Title, e1().Body, e1().Tags, e1().Created)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	database, _ = sql.Open("sqlite3", "./justblog.db")
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
	e := Entry{
		Title:   r.FormValue("title"),
		Body:    r.FormValue("body"),
		Tags:    r.FormValue("tags"),
		Created: time.Now(),
	}
	statement, err := database.Prepare("INSERT INTO entries (title, body, tags, created) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec(e.Title, e.Body, e.Tags, e.Created)
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
		rows.Scan(&e.Id, &e.Title, &e.Body, &e.Tags, &e.Created)
		e.CreatedText = e.Created.Format(time.RFC1123)
		data.Entries = append(data.Entries, e)
	}

	box := packr.NewBox("./templates")
	s, _ := box.FindString("index.html")
	tmpl, _ := template.New("index").Parse(s)
	tmpl.Execute(w, data)
}

func e1() Entry {
	return Entry{
		Title:   "Gedanke oder Gefühl",
		Body:    "Es geht nicht darum, sich auf den Atem zu konzentrieren. Das wäre viel zu anstrengend. Es geht vielmehr darum, Ablenkungen wahrzunehmen. Ein Ablenkung kann ein Gedanke oder ein Gefühl sein. Sobald ich einen Gedanken bewusst wahrnehme, mache ich einen mentalen Haken Gedanke und kehre bewusst zum Atem zurück. Der Atem ist nur ein Anker, zu dem ich immer zurückkehren kann.",
		Created: time.Now(),
		Tags:    "#Meditation #Privat",
	}
}
