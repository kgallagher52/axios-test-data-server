package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

type singleItem struct {
	ID     int    `json-:"id"`
	Total  string `json-:"sub_total"`
	Name   string `json-:"item_name"`
	Handle string `json-:"item_handle"`
	SKU    string `json-:"item_sku"`
	Price  string `json-:"item_price"`
	Image  string `json-:"item_image"`
}
type serverError struct {
	Message string
	Status  int
}

var db *sql.DB

func main() {
	pgURL, err := pq.ParseURL(os.Getenv("ELEPHANTSQL_URL"))
	if err != nil {
		log.Fatal(err)
	}

	db, err = sql.Open("postgres", pgURL)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()

	r := mux.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					switch t := err.(type) {
					case string:
						http.Error(w, t, http.StatusInternalServerError)
					case error:
						http.Error(w, t.Error(), http.StatusInternalServerError)
					case serverError:
						http.Error(w, t.Message, t.Status)
					default:
						http.Error(w, "unknown error", http.StatusInternalServerError)
					}
				}
			}()
			next.ServeHTTP(w, r)
		})
	})

	r.HandleFunc("/weboostCartData", getTestCartData).Methods("GET")

	// Telling the server what it accepts
	methods := handlers.AllowedMethods([]string{"GET", "PUT", "POST", "DELETE"})
	headers := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	origins := handlers.AllowedOrigins([]string{
		"http://localhost:3000", "http://localhost:3000/cart", "http://localhost:3000/", "localhost:3000", "http://www.axiostestdataserver.com/weboostCartData"})

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT not set")
	}

	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(methods, origins, headers)(r)))
}

func respondWithError(w http.ResponseWriter, status int, error serverError) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
}

func responseJSON(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

func getTestCartData(w http.ResponseWriter, r *http.Request) {
	var items []singleItem
	var error serverError

	rows, err := db.Query("SELECT * FROM test_items_weboost_cart")
	if err != nil {
		error.Message = fmt.Sprintf(`Server Error %s`, err)
		respondWithError(w, http.StatusInternalServerError, error)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var total string
		var name string
		var handle string
		var sku string
		var price string
		var image string

		err = rows.Scan(&id, &total, &name, &handle, &sku, &price, &image)
		if err != nil {
			error.Message = fmt.Sprintf(`Scan Error %s`, err)
			respondWithError(w, http.StatusInternalServerError, error)
			return
		}

		a := singleItem{ID: id, Total: total, Name: name, Handle: handle, SKU: sku, Price: price, Image: image}
		items = append(items, a)
	}

	err = rows.Err()
	if err != nil {
		error.Message = fmt.Sprintf(`Rows Error %s`, err)
		respondWithError(w, http.StatusInternalServerError, error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	responseJSON(w, items)
}
