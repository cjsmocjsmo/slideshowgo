package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type ImageData struct {
	Name        string
	Path        string
	Http        string
	Idx         int
	Orientation string
}

// Global variable to store parsed templates
var templates *template.Template
var dbpath = "/home/pimedia/imagesDB"
var imagedir = "/home/pimedia/Pictures/"

// Global variables for slideshow control
var currentImageIdx int = 1
var imageMutex sync.RWMutex

func init() {
	// Parse all templates in the "templates" directory.
	// template.Must panics if there's an error, which is good for quick startup
	// errors for templates. In a larger app, you might handle errors more gracefully.
	templates = template.Must(template.ParseGlob("templates/*.html"))
}

func db_count() int {
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		log.Printf("Error opening count database: %v", err)
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

func get_db_image(idx int) (ImageData, error) {
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return ImageData{}, err
	}
	defer db.Close()

	// var img ImageData
	// query := "SELECT name, path, http, idx, orientation FROM images WHERE idx = ?"
	// err = db.QueryRow(query, idx).Scan(&img.Name, &img.Path, &img.Http, &img.Idx, &img.Orientation)
	// if err != nil {
	// 	log.Printf("Error querying get_db_image: %v", err)
	// 	return ImageData{}, err
	// }
	// return img, nil
	var img ImageData
	query := "SELECT name, http, idx, orientation FROM images WHERE idx = ?"
	err = db.QueryRow(query, idx).Scan(&img.Name, &img.Http, &img.Idx, &img.Orientation)
	if err != nil {
		log.Printf("Error querying get_db_image: %v", err)
		return ImageData{}, err
	}
	return img, nil
}

var dbcount = db_count()

// startSlideshow starts the automatic slideshow timer
func startSlideshow() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			imageMutex.Lock()
			currentImageIdx++
			if currentImageIdx > dbcount {
				currentImageIdx = 1
			}
			imageMutex.Unlock()
			log.Printf("Slideshow advanced to image %d", currentImageIdx)
		}
	}()
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	imageMutex.RLock()
	idx := currentImageIdx
	imageMutex.RUnlock()

	fmt.Println("db_count:", dbcount)

	data, err1 := get_db_image(idx)
	if err1 != nil {
		log.Printf("Error getting image from database: %v", err1)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
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
	// Start the slideshow timer
	startSlideshow()

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
