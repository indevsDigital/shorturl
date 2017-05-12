package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	hashids "github.com/speps/go-hashids"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

func init() {
	// create a sql connect
	db, err = sql.Open("mysql", "root:duncan@tcp(127.0.0.1:3306)/shorturl")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`create table if not exists shorturl (id int NOT NULL AUTO_INCREMENT,dataid varchar(100),short_url varchar(100),long_url varchar(255),primary key(id))`)
	if err != nil {
		log.Fatal(err)
	}
	consolewelcome := `
	~
	~ URU shortner made with love
	~
	`
	fmt.Println(consolewelcome)
}

type myURL struct {
	ID       string `json:"id,omitempty"`
	LongURL  string `json:"longurl,omitempty"`
	ShortURL string `json:"shorturl,omitempty"`
}

func expandEndPoint(w http.ResponseWriter, r *http.Request) {
	selecturl, _ := db.Prepare(`select dataid , short_url, long_url from shorturl where short_url=?`)
	params := r.URL.Query()
	rows := selecturl.QueryRow(params.Get("shorturl"))
	var row myURL
	err = rows.Scan(&row.ID, &row.ShortURL, &row.LongURL)
	if err != nil {

	}

	json.NewEncoder(w).Encode(row)

}

func createEndPoint(w http.ResponseWriter, r *http.Request) {
	var url myURL
	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		log.Println("called")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	selecturl, _ := db.Prepare(`select dataid , short_url, long_url from shorturl where long_url=?`)
	rows := selecturl.QueryRow(url.LongURL)

	var row myURL
	err = rows.Scan(&row.ID, &row.ShortURL, &row.LongURL)
	if row == (myURL{}) {
		hd := hashids.NewData()
		h := hashids.NewWithData(hd)
		now := time.Now()
		url.ID, _ = h.Encode([]int{int(now.Unix())})
		url.ShortURL = "http://localhost:8080/" + url.ID
		insert, _ := db.Prepare(`insert into shorturl (dataid , short_url, long_url) values(?,?,?)`)

		_, err = insert.Exec(url.ID, url.ShortURL, url.LongURL)

	} else {
		url = row
	}
	log.Println("called")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Request-URL", "*")
	json.NewEncoder(w).Encode(url)

}

func rootEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var url myURL

	selecturl, _ := db.Prepare(`select dataid , short_url, long_url from shorturl where dataid=?`)
	rows := selecturl.QueryRow(params["id"])

	err = rows.Scan(&url.ID, &url.ShortURL, &url.LongURL)
	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, url.LongURL, 301)

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/create", createEndPoint).Methods("POST")
	router.HandleFunc("/expand/", expandEndPoint).Methods("GET")
	router.HandleFunc("/{id}", rootEndPoint).Methods("GET")
	fmt.Println("starting server at port http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
