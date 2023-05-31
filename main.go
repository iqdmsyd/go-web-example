package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// process dynamic requests
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprint(w, "Hello Dek")
	// })

	// server static assets
	// fs := http.FileServer(http.Dir("static/"))
	// http.Handle("/static/", http.StripPrefix("/static/", fs))

	// accept connections
	// http.ListenAndServe(":8080", nil)

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello Dek")
	})

	fs := http.FileServer(http.Dir("static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	router.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(w, "You've requested the book: %s on page: %s\n", title, page)
	})

	router.HandleFunc("/books/{title}", CreateBook).Methods("POST")
	router.HandleFunc("/books/{title}", ReadBook).Methods("GET")
	router.HandleFunc("/books/{id}", UpdateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", DeleteBook).Methods("DELETE")

	router.HandleFunc("/records", GetAllRecord).Methods("GET")
	router.HandleFunc("/records/{id}", GetRecordById).Methods("GET")
	router.HandleFunc("/records/{band}/{song}", CreateRecord).Methods("POST")
	router.HandleFunc("/records/{id}", DeleteRecord).Methods("DELETE")

	http.ListenAndServe(":8080", router)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]

	fmt.Fprintf(w, "POST /books/%s\n", title)
}

func ReadBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]

	fmt.Fprintf(w, "GET /books/%s\n", title)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Fprintf(w, "PUT /books/%s\n", id)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Fprintf(w, "DEL /books/%s\n", id)
}

type record struct {
	Id   int32
	Band string
	Song string
}

func GetAllRecord(w http.ResponseWriter, r *http.Request) {
	// connection string
	connStr := "postgres://mydb:mydb@localhost"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// query data
	rows, err := db.Query("SELECT id, band, song FROM records")
	if err != nil {
		errRes := map[string]string{"error": "Data not found"}
		errJSON, _ := json.Marshal(errRes)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(errJSON)
		return
	}
	defer rows.Close()

	// scan query result
	var result []record
	for rows.Next() {
		var r record
		err := rows.Scan(&r.Id, &r.Band, &r.Song)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, r)
	}

	// encode result to json
	resJSON, _ := json.Marshal(result)

	// set response header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

func GetRecordById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	connStr := "postgres://mydb:mydb@localhost"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	var result record
	row := db.QueryRow("SELECT id, band, song FROM records WHERE id = $1", id).Scan(&result.Id, &result.Band, &result.Song)

	if row != nil {
		errRes := map[string]string{"error": "Data not found"}
		errJSON, _ := json.Marshal(errRes)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(errJSON)
		return
	}

	resJSON, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

func CreateRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	band := vars["band"]
	song := vars["song"]

	// connection string
	connStr := "postgres://mydb:mydb@localhost"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	var id int32
	row := db.QueryRow(`INSERT INTO records(band, song)
		VALUES($1, $2) RETURNING id`, band, song).Scan(&id)
	if row != nil {
		log.Fatal(row)
	}

	result := record{Id: id, Band: band, Song: song}
	resJSON, _ := json.Marshal(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resJSON)
}

func DeleteRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	connStr := "postgres://mydb:mydb@localhost"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	var rec record
	row := db.QueryRow("SELECT id, band, song FROM records WHERE id = $1", id).Scan(&rec.Id, &rec.Band, &rec.Song)
	if row != nil {
		errRes := map[string]string{"error": "Data not found"}
		errJSON, _ := json.Marshal(errRes)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(errJSON)
		return
	}

	db.QueryRow("DELETE FROM records WHERE id = $1", id)
	response := map[string]string{"data": ""}
	resJSON, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}
