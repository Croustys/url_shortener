package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

func url_already_exists(w http.ResponseWriter, r *http.Request) {
	json_resp, err := json.Marshal(map[string]string{"status_message": "shortened url already exists"})
	if err != nil {
		log.Println(err.Error())
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Write(json_resp)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var url short_request

		err := json.NewDecoder(r.Body).Decode(&url)
		if err != nil {
			log.Println(err)
		}

		var id string

		for k, v := range urls {
			if v == url.Url {
				id = k
			}
		}

		if url.Custom != "" && id == "" {

			if urls[url.Custom] != "" {
				url_already_exists(w, r)
				return
			}

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
		if urls[idx] != "" {
			http.Redirect(w, r, urls[idx], http.StatusSeeOther)
		}
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/", handler)

	log.Println("Listening on", PORT)
	http.ListenAndServe(":"+PORT, nil)
}
