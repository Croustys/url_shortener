package main

import "testing"

func TestGet_slug(t *testing.T) {
	got := get_slug("/25e33860-e9b6-4d9b-8aaa-9b8a881aa6bc").String()
	want := "25e33860-e9b6-4d9b-8aaa-9b8a881aa6bc"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestCreateUrl(t *testing.T) {
	url := "https://www.google.com/"
	created_url_uuid := createUrl(url)
	redirect_url := urls[created_url_uuid]

	if redirect_url != url {
		t.Errorf("Created UUID does not match its associated url")
	}
}
