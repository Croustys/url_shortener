package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var urls map[string]string = make(map[string]string)
var i int = 0

const PORT string = "3600"

type short_request struct {
	Url    string `json:"url"`
	Custom string `json:"custom"`
}

func shorten(url string) string {
	idx := strconv.Itoa(i)
	i++
	urls[idx] = url
	return idx
}

func to_json(key string, value string) []byte {
	json, err := json.Marshal(map[string]string{key: value})
	if err != nil {
		log.Println(err.Error())
	}

	return json
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		handleCreation(w, r)
	} else if r.Method == "GET" {
		handleRedirect(w, r)
	}
}

func handleCreation(w http.ResponseWriter, r *http.Request) {
	var url short_request

	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		log.Println(err)
	}

	if _, exists := urls[url.Custom]; exists {
		json := to_json("status_message", "Custom url already exists")
		w.WriteHeader(http.StatusConflict)
		w.Write(json)
		return
	}

	if url.Custom != "" {
		custom_url := url.Custom
		urls[custom_url] = url.Url

		json := to_json("url", custom_url)
		w.WriteHeader(http.StatusOK)
		w.Write(json)
		return
	}

	for k, v := range urls {
		if v == url.Url {
			json := to_json("url", k)
			w.WriteHeader(http.StatusOK)
			w.Write(json)
			return
		}
	}

	short_url := shorten(url.Url)
	json := to_json("url", short_url)
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}
func handleRedirect(w http.ResponseWriter, r *http.Request) {
	subRoute := r.URL.Path[1:]
	idx := strings.Split(subRoute, "/")[0]

	if redirect_url, exists := urls[idx]; exists {
		http.Redirect(w, r, redirect_url, http.StatusSeeOther)
		return
	}

	json := to_json("status_message", "Url associated to the requested shortUrl does not exist")
	w.WriteHeader(http.StatusNotFound)
	w.Write(json)
}

func main() {
	http.HandleFunc("/", handler)

	log.Println("Listening on", PORT)
	http.ListenAndServe("localhost:"+PORT, nil)
}
