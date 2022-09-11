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
	Url    string `json:"url"`
	Custom string `json:"custom"`
}

func shorten(url string) string {
	idx := strconv.Itoa(i)
	i++
	urls[idx] = url
	return idx
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
	w.Header().Set("Content-Type", "application/json")

	var url short_request

	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		log.Println(err)
	}

	if r.Method == "POST" {
		var id string

		for k, v := range urls {
			if v == url.Url {
				id = k
			}
		}

		if url.Custom != "" && id == "" {
			id = url.Custom
			urls[id] = url.Url
		} else if id == "" {
			id = shorten(url.Url)
		}

		json_resp, err := json.Marshal(map[string]string{"url": id})
		if err != nil {
			log.Println(err.Error())
		}

		w.Write(json_resp)
		return
	}
	if r.Method == "GET" {
		idx := r.URL.Query()["r"][0]
		http.Redirect(w, r, urls[idx], http.StatusSeeOther)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/", handler)

	log.Println("Listening on", PORT)
	http.ListenAndServe(":"+PORT, nil)
}
