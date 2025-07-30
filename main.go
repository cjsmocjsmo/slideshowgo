package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
	// "time"
)

// TemplateData struct for passing data to templates
// type TemplateData struct {
// 	Title       string
// 	PageName    string
// 	CurrentTime string
// }

type ImageData struct {
	Name        string
	Path        string
	Http        string
	Idx         int
	Orientation string
}

// Global variable to store parsed templates
var templates *template.Template
var dbpath = "/home/pimedia/Pictures/imagesDB"
var imagedir = "/home/pimedia/Pictures/"

func init() {
	// Parse all templates in the "templates" directory.
	// template.Must panics if there's an error, which is good for quick startup
	// errors for templates. In a larger app, you might handle errors more gracefully.
	templates = template.Must(template.ParseGlob("templates/*.html"))
}

func db_count() int {
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return 0
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM images").Scan(&count)
	if err != nil {
		log.Printf("Error querying count: %v", err)
		return 0
	}
	return count
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := ImageData{
		Name:        "SlideShow",
		Path:        "/static/Pics1/images_part_001/83bcf227931a9595.jpg",
		Http:        "http://10.0.4.41:8080/static/Pics1/images_part_001/83bcf227931a9595.jpg",
		Idx:         1,
		Orientation: "landscape",
	}

	err := templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
	}
}

// serveStaticFiles sets up a file server for static assets (like CSS, JS, images).
func serveStaticFiles(router *mux.Router) {
	// Serve static files from /home/pimedia/Pictures/
	staticFileServer := http.FileServer(http.Dir("/home/pimedia/Pictures/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticFileServer))
}

func main() {
	router := mux.NewRouter()

	// Register handlers for HTML templates
	router.HandleFunc("/", homeHandler).Methods("GET")

	// Serve static files (optional, but good practice for real apps)
	// If you have CSS, JS, images, etc., put them in a 'static' folder.
	// You might create a `static` directory like `my-web-app/static/css/style.css`
	serveStaticFiles(router)

	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
