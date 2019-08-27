# Cache Money

## Usage

For use with GitHub Actions, add a step to install the caching tool and wrap the desired command to store and restore the cache:

```
- name: Install dependencies
  run: |
    go get github.com/RobotsAndPencils/cache-money-client
    /home/runner/go/bin/cache-money-client restore
    # npm install
    /home/runner/go/bin/cache-money-client store
  env:
    cache_key: v1-{{ checksum "package-lock.json" }}
    cache_path: node_modules
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
