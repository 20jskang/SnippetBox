package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox: as the response body
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	// w.Write([]byte("Hello from Snippetbox"))

	// Initialise a slice containing the paths to the two files. It's important
	// to note that the file contaiing our base template must be the *first*
	// file in the slice.
	files := []string{
		"./assets/html/base.tmpl.html",
		"./assets/html/pages/home.html",
	}
	
	// Use the template.ParseFiles() function to read the template file into a
	// template set. Notice that we use ... to pass the contents of the files
	// slice as variadic arguments.
	//  If there's an error, we log the detailed error message, use
	// the http.Error() function to send an Internal Server Error response to the
	// user, and then return from the handler so no subsequent code is executed.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Then we use the Execute() method on the template set to write the
	// template content as the response body. The last parameter to Execute()
	// represents any dynamic data that we want to pass in, which for now we'll
	// leave as nil.
	// 
	// err = ts.Execute(w, nil)
	// if err != nil {
	// 	log.Print(err.Error())
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// }
	// Use the ExecuteTemplate() method to write the content of the "base"
	// template as the response body
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	// Extract the value of the id wildcard from the request using r.PathValue()
	// and try to convert it to an integer using the strconv.Atoi function. If
	// it can't be converted into an integer, or the value is less than 1, we
	// return a 404 page not found response.
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Use the fmt.Fprintf() function to interpolate the id value with a
	// message, then write it as the HTTP response.
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Use the w.WriteHeader() method to send a 201 status code.
	w.WriteHeader(http.StatusCreated)

	// Then use w.Write() method to write the response body as normal.
	w.Write([]byte("Save a new snippet..."))
}
