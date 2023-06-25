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

## URL Query Strings

- to extract id from th url `/snippet?id=1` => `r.URL.Query().Get("id")`

```
package main
import (
"fmt" // New import
"log"
"net/http"
"strconv" // New import
)
...
func showSnippet(w http.ResponseWriter, r *http.Request) {
	// Extract the value of the id parameter from the query string and try to
	// convert it to an integer using the strconv.Atoi() function. If it can't
	// be converted to an integer, or the value is less than 1, we return a 404 page
    // not found response.

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	// Use the fmt.Fprintf() function to interpolate the id value with our response
    // and write it to the http.ResponseWriter.
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

```

## The io.Writer Interface

- If you take a look at the documentation for the fmt.Fprintf() function you’ll notice that it takes an io.Writer as the first paramete

- `func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error)`
- but above in the code snippet we passed it our http.ResponseWriter object instead — and it worked fine.

- We’re able to do this because the io.Writer type is an interface, and the http.ResponseWriter object satisfies the interface because it has a w.Write() method.

## Templating

```
{{define "base"}}
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <title>{{template "title" .}} - Snippetbox</title>
  </head>
  <body>
    <header>
      <h1><a href="/">Snippetbox</a></h1>
    </header>
    <nav>
      <a href="/">Home</a>
    </nav>
    <section>{{template "body" .}}</section>
  </body>
</html>
{{end}}

```

- Here we’re using the {{define "base"}}...{{end}} action to define a distinct named template called base, which contains the content we want to appear on every page.
- The {{template "title" .}} and {{template "body" .}} actions denote that we want to invoke other named templates (called title and body) at a particular point in the HTML.

## FileServer

- We can serve static files with the help of inbuilt fileServer function

```
package main

import (
	"log"
	"net/http"
)
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project // directory root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))



	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}


```

We can then use the served files by adding relative path like this

```

<link rel='stylesheet' href='/static/css/main.css'>
<link rel='shortcut icon' href='/static/img/favicon.ico' type='image/x-icon'>
<!-- Also link to some fonts hosted by Google -->
<link rel='stylesheet' href='https://fonts.googleapis.com/css?family=Ubuntu+Mono:400,700'>

```

- Go fileServer sanitizes all request paths by running them through the `path.Clean()` function before searching for a file. This removes any . and .. elements from the URL path, which helps to stop directory traversal attacks. This feature is particularly useful if you’re using the fileserver in conjunction with a router that doesn’t automatically sanitize URL paths.

## Serve single file

- Sometimes you might want to serve a single file from within a handler. For this there’s the http.ServeFile() function, which you can use like so:

```
func downloadHandler(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "./ui/static/file.zip")
}
```

- But be aware: http.ServeFile() does not automatically sanitize the file path. If you’re constructing a file path from untrusted user input, to avoid directory traversal attacks you must sanitize the input with filepath.Clean() before using it.

## Handler

- A handler is an object which satisfies the http.Handler interface:

```
type Handler interface {
  ServeHTTP(ResponseWriter, *Request)
}
```

- a handler an object must have a ServeHTTP() method with the exact signature:

```
ServeHTTP(http.ResponseWriter, *http.Request)
```

## Command line flags

- we can send variable values from command line
- addr is a variable whose value is being passed here

=> `go run cmd/web/* -addr=":80"`

```
  // Define a new command-line flag with the name 'addr', a default value of ":4000"
  // and some short help text explaining what the flag controls. The value of the
  // flag will be stored in the addr variable at runtime.
  addr := flag.String("addr", ":4000", "HTTP network address")


	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any errors are
	// encountered during parsing the application will be terminated.
	flag.Parse()


  // The value returned from the flag.String() function is a pointer to the flag
  // value, not the value itself. So we need to dereference the pointer (i.e.
  // prefix it with the * symbol) before using it.
  log.Printf("Starting server on %s", *addr)


```

## Automated Help

- -help will list all of the possible flags that can be types in the run command

```
go run cmd/web/* -help
Usage of /tmp/go-build786121279/b001/exe/handlers:
     -addr string
        HTTP network address (default ":4000") exit status 2
```

## Leveled logging

- Go’s standard logger prefixes the message with the local date and time.
- We can create our own loggers
- We will use `log.New()` function to create two new custom loggers.
- If you want to include the full file path in your log output, instead of just the file name, you can use the `log.Llongfile` flag instead of `log.Lshortfile`

```
	// Use log.New() to create a logger for writing information messages. This takes
	// three parameters: the destination to write the logs to (os.Stdout), a string
	// prefix for message (INFO followed by a tab), and flags to indicate what
	// additional information to include (local date and time). Note that the flags
	// are joined using the bitwise OR operator |.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// Create a logger for writing error messages in the same way, but use stderr as
	// the destination and use the log.Lshortfile flag to include the relevant
	// file name and line number.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

    // Write messages using the two new loggers, instead of the standard logger.
    infoLog.Printf("Starting server on %s", *addr)
    err := http.ListenAndServe(*addr, mux) errorLog.Fatal(err)
```

