# blg
This is my take on blogs. [screenshots](SCREENSHOT.md)

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
500 - Error in database(could mean that email already exists or post already exists)
```
All parameters are Captialized, and need to be provided in a `application/www-x-form-urlencoded`-type request

### API-routes
```
# User
GET:    /api/user - Returns the current user's information
POST:   /api/user - Creates a new user(need to provide Email, Username, Password as parameters)
PATCH:  /api/user - Updates the current user's information(provide one of Firstname, Lastname, Description, Email, Username or Password)
DELETE: /api/user - Delete's the current user

# Auth
POST: /api/auth/login   - Logs in(need to provide UserID(which is basically username or email) and Password)
POST: /api/auth/logout  - Logs out

# Post:
GET:    /api/post/:title  -  Returns a post by it's title
POST:   /api/post         -  Creates a new post(need to provide Title, Description, Content)
PATCH:  /api/post/:title  -  Updates a post(provide one of, Title, Description, Content, Public(bool))
DELETE: /api/post/:title  -  Deletes a post

GET: /api/post-search/:field - Searches for a post by that a name or description similar to :field
```

## Testing
You could test the database through `go test -u`, and the API through postman
