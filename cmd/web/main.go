package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

// Define an application struct to hold the application-wide dependencies for the 
// web application. For now we'll only include the structured logger, but we'll 
// add more to this as development progresses.
type application struct {
	logger *slog.Logger
}

func main() {
	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the
	// flag will be stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Importantly, we use flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any errors
	// are encountered during parsing the application will be terminated.
	flag.Parse()

	// Use the slog.New() function to initialise a new structured logger, which
	// writes to the standard out stream and uses the default settings.
	// logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)) for JSON output format
	// logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug,}))
	// logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true,}))
	// 
	// go run ./cmd/web >> /tmp/web.log
	// Redirects the standard out stream to an on-disk file when starting the application
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Initialise a new instance of our application struct, containing the dependencies
	app := &application{
		logger: logger,
	}
	// Use the http.NewServeMux() function to initialise a new servemux
	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./assets/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir("./assets/static/"))

	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Register the other application routes as normal.
	// Swap the route declaration to use the application struct's methods as the 
	// handler function
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	// Print a log message to say that the server is starting
	//
	// The value returned from flag.String() function is a pointer to the flag
	// value, not the value itself. So in this code, that means the addr variable
	// is actually a pointer, and we need to dereference it before using it.
	// Note that we're using the log.Printf() function to interpolate the address
	// with the log message.
	// 
	// Use the Info() method to log the starting server message at Info severity
	// (along with the listen address as an attribute).
	logger.Info("Starting server", "addr", *addr)

	// Use the http.ListenAndServe() function to start a new web server. We pass in
	// two parameters: the TCP network address to listen on (in this case ":4000")
	// and the servemux we just created. If http.ListenAndServe() returns an error
	// we use the log.Fatal() function to log the error message and terminate the
	// program. Note that any error returned by http.ListenAndServe() is always
	// non-nil.
	err := http.ListenAndServe(*addr, mux)

	// And we also use the Error() method to log any error message returned by
	// http.ListenAndServe() at Error severity (with no additional attributes),
	// and then call os.Exit(1) to terminate the application with exit code 1.
	logger.Error(err.Error())
	os.Exit(1)
}
