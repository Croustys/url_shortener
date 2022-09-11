package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

var urls map[string]string = make(map[string]string)
var i int = 0
var PORT string = "8080"

type short_request struct {
	Url string `json:"url"`
}

func shorten(url string) string {
	for k, v := range urls {
		if v == url {
			return k
		}
	}

	idx := strconv.Itoa(i)
	i++
	urls[idx] = url
	return idx
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "https://shortrl.netlify.app")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
	w.Header().Set("Content-Type", "application/json")

	var url short_request

	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		log.Println(err)
	}

	if r.Method == "POST" {
		id := shorten(url.Url)

		json_resp, err := json.Marshal(map[string]string{"url": id})
		if err != nil {
			log.Println(err.Error())
		}

		w.Write(json_resp)
	} else if r.Method == "GET" {
		idx := r.URL.Query()["r"][0]
		http.Redirect(w, r, urls[idx], http.StatusSeeOther)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Wrong Request"))
	}
}

func main() {
	http.HandleFunc("/", handler)

	log.Println("Listening on", PORT)
	http.ListenAndServe(":"+PORT, nil)
}
