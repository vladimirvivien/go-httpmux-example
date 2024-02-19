# go-httpmux-example

Simple exmaple to show how to use of new http.ServeMux router in Go v1.22.x.

## Pre-requisite
* First, update your Go tools to version v1.22.0 or higher
* If you are planning to use the new enhance routing feaure in an existing Go module/project, **make sure to change the Go version in your go.mod file**  to `go 1.22.0` or whatever future version you are using (to avoid frustration trying to figure out why the feature is not working).


## Enhanced Pattern Matching for Go 1.22 http.ServeMux

Go's standard library has long provided a powerful and efficient way to route web requests with the http package. Go 1.22 introduced a new pattern matching feature that makes HTTP request routing more flexible.

* As a convenience, you can embed allowed HTTP methods (e.g., "GET", "POST") directly in the route along with the path. 
* Path patterns can wildcards with named variables denoted using curly braces. 
* Mechanism to extract path variables from the URL and made available to the handler.

## Benefits
* Combining the allowed HTTP method along with the path provides a clean, expressive, and RESTful routing patterns.
* Support for path variables reduces the logic necessary to extract values from the URL within your handler functions.


## Example
The code is a simple demostration of the new feature:

```go
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", getTasks)
	mux.HandleFunc("GET /tasks/{id}", getTask)
	mux.HandleFunc("POST /tasks/create", postTask)
	fmt.Println("Starting server on port 8080")
	http.ListenAndServe(":8080", mux)
}
```

Method call `mux.HandleFunc("GET /tasks", getTasks)` does two things:
* It specifies that the request will only work for a GET request.
* And will only work for request with paths that start with `/tasks`

As mentioned, the request pattern can include a request path variable: `mux.HandleFunc("GET /tasks/{id}", getTask)`. In your handler, you can now use the new `Request.PathValue` method to extract the path variable. 

For instance, `curl http://localhost:8080/tasks/four` would trigger the following handler:

```go
func getTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	lock.RLock()
	task, ok := tasks[idStr]
	lock.RUnlock()

	if !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Error encountered", http.StatusInternalServerError)
	}
}
```

## References
You can read more about the new `http.ServeMux` routing enhancements to get a better picture of its capabilities:

* Routing enhancements for Go v1.22.0 - [Go blog](https://tip.golang.org/blog/routing-enhancements)
* net/http.ServMux source code [documentation](https://pkg.go.dev/net/http@master#ServeMux)
