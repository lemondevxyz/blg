# blg
This is my take on blogs.

## Features
- Creation, Deletion, Modification of posts
- Multi-user posts
- Search
- API

## Config
The config is laid out in file `config.yaml`, where you have three fields to customize from:
```
# domain used to link articles
domain: localhost:8008
# what port to run the http server in
port: 8080
# session key secret, 32 length
secret: H2yTFfUtlUn3Qqb2MiPxr6Rt0F3blSwm
```

## API
When consuming the API, errors will be returned without a response body. only a response code. Here is what they mean:
```
400 - Missing parameters in request
403 - You are not logged in or You are logged in
404 - Post/User doesn't exist
500 - Error in database
```

## Testing
You could test the database through `go test -u`, and the API through postman
