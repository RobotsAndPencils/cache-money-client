# Cache Money

## Usage

For use with GitHub Actions, add a step to install the caching tool:

```
- name: Install caching tool
  run: go get github.com/RobotsAndPencils/cache-money-client
```

And wrap installation with a commands to store and restore the cache:

```
- name: Install dependencies
  run: |
    cache restore
    npm install
    cache store
  with:
    key: v1-{{ checksum "package-lock.json" }}
    path: node_modules
    token: ${{ secrets.TOKEN }}
    endpoint: https://{{host}}/api
```

## Server

Cache Money client communicates with a REST-based server with the following APIs.

### Authorization header for all requests

Header: `Authorization: {{token}}`

### Check if key exists in the cache

`HEAD {{endpoint}}/{{key}}`

- 204 No content
- 404 Not found
- 500 Internal server error

### Upload a file to the cache

`PUT {{endpoint}}/{{key}}`

Header: `Content-Type: application/zip`

- 200 OK
- 500 Internal server error

### Download a file from the cache

`GET {{endpoint}}/{{key}}`

- 200 OK
- 404 Not found
- 500 Internal server error
