package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "root"
	DB_PASSWORD = "root"
	DB_NAME     = "goplay"
)

type response struct {
	Type    string  `json:"type"`
	Data    []album `json:"data"`
	Message string  `json:"message"`
}

type album struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/albums", getAlbums).Methods("GET")
	router.HandleFunc("/albums", postAlbum).Methods("POST")

	http.Handle("/", router)

	fmt.Println("Running on localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func setupDB() *sql.DB {
	conn := fmt.Sprintf(
		"postgresql://%s:%s@localhost:5432/%s?sslmode=disable",
		DB_USER,
		DB_PASSWORD,
		DB_NAME,
	)
	db, err := sql.Open("postgres", conn)

	checkErr(err)

	return db
}

// Using gorilla/mux
func getAlbums(w http.ResponseWriter, req *http.Request) {
	log.Println("Received get request")
	db := setupDB()

	log.Println("Fetching albums")
	rows, err := db.Query("SELECT * FROM albums")
	defer rows.Close()

	checkErr(err)

	var albums []album

	for rows.Next() {
		var (
			id     int
			title  string
			artist string
		)

		err = rows.Scan(&id, &title, &artist)

		log.Println("Parsing album " + title + " by " + artist)

		checkErr(err)

		albums = append(albums, album{title, artist})
	}

	res := response{Type: "success", Data: albums}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func postAlbum(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	checkErr(err)

	var newAlbum album
	err = json.Unmarshal(body, &newAlbum)
	checkErr(err)

	var res response
	if newAlbum.Artist == "" || newAlbum.Title == "" {
		res = response{Type: "Error", Message: "This method is not supported"}
	} else {
		db := setupDB()

		var lastInsertID int
		err = db.QueryRow(
			"INSERT INTO albums(title, artist) VALUES($1, $2) returning id;",
			newAlbum.Title,
			newAlbum.Artist,
		).Scan(&lastInsertID)

		checkErr(err)

		res = response{Type: "success", Message: "The album has been inserted successfully!"}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// Using the standard library.
// For comparison

var a1 = album{
	Title:  "Sad and Sorrow",
	Artist: "James Blunt",
}

var a2 = album{
	Title:  "Happy and Glad",
	Artist: "Bruno Mars",
}

var localAlbums = []album{a1, a2}

func albumsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var j []byte
	var err error

	if req.Method == "POST" {
		j, err = io.ReadAll(req.Body)
		checkErr(err)

		var newAlbum album
		err = json.Unmarshal(j, &newAlbum)
		checkErr(err)

		localAlbums = append(localAlbums, newAlbum)
	} else if req.Method == "GET" {
		j, err = json.Marshal(localAlbums)
		checkErr(err)
	} else {
		e := response{Type: "Error", Message: "This method is not supported"}
		j, err = json.Marshal(e)
		checkErr(err)
	}

	w.Write(j)
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
