package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

var urls map[uuid.UUID]string = make(map[uuid.UUID]string)
var redirect_path string = "/"
var PORT int = 8080

type short_request struct {
	Url string
}

func shorten(w http.ResponseWriter, req *http.Request) {
	var url short_request
	err := json.NewDecoder(req.Body).Decode(&url)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(url.Url)
	new_uuid := createUrl(url.Url)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make(map[string]string)
	resp["url"] = new_uuid.String()
	json_resp, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json_resp)
}

func redirect(w http.ResponseWriter, req *http.Request) {
	strUrl := req.URL.String()
	slug := get_slug(strUrl)

	http.Redirect(w, req, urls[slug], http.StatusSeeOther)
}

func createUrl(url string) uuid.UUID {
	uuid := uuid.New()
	urls[uuid] = url
	return uuid
}

func get_slug(input string) uuid.UUID {
	if len(input) <= 1 {
		return uuid.UUID{}
	}
	uuid, err := uuid.Parse(input[1:])
	if err != nil {
		fmt.Println(err.Error())
	}
	return uuid
}

func main() {
	http.HandleFunc("/api/shorten", shorten)
	http.HandleFunc(redirect_path, redirect)

	http.ListenAndServe(":"+strconv.Itoa(PORT), nil)
	fmt.Println("Server running")
}
