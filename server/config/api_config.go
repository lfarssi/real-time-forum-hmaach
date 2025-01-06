package config

var APIs = map[string]string{
	"GET  /":                "List all available routes and their descriptions.",
	"GET  /posts":           "Retrieve a list of all posts.",
	"GET  /posts/{id}":      "Retrieve details of a specific post. Parameters: {id} (integer) - the ID of the post.",
	"POST /posts/create":    "Create a new post. Parameters: title (string), content (string), authorID (integer).",
	"POST /comments/create": "Create a new comment on a post. Parameters: postID (integer), content (string), authorID (integer).",
	"POST /register":        "Register a new user. Parameters: nickname (string), email (string), password (string).",
	"POST /login":           "Log in an existing user. Parameters: email (string), password (string).",
	"POST /logout":          "Log out the current user. No parameters required.",
}
