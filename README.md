# URL Shortener Written in Go

## Usage

`POST` `/api/shorten` with the request body containing a JSON with the URL we'd like to shorten. `Returns` the slug you can be redirected with <br /><br />
`GET` `/{slug}` redirects you to the url you've given in the post request<br />

### Examples

`POST` `api/shorten` `Body`: `'URL': 'https://www.google.com/'`<br />
`Returns` a random UUID slug `(ea5ec738-e36c-4ed4-be42-73df0fdd5e2f)`

`GET` `/ea5ec738-e36c-4ed4-be42-73df0fdd5e2f` redirects you to `'https://www.google.com/'`

### Hosted on Digital Ocean

### Frontend

Currently in development using Solidjs<br />
[Git Repo](https://github.com/Croustys/url-shortener-frontend)
