package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
)

// TemplateData struct for passing data to templates
type TemplateData struct {
	Title       string
	PageName    string
	CurrentTime string
}

// Global variable to store parsed templates
var templates *template.Template

func init() {
	dbpath := "/home/whitepi/go/slideshowgo/imagesDB"
	imagedir := "/home/whitepi/Pictures/"
	Walk_Img_Dir(dbpath, imagedir)
	// Parse all templates in the "templates" directory.
	// template.Must panics if there's an error, which is good for quick startup
	// errors for templates. In a larger app, you might handle errors more gracefully.
	templates = template.Must(template.ParseGlob("templates/*.html"))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		Title:       "Home - My Go App",
		PageName:    "Home",
		CurrentTime: time.Now().Format("Mon Jan 2 15:04:05 MST 2006"),
	}
	err := templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		Title:    "About Us - My Go App",
		PageName: "About",
	}
	err := templates.ExecuteTemplate(w, "about.html", data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
	}
}

// serveStaticFiles sets up a file server for static assets (like CSS, JS, images).
func serveStaticFiles(router *mux.Router) {
    // Create a file server for the "static" directory.
    // Ensure you create a "static" directory in your project if you use this.
    // Example: my-web-app/static/css/style.css
    staticFileServer := http.FileServer(http.Dir("static"))

    // Use PathPrefix to match any request starting with /static/
    // StripPrefix removes the /static/ part from the URL path before
    // FileServer looks for the file on disk.
    router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticFileServer))
}


func main() {
	router := mux.NewRouter()

	// Register handlers for HTML templates
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/about", aboutHandler).Methods("GET")

	// Serve static files (optional, but good practice for real apps)
	// If you have CSS, JS, images, etc., put them in a 'static' folder.
	// You might create a `static` directory like `my-web-app/static/css/style.css`
	serveStaticFiles(router)


	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
