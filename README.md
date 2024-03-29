# Cache Money

## Usage

For use with GitHub Actions, add a step to install the caching tool and wrap the desired command to store and restore the cache:

```
- name: Install dependencies
  run: |
    go get github.com/RobotsAndPencils/cache-money-client
    /home/runner/go/bin/cache-money-client restore
    npm install
    /home/runner/go/bin/cache-money-client store
  env:
    CACHE_KEY: v1-{{ checksum "package-lock.json" }}-${{ runner.os }}
    CACHE_PATH: node_modules
    ENDPOINT: https://{{host}}/api
    TOKEN: ${{ secrets.TOKEN }}
```

TODO: release static binaries for each operating system

## Server

Cache Money client communicates with a REST-based server with the following API.

### Authorization header for all requests

Header:

- `Authorization: {{token}}`

### Check if key exists in the cache

`HEAD {{endpoint}}/{{key}}`

Status:

- 204 No content
- 404 Not found
- 500 Internal server error

### Upload a file to the cache

`PUT {{endpoint}}/{{key}}`

Headers:

- `Content-Type: application/zip`
- `Content-Length: 123`

Status:

- 200 OK
- 500 Internal server error

### Download a file from the cache

`GET {{endpoint}}/{{key}}`

Status:

- 200 OK
- 404 Not found
- 500 Internal server error
