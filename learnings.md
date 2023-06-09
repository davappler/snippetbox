## Fixed Path and Subtree Patterns

- Go’s servemux supports two different types of URL patterns: fixed paths and subtree paths. Fixed paths don’t end with a trailing slash, whereas subtree paths do end with a trailing slash.

- Our two new patterns — "/snippet" and "/snippet/create" — are both examples of fixed paths. In Go’s servemux, fixed path patterns like these are only matched (and the corresponding handler called) when the request URL path exactly matches the fixed path.
- That path is only triggered if the whole URL is an exact match of the path
- In contrast, our pattern "/" is an example of a subtree path (because it ends in a trailing slash). Another example would be something like "/static/".
- Subtree path patterns are matched/triggered whenever the start of a request URL path matches the subtree path.
- If it helps your understanding, you can think of subtree paths as acting a bit like they have a wildcard at the end, like "/**" or "/static/**".
- This helps explain why the "/" pattern is acting like a catch-all. The pattern essentially means match a single slash, followed by anything (or nothing at all) .

## Restricting the Root URL Pattern

So what if you don’t want the "/" pattern to act like a catch-all?
For instance, in the application we’re building we want the home page to be displayed if — and only if — the request URL path exactly matches "/". Otherwise, we want the user to receive a 404 page not found response.
It’s not possible to change the behavior of Go’s servemux to do this, but you can include a simple check in the home hander which ultimately has the same effect:

```
package main
...
func home(w http.ResponseWriter, r *http.Request) {
    // Check if the current request URL path exactly matches "/".
    // If it doesn't, use // the http.NotFound() function to send a 404 response to the client.
    // Importantly, we then return from the handler.
    // If we don't return the handler
    // would keep executing and also write the "Hello from SnippetBox" message.


    if r.URL.Path != "/" {
    http.NotFound(w, r)
    return
    }
    w.Write([]byte("Hello from Snippetbox"))
}


```

## The DefaultServeMux (Usage not recommended)

- If you’ve been working with Go for a while you might have come across the http.Handle() and http.HandleFunc() functions. These allow you to register routes without declaring a servemux, like this:

```
func main() {
    http.HandleFunc("/", home)
    http.HandleFunc("/snippet", showSnippet)
    http.HandleFunc("/snippet/create", createSnippet)

    log.Println("Starting server on :4000")
    err := http.ListenAndServe(":4000", nil)
    log.Fatal(err)
}
```

- Behind the scenes, these functions register their routes with something called the DefaultServeMux.
- There’s nothing special about this — it’s just regular servemux like we’ve already been using, but which is initialized by default and stored in a net/http global variable.
- It is not recommended to use it because DefaultServeMux is a global variable, any package can access it and register a route — including any third-party packages that your application imports. If one of those third-party packages is compromised, they could use DefaultServeMux to expose a malicious handler to the web.

## Servemux Features and Quirks

- In Go’s servemux, longer URL patterns always take precedence over shorter ones. So, if a servemux contains multiple patterns which match a request, it will always dispatch the request to the handler corresponding to the longest pattern. This has the nice side- effect that you can register patterns in any order and it won’t change how the servemux behaves.
- Request URL paths are automatically sanitized. If the request path contains any . or .. elements or repeated slashes, it will automatically redirect the user to an equivalent clean URL. For example, if a user makes a request to /foo/bar/..//baz they will automatically be sent a 301 Permanent Redirect to /foo/baz instead.
- If a subtree path has been registered and a request is received for that subtree path without a trailing slash, then the user will automatically be sent a
  301 Permanent Redirect to the subtree path with the slash added. For example, if you have registered the subtree path /foo/, then any request to /foo will be redirected to /foo/ .

## What Are Canonical URLs?

- The canonical URL is the URL for the master copy of a page when you have duplicate versions of that page.
- Since the duplicate pages are on different URLs, you can set a canonical URL so Google knows that that page is the original. Or the most representative.
- Canonical URLs help Google better understand your site. And that could help your site rank higher in search results.

## Host Name Matching

- It’s possible to include host names in your URL patterns. This can be useful when you want to redirect all HTTP requests to a canonical URL, or if your application is acting as the back end for multiple sites or services.
- When it comes to pattern matching, any host-specific patterns will be checked first and if there is a match the request will be dispatched to the corresponding handler. Only when there isn’t a host-specific match found will the non-host specific patterns also be checked.