Tip:

- By default, if Go’s HTTP server encounters an error it will log it using the standard logger. We can make it use our `errorLogger` we created.
- We need to initialize a new `http.Server` struct with the settings for our server and then we can use this new struct instead of using the `http.ListenAndServe()`

```
  // Initialize a new http.Server struct. We set the Addr and Handler fields so
  // that the server uses the same network address and routes as before, and set
  // the ErrorLog field so that the server now uses the custom errorLog logger in
  // the event of any problems.
  srv := &http.Server{
    Addr: *addr, ErrorLog: errorLog, Handler: mux,
}
  infoLog.Printf("Starting server on %s", *addr)
  // Call the ListenAndServe() method on our new http.Server struct.
  err := srv.ListenAndServe()
  errorLog.Fatal(err)

```

- As a rule of thumb, you should avoid using the Panic() and Fatal() variations outside of your main() function — it’s good practice to return errors instead, and only panic or exit directly from main().

## Writing logs to a file

```
  f, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0666)
  if err != nil {
    log.Fatal(err)
  }
  defer f.Close()
  infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)
```

## Dependencies injection

- Since we have now created new custom logger functions in main.go file we want them to also be used in handler functions.
- In order to link the dependencies in two files, we create a new `application` struct and add all of the handlers as methods/functions to this struct and we also add those two new logger functions as methods/functions to this new struct, now they are all linked to each other through this new struct.
- But this technique only works when all the handlers are situated in a single file.
- When the handlers are located in multiple files then we should follow something like this -> https://gist.github.com/alexedwards/5cd712192b4831058b21

## Status code Constants

- Instead of using 404 or 500 we can use the constants defined here => https://pkg.go.dev/net/http#pkg-constants
- Example =>

```
func (app *application) notFound(w http.ResponseWriter) {
  app.clientError(w, http.StatusNotFound)  // StatusNotFound is a constant for 404
}

```

## Database connection

```

import (
  "database/sql" // New import "flag"
  "log"
  "net/http"
  "os"
  _ "github.com/go-sql-driver/mysql" // New import
)


func main() {

  dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
  flag.Parse()


  db, err := openDB(*dsn)
  if err != nil {
    errorLog.Fatal(err)
  }

  // We also defer a call to db.Close(), so that the connection pool is closed
  // before the main() function exits.
  defer db.Close()
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil { return nil, err
	}
	if err = db.Ping(); err != nil {
	return nil, err }
	return db, nil
}

```

- Notice how the import path for our driver is prefixed with an underscore? This is because our main.go file doesn’t actually use anything in the mysql package. So if we try to import it normally the Go compiler will raise an error. However, we need the driver’s init() function to run so that it can register itself with the database/sql package. The trick to getting around this is to alias the package name to the blank identifier. This is standard practice for most of Go’s SQL drivers.

- The sql.Open() function doesn’t actually create any connections, all it does is initialize the pool for future use. Actual connections to the database are established lazily, as and when needed for the first time. So to verify that everything is set up correctly we need to use the db.Ping() method to create a connection and check for any errors.

- At this moment in time, the call to defer db.Close() is a bit superfluous. Our application is only ever terminated by a signal interrupt (i.e. Ctrl+c) or by errorLog.Fatal(). In both of those cases, the program exits immediately and deferred functions are never run. But including db.Close() is a good habit to get into and it could be beneficial later in the future if you add a graceful shutdown to your application.

## Designing a Database Model

- We’ll start by using the pkg/models/models.go file to define the top-level data types that
  our database model will use and return.
- pkg/models/mysql/snippets.go file, which will contain the code specifically for working with the snippets in our MySQL database, In this file we’re going to define a new SnippetModel type and implement some methods on it to access and manipulate the database. Like so:

```
// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
  DB *sql.DB
}

```

Using the SnippetModal

- To use this model in our handlers we need to establish a new SnippetModel struct in main()
- Then inject it as a dependency via the application struct

```
// Add a snippets field to the application struct. This will allow us to
// make the SnippetModel object available to our handlers.
type application struct {
  errorLog *log.Logger
  infoLog *log.Logger 
  snippets *mysql.SnippetModel
}

.
.
.


// Initialize a mysql.SnippetModel instance and add it to the application
// dependencies.
app := &application{
  errorLog: errorLog,
  infoLog: infoLog,
  snippets: &mysql.SnippetModel{DB: db},
}

```

Benefits of This Structure

