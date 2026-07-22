# URL Shortener Written in Go

Runs on [Cloudflare Workers](https://developers.cloudflare.com/workers/) as a WebAssembly
module, with [Workers KV](https://developers.cloudflare.com/kv/) for storage. The Go code is
compiled with `GOOS=js GOARCH=wasm` and served through
[`syumai/workers`](https://github.com/syumai/workers).

## Usage

### `POST`

- `/` with the request body containing a JSON with the URL we'd like to shorten. `Returns` the slug you can be redirected with <br /><br />

### `GET`

- `/{slug}` or `/?r={slug}` redirects you to the url you've given in the post request<br />
- You can pass a `'custom'` field as well in the JSON body which will return a Custom redirect URL

## Examples

`POST` `/` `Body`: `{'url': 'https://www.google.com/'}`
`Returns` a random 6-character slug `on4Kl7`

`POST` `/` `Body`: `{'url': 'https://www.google.com/', 'custom': 'ggl'}`
`Returns` `ggl`

`GET` `/on4Kl7` redirects you to `'https://www.google.com/'`
`GET` `/?r=ggl` redirects you to `'https://www.google.com/'`

Submitting the same URL twice returns the slug it was first given. A custom slug that is
already taken returns `409`. A URL that is not absolute `http`/`https` returns `400`.

## Requirements

- Node.js
- Go 1.24 or later

## Setup

```console
npm install
npx wrangler kv namespace create URLS
```

Put the printed namespace id into the `kv_namespaces` entry in `wrangler.jsonc`.

## Development

```console
npm start          # wrangler dev, KV simulated locally under .wrangler/
npm run build      # build build/app.wasm + the JS shim
npm run deploy     # deploy to Cloudflare
```

Local KV state persists in `.wrangler/state` between runs. To start clean, delete it.

### Editor setup

The whole module is `js/wasm`-only, so a host-platform build fails on `syscall/js`. Point
your language server at the target platform, e.g. for VS Code:

```jsonc
"gopls": { "build.env": { "GOOS": "js", "GOARCH": "wasm" } }
```

## Storage layout

A single KV namespace holds two kinds of keys:

- `s:{slug}` → target URL, read on redirect
- `u:{sha256(url)}` → slug, the reverse index that makes repeated submissions return the
  same slug without scanning the namespace

Note that KV is eventually consistent: a freshly created slug may take a few seconds to
resolve from an edge location other than the one that wrote it.

## Bundle size

The standard Go toolchain produces a ~6.5MB `app.wasm`, ~1.8MB compressed — under the 3MB
free-plan limit. If added dependencies push it over, switch to the
[TinyGo template](https://github.com/syumai/workers/tree/main/_templates/cloudflare/worker-tinygo)
build (`workers-assets-gen -mode=tinygo` plus a `tinygo build` step).

### Hosted on ~~[Fly.io](https://fly.io/)~~ ~~[Oracle](https://www.oracle.com/id/)~~ [Cloudflare](https://workers.cloudflare.com/)

# Production domain

[`links.barabasakos.hu`](https://links.barabasakos.hu)

### Frontend

Currently in development using Solidjs<br />
[Git Repo](https://github.com/Croustys/url-shortener-frontend)