```
mux := http.NewServeMux()
mux.HandleFunc("foo.example.org/", fooHandler)
mux.HandleFunc("bar.example.org/", barHandler)
mux.HandleFunc("/baz", bazHandler)
```

## Specifying POST request

- It’s only possible to call `w.WriteHeader()` once per response, and after the status code has been written it can’t be changed
- If you try to call w.WriteHeader() a second time Go will log a warning message.
- It should be called before `w.Write()` otherwise when we start using `w.Write()`, go will assume that status is 200 ok.
- Whenever the status has to be different than 200 ok, we should specify if with the help of `w.WriteHeader()`

```

func createSnippet(w http.ResponseWriter, r *http.Request) {
    // Use r.Method to check whether the request is using POST or not.
    // If it's not, use the w.WriteHeader() method to send a 405 status code  and
    // the w.Write() method to write a "Method Not Allowed" response body. We
    // then return from the function so that the subsequent code is not executed.

    if r.Method != "POST" {
        w.WriteHeader(405)
        w.Write([]byte("Method Not Allowed"))
        return
    }
    w.Write([]byte("Create a new snippet..."))
}

```

## Customize the response header

- Another improvement we can make is to include an Allow: POST header with every 405 Method Not Allowed response to let the user know which request methods are supported for that particular URL.

- We should use `w.Header().Set("Allow", "POST")` before `w.WriteHeader(405)` or `w.Write([]byte("Method Not Allowed"))` otherwise the changes will be not reflected to the response that the user will receive.

```
if r.Method != "POST" {
    // Use the Header().Set() method to add an
    // 'Allow: POST' header to the // response header map.
    // The first parameter is the header name, and
    // the second parameter is the header value.

    w.Header().Set("Allow", "POST")
    w.WriteHeader(405)
    w.Write([]byte("Method Not Allowed"))
    return
}

```

## http.Error

- We can use this inbuilt function to call the `w.Write()` and `w.WriteHeader`
- `http.Error(w, "Method Not Allowed", 405)` => this will call them with the given parameters.

# Manipulating the Header Map

- `w.Header().Set()` is used to set a header in the response header map.
- Jut like that we have more functions that can modify the header map.


```

// Set a new cache-control header. If an existing "Cache-Control" header exists 
// it will be overwritten.
w.Header().Set("Cache-Control", "public, max-age=31536000")


// In contrast, the Add() method appends a new "Cache-Control" header and can 
// be called multiple times.
w.Header().Add("Cache-Control", "public")
w.Header().Add("Cache-Control", "max-age=31536000")

// Delete all values for the "Cache-Control" header.
w.Header().Del("Cache-Control")

// Retrieve the first value for the "Cache-Control" header.
w.Header().Get("Cache-Control")


```



- Go will automatically set three system-generated headers for you: `Date` and `Content-Length` and `Content-Type`.
- GO can’t distinguish JSON from plain text. And, by default,JSON responses will be sent with a `Content-Type: text/plain; charset=utf-8` header.
- We can manually set it for JSON like this 

```
w.Header().Set("Content-Type", "application/json") 
w.Write([]byte(`{"name":"Alex"}`))

```



## Header Canonicalization

- By default when we add a header in the header map, the `key` of the key-value pair will by default be canonicalized.
- That means if we add a key like this => `w.Header().Add("foo-bar", "public")` , here key is `foo-bar` it will be stored in canonicalized manner and will become => `Foo-Bar`.
- It means the first letter and the letters after the dash all of them will be capitalized.

- We can avoid this default behaviour by manually adding a key to the object that `w.Header()` returns => `w.Header()["X-XSS-Protection"] = []string{"1; mode=block"}`
- Here we are adding `"X-XSS-Protection"` manually so it will remain as it is exactly. 
- On the right hand side we have `[]string{"1; mode=block"}`, this is an array of string and the first value will be => `"1; mode=block"`


## Suppressing System-Generated Headers
- The Del() method doesn’t remove system-generated headers. To suppress these, you need to access the underlying header map directly and set the value to nil. If you want to suppress the Date header, for example, you need to write:

```
w.Header()["Date"] = nil
```