- There’s a clean separation of concerns. Our database logic isn’t tied to our handlers which means that handler responsibilities are limited to HTTP stuff (i.e. validating requests and writing responses). This will make it easier to write tight, focused, unit tests in the future.
- By creating a custom SnippetModel type and implementing methods on it we’ve been able to make our model a single, neatly encapsulated object, which we can easily initialize and then pass to our handlers as a dependency. Again, this makes for easier to maintain, testable code.
- Because the model actions are defined as methods on an object — in our case SnippetModel — there’s the opportunity to create an interface and mock it for unit testing purposes.
- We have total control over which database is used at runtime, just by using the command-line flag.
- And finally, the directory structure scales nicely if your project has multiple back ends. For example, if some of your data is held in Redis you could put all the models for it in a pkg/models/redis package.

## Executing SQL Statements

- Notice how in this query we’re using the ? character to indicate placeholder parameters for the data that we want to insert in the database? Because the data we’ll be using will ultimately be untrusted user input from a form, it’s good practice to use placeholder parameters instead of interpolating data in the SQL query.
- In the code we constructed our SQL statement using placeholder parameters, where ? acted as a placeholder for the data we want to insert.
The reason for using placeholder parameters to construct our query (rather than string interpolation) is to help avoid SQL injection attacks from any untrusted user-provided input.


```
	stmt := `INSERT INTO snippets (title, content, created, expires)
			  VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
```

Go provides three different methods

- DB.Query() is used for SELECT queries which return multiple rows.
- DB.QueryRow() is used for SELECT queries which return a single row.
- DB.Exec() is used for statements which don’t return rows (like INSERT and DELETE).

```


stmt := `INSERT INTO snippets (title, content, created, expires)
VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

  // Use the Exec() method on the embedded connection pool to execute
  // statement. The first parameter is the SQL statement, followed by
  // title, content and expiry values for the placeholder parameters.
  // method returns a sql.Result object, which contains some basic
  // information about what happened when the statement was executed.

  result, err := m.DB.Exec(stmt, title, content, expires)
  if err != nil {
    return 0, err
  }


  // Use the LastInsertId() method on the result object to get the ID
  // newly inserted record in the snippets table.
  id, err := result.LastInsertId()
  if err != nil {
    return 0, err
  }

  // The ID returned has the type int64, so we convert it to an int type
  // before returning.
  return int(id), nil


```


The sql.Result interface returned by DB.Exec(). This provides two methods:

- LastInsertId() — which returns the integer (an int64) generated by the database in response to a command. 
- RowsAffected() — which returns the number of rows (as an int64) affected by the statement.
- `It’s important to note that not all drivers and databases support these two methods`
- 

Also, it is perfectly acceptable (and common) to ignore the sql.Result return value if you
don’t need it. Like so:

```
_, err := m.DB.Exec("INSERT INTO ...", ...)
```



Behind the scenes, the DB.Exec() method works in three steps:

- It creates a new prepared statement on the database using the provided SQL statement. The database parses and compiles the statement, then stores it ready for execution.
- In a second separate step, Exec() passes the parameter values to the database. The database then executes the prepared statement using these parameters. Because the parameters are transmitted later, after the statement has been compiled, the database treats them as pure data. They can’t change the intent of the statement. So long as the original statement is not derived from an untrusted data, injection cannot occur.
- It then closes (or deallocates) the prepared statement on the database.


## Single-record SQL Queries

- With the help of QueryRow we can query a single row

```


  stmt := `SELECT id, title, content, created, expires FROM snippets
  WHERE expires > UTC_TIMESTAMP() AND id = ?`
  // Use the QueryRow() method on the connection pool to execute our
  // SQL statement, passing in the untrusted id variable as the value for the 
  // placeholder parameter. This returns a pointer to a sql.Row object which 
  // holds the result from the database.
  row := m.DB.QueryRow(stmt, id)


	// Initialize a pointer to a new zeroed Snippet struct.
	s := &models.Snippet{}
	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct. Notice that the arguments
	// to row.Scan are *pointers* to the place you want to copy the data into,
	// and the number of arguments must be exactly the same as the number of
	// columns returned by your statement. If the query returns no rows, then
	// row.Scan() will return a sql.ErrNoRows error. We check for that and return 
	// our own models.ErrNoRecord error instead of a Snippet object.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
	return nil, models.ErrNoRecord } else if err != nil {
	return nil, err }
	// If everything went OK then return the Snippet object.
	return s, nil


```

- Aside: You might be wondering why we’re returning the models.ErrNoRecord error instead of sql.ErrNoRows directly. The reason is to help encapsulate the model completely, so that our application isn’t concerned with the underlying datastore or reliant on datastore- specific errors for its behavior.


# TIP
- Behind the scenes of rows.Scan() your driver will automatically convert the raw output from the SQL database to the required native Go types. So long as you’re sensible with the types that you’re mapping between SQL and Go, these conversions should generally Just Work. Usually:
  - CHAR, VARCHAR and TEXT map to string. 
  - BOOLEAN maps to bool.
  - INT maps to int; 
  - BIGINT maps to int64. 
  - DECIMAL and NUMERIC map to float.
  - TIME, DATE and TIMESTAMP map to time.Time.