package main

import (
	"fmt"
	"io"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func insertHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(io.LimitReader(r.Body, 2048))
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	if len(body) == 0 {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	shortUrl, err := InsertUrl(body)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error creating URL", http.StatusInternalServerError)
	} else {
		w.Write(shortUrl)
	}
}

func fetchHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	if len(body) == 0 {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}
	url, err := GetLongUrl(body)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error fetching URL", http.StatusInternalServerError)
		return
	}

	if url == nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Write(url)
	}
}

func main() {
	err := SetupDBConnection()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer CloseDB()

	http.HandleFunc("/set", insertHandler)
	http.HandleFunc("/get", fetchHandler)
	fmt.Printf("Starting server on port :8080\n")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error")
	}
}
