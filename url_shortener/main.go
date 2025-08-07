package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID          string    `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	CreatonDate time.Time `json:"creation_date"`
}

var urlDB = make(map[string]URL)

func generateShortURL(OriginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL)) // converts the OriginalURL string to a byte size
	data := hasher.Sum(nil)
	hash := hex.EncodeToString(data)
	//fmt.Println("hash:", hash)
	return hash[:8]
}

func createURL(originalURL string) string {
	shortURL := generateShortURL(originalURL)
	id := shortURL
	urlDB[id] = URL{
		ID:          id,
		OriginalURL: originalURL,
		ShortURL: shortURL,
		CreatonDate: time.Now(),
	}
	return shortURL
}

func getURL(id string) (URL, error) {
	url, ok := urlDB[id]

	if !ok {
		return URL{}, errors.New("URL not found")
	}

	return url, nil
}

func RootPageURL(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w,"Hello EveryOne")
}

func ShortURLHandler(w http.ResponseWriter, r *http.Request){
	var data struct{
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil{
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	shortURL_ := createURL(data.URL)
	//fmt.Fprintf(w, shortURL)
	response := struct{
		ShortURL string `json:"short_url"`
	}{ShortURL: shortURL_}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func redirectURLHandler(w http.ResponseWriter, r *http.Request){
	id := r.URL.Path[len("/redirect/"):]
	url, err := getURL(id)
	if err != nil{
		http.Error(w, "Invalid request", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}
func main() {

	// Register the handler function to handle all request to the root URL ("/")
	http.HandleFunc("/", RootPageURL)
	http.HandleFunc("/shorten", ShortURLHandler)
	http.HandleFunc("/redirect/", redirectURLHandler)

	fmt.Println("Starting server on port 3000...")
	//start the http server on port 3000
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error on starting the server:", err)
	}
}
