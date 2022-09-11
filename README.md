# URL Shortener Written in Go

## Usage

### `POST`

- `/` with the request body containing a JSON with the URL we'd like to shorten. `Returns` the slug you can be redirected with <br /><br />

### `GET`

- `/?r={slug}` redirects you to the url you've given in the post request<br />
- You can pass a `'custom'` field as well in the JSON body which will return a Custom redirect URL

## Examples

`POST` `/` `Body`: `'url': 'https://www.google.com/'`
`Returns` a random number slug `12`

`POST` `/` `Body`: `{'url': 'https://www.google.com/', 'custom": 'ggl'}`
`Returns` `ggl`

`GET` `/?r=12` redirects you to `'https://www.google.com/'`
`GET` `/?r=ggl` redirects you to `'https://www.google.com/'`

### Hosted on [Fly.io](https://fly.io/)

# Production domain

[`links.barabasakos.hu`](https://links.barabasakos.hu)

### Frontend

Currently in development using Solidjs<br />
[Git Repo](https://github.com/Croustys/url-shortener-frontend)
