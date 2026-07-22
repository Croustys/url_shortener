package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/syumai/workers"
	"github.com/syumai/workers/cloudflare/kv"
)

const KV_BINDING string = "URLS"

const (
	slug_prefix string = "s:"
	url_prefix  string = "u:"
)

const slug_alphabet string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const slug_length int = 6
const slug_attempts int = 5

var custom_pattern = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)

var errNoSlug = errors.New("no unused slug found")

type short_request struct {
	Url    string `json:"url"`
	Custom string `json:"custom"`
}

func to_json(key string, value string) []byte {
	json, err := json.Marshal(map[string]string{key: value})
	if err != nil {
		log.Println(err.Error())
	}

	return json
}

func write_json(w http.ResponseWriter, status int, key string, value string) {
	w.WriteHeader(status)
	w.Write(to_json(key, value))
}

func lookup(ns *kv.Namespace, key string) (string, bool, error) {
	value, err := ns.GetString(key, nil)
	if err != nil {
		return "", false, err
	}

	if value == "<null>" || value == "<undefined>" {
		return "", false, nil
	}

	return value, true, nil
}

func url_key(target string) string {
	sum := sha256.Sum256([]byte(target))
	return url_prefix + hex.EncodeToString(sum[:])
}

func random_slug() (string, error) {
	buf := make([]byte, slug_length)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	for i, b := range buf {
		buf[i] = slug_alphabet[int(b)%len(slug_alphabet)]
	}

	return string(buf), nil
}

func valid_target(target string) bool {
	parsed, err := url.Parse(target)
	if err != nil {
		return false
	}

	return (parsed.Scheme == "http" || parsed.Scheme == "https") && parsed.Host != ""
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	urls, err := kv.NewNamespace(KV_BINDING)
	if err != nil {
		log.Println(err)
		write_json(w, http.StatusInternalServerError, "status_message", "KV namespace is unavailable")
		return
	}

	switch r.Method {
	case "POST":
		handleCreation(w, r, urls)
	case "GET":
		handleRedirect(w, r, urls)
	}
}

func handleCreation(w http.ResponseWriter, r *http.Request, urls *kv.Namespace) {
	var request short_request

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println(err)
		write_json(w, http.StatusBadRequest, "status_message", "Request body is not valid JSON")
		return
	}

	if !valid_target(request.Url) {
		write_json(w, http.StatusBadRequest, "status_message", "Url must be an absolute http or https URL")
		return
	}

	if request.Custom != "" {
		handleCustom(w, request, urls)
		return
	}

	if existing, found, err := lookup(urls, url_key(request.Url)); err != nil {
		log.Println(err)
		write_json(w, http.StatusInternalServerError, "status_message", "Failed to read from KV")
		return
	} else if found {
		write_json(w, http.StatusOK, "url", existing)
		return
	}

	slug, err := unused_slug(urls)
	if err != nil {
		log.Println(err)
		write_json(w, http.StatusInternalServerError, "status_message", "Failed to allocate a slug")
		return
	}

	if err := urls.PutString(slug_prefix+slug, request.Url, nil); err != nil {
		log.Println(err)
		write_json(w, http.StatusInternalServerError, "status_message", "Failed to write to KV")
		return
	}

	if err := urls.PutString(url_key(request.Url), slug, nil); err != nil {
		log.Println(err)
	}

	write_json(w, http.StatusOK, "url", slug)
}

func handleCustom(w http.ResponseWriter, request short_request, urls *kv.Namespace) {
	if !custom_pattern.MatchString(request.Custom) {
		write_json(w, http.StatusBadRequest, "status_message", "Custom url may only contain letters, digits, hyphens and underscores")
		return
	}

	_, taken, err := lookup(urls, slug_prefix+request.Custom)
	if err != nil {
		log.Println(err)
		write_json(w, http.StatusInternalServerError, "status_message", "Failed to read from KV")
		return
	}

	if taken {
		write_json(w, http.StatusConflict, "status_message", "Custom url already exists")
		return
	}

	if err := urls.PutString(slug_prefix+request.Custom, request.Url, nil); err != nil {
		log.Println(err)
		write_json(w, http.StatusInternalServerError, "status_message", "Failed to write to KV")
		return
	}

	if _, found, err := lookup(urls, url_key(request.Url)); err == nil && !found {
		if err := urls.PutString(url_key(request.Url), request.Custom, nil); err != nil {
			log.Println(err)
		}
	}

	write_json(w, http.StatusOK, "url", request.Custom)
}

func unused_slug(urls *kv.Namespace) (string, error) {
	var err error

	for attempt := 0; attempt < slug_attempts; attempt++ {
		var slug string

		slug, err = random_slug()
		if err != nil {
			continue
		}

		_, taken, lookup_err := lookup(urls, slug_prefix+slug)
		if lookup_err != nil {
			err = lookup_err
			continue
		}

		if !taken {
			return slug, nil
		}
	}

	if err != nil {
		return "", err
	}

	return "", errNoSlug
}

func handleRedirect(w http.ResponseWriter, r *http.Request, urls *kv.Namespace) {
	idx := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")[0]

	if idx == "" {
		idx = r.URL.Query().Get("r")
	}

	if idx != "" {
		redirect_url, found, err := lookup(urls, slug_prefix+idx)
		if err != nil {
			log.Println(err)
			write_json(w, http.StatusInternalServerError, "status_message", "Failed to read from KV")
			return
		}

		if found {
			http.Redirect(w, r, redirect_url, http.StatusSeeOther)
			return
		}
	}

	write_json(w, http.StatusNotFound, "status_message", "Url associated to the requested shortUrl does not exist")
}

func main() {
	http.HandleFunc("/", handler)

	workers.Serve(nil)
}
