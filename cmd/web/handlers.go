package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)



func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}


	// Initialize a slice containing the paths to the two files. Note that the 
	// home.page.tmpl file must be the *first* file in the slice.
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl", 
		"./ui/html/footer.partial.tmpl",
	}


	// Notice that we can pass the slice of file paths // as a variadic parameter?
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
		return
	}



	// We then use the Execute() method on the template set to write the template 
	// content as the response body. The last parameter to Execute() represents any
	// dynamic data that we want to pass in, which for now we'll leave as nil.

	err = ts.Execute(w, nil)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err) 
	}
}


func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w) 
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}




func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed) // Use the clientError() helper. return
		}
	w.Write([]byte("Create a new snippet...")) 
}