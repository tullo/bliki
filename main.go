package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type Entry struct {
	Id          int
	Title       string
	Body        template.HTML
	Tags        string
	Public      int
	Created     time.Time
	CreatedText string
}

type IndexData struct {
	Entries []Entry
}

var database *sql.DB

func initDatabase() {
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS entries (id INTEGER PRIMARY KEY, title TEXT, body TEXT, tags TEXT, created DATETIME, public INTEGER)")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	database, _ = sql.Open("sqlite3", "./bliki.db")
	initDatabase()
	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/new", basicAuth(newHandler))
	router.HandleFunc("/delete/{id}", basicAuth(deleteHandler))
	router.HandleFunc("/create", basicAuth(createHandler))
	router.HandleFunc("/update/{id}", basicAuth(updateHandler))
	router.HandleFunc("/edit/{id}", basicAuth(editHandler))
	router.HandleFunc("/admin", basicAuth(adminHandler))
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("listing on port %s...", port)
	_ = http.ListenAndServe(":"+port, router)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	title := r.FormValue("title")
	body := r.FormValue("body")
	tags := r.FormValue("tags")
	public := r.FormValue("public")
	if len(public) == 0 {
		public = "0"
	}
	statement, err := database.Prepare("UPDATE entries SET title = ?, body = ?, tags = ?, public = ? WHERE ID = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec(title, body, tags, public, id)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	statement, _ := database.Prepare("DELETE FROM entries WHERE id = ?")
	statement.Exec(id)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	stmt, err := database.Prepare("SELECT id, title, body, tags, created, public FROM entries WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := stmt.QueryRow(id)
	e := Entry{}
	var body string
	row.Scan(&e.Id, &e.Title, &body, &e.Tags, &e.Created, &e.Public)
	e.Body = template.HTML(body)

	box := packr.New("bliki", "./templates")
	s, err := box.FindString("edit.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl, _ := template.New("edit").Parse(s)
	tmpl.Execute(w, e)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	body := r.FormValue("body")
	tags := r.FormValue("tags")
	public := r.FormValue("public")
	if len(public) == 0 {
		public = "0"
	}
	created := time.Now()
	statement, err := database.Prepare("INSERT INTO entries (title, body, tags, created, public) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec(title, body, tags, created, public)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func newHandler(w http.ResponseWriter, request *http.Request) {
	box := packr.New("bliki", "./templates")
	s, err := box.FindString("new.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl, _ := template.New("new").Parse(s)
	tmpl.Execute(w, nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Query("SELECT id, title, body, tags, created, public FROM entries WHERE public = true order by created desc ")
	if err != nil {
		log.Fatal(err)
	}
	data := IndexData{}
	for rows.Next() {
		e := Entry{}
		var body string
		rows.Scan(&e.Id, &e.Title, &body, &e.Tags, &e.Created, &e.Public)
		e.CreatedText = e.Created.Format(time.RFC1123)
		e.Body = template.HTML(strings.Replace(body, "\r\n", "<br>", -1))
		data.Entries = append(data.Entries, e)
	}

	box := packr.New("bliki", "./templates")
	s, _ := box.FindString("bliki.html")
	tmpl, _ := template.New("bliki").Parse(s)
	tmpl.Execute(w, data)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Query("SELECT id, title, body, tags, created, public FROM entries order by created desc ")
	if err != nil {
		log.Fatal(err)
	}
	data := IndexData{}
	for rows.Next() {
		e := Entry{}
		var body string
		rows.Scan(&e.Id, &e.Title, &body, &e.Tags, &e.Created, &e.Public)
		e.CreatedText = e.Created.Format(time.RFC1123)
		e.Body = template.HTML(strings.Replace(body, "\r\n", "<br>", -1))
		data.Entries = append(data.Entries, e)
	}

	box := packr.New("bliki", "./templates")
	s, _ := box.FindString("admin.html")
	tmpl, _ := template.New("admin").Parse(s)
	tmpl.Execute(w, data)
}

func basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if username, password, ok := r.BasicAuth(); ok {
			if username == os.Getenv("USERNAME") && password == os.Getenv("PASSWORD") {
				next.ServeHTTP(w, r)
				return
			}
		}
		w.Header().Set("WWW-Authenticate", "Basic realm=\"restvoice.org\"")
		w.WriteHeader(http.StatusUnauthorized)
	}
}
