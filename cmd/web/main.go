package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"davappler/snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include fields for the two custom loggers, but // we'll add more to it as the build progresses.
type application struct {
	errorLog *log.Logger
	infoLog *log.Logger 
	snippets *mysql.SnippetModel
}

func main() {



	// Define a new command-line flag with the name 'addr', a default value of ":4000" 
	// and some short help text explaining what the flag controls. The value of the
	// flag will be stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Define a new command-line flag for the MySQL DSN string.
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any errors are 
	// encountered during parsing the application will be terminated.
	flag.Parse()



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


	db, err := openDB(*dsn)
	if err != nil { 
		errorLog.Fatal(err)
	}

	defer db.Close()

	app := &application{ 
		errorLog: errorLog, 
		infoLog: infoLog,
		snippets: &mysql.SnippetModel{DB: db},
	}

	srv := &http.Server{ 
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(), // Call the new app.routes() method 
	}


	// Write messages using the two new loggers, instead of the standard logger.
	infoLog.Printf("Starting server on %s", *addr) 
	errr := srv.ListenAndServe()
	errorLog.Fatal(errr)

}



func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil { return nil, err
	}
	if err = db.Ping(); err != nil {
	return nil, err }
	return db, nil 
}


