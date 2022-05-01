package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

var urls map[uuid.UUID]string = make(map[uuid.UUID]string)
var redirect_path string = "/"
var PORT int = 7823

func shorten(w http.ResponseWriter, req *http.Request) {
	var url map[string]string
	err := json.NewDecoder(req.Body).Decode(&url)
	if err != nil {
		fmt.Println(err)
	}
	new_uuid := createUrl(url["url"])

	json_response, err := json.Marshal(fmt.Sprintf("localhost:%d%s%s", PORT, redirect_path, new_uuid))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json_response)
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
func get_slug(url string) uuid.UUID {
	arr := strings.Split(url, "")
	slicedArr := arr[len(redirect_path):]
	slicedStr := strings.Join(slicedArr[:], "")
	slug, err := uuid.Parse(slicedStr)
	if err != nil {
		panic(err.Error())
	}
	return slug
}

func main() {
	http.HandleFunc("/api/shorten", shorten)
	http.HandleFunc(redirect_path, redirect)

	http.ListenAndServe(":"+strconv.Itoa(PORT), nil)
	fmt.Println("Server running")
}
