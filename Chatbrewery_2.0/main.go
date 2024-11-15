package main

import (
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User-Datenbankmodell
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Email    string
	Password string
}

var db *gorm.DB

// Datenbank initialisieren
func initDatabase() {
	dsn := "root:@tcp(127.0.0.1:3306)/chatbrewery_2.0?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Tabellen migrieren
	db.AutoMigrate(&User{})
	log.Println("Database connected and User table created.")
}

// Startpunkt
func main() {
	// Initialisiere die Datenbank
	initDatabase()

	// Router initialisieren
	r := http.NewServeMux()

	// Routen registrieren
	r.HandleFunc("/", renderIndex)
	r.HandleFunc("/signup", handleSignup)
	r.HandleFunc("/login", handleLogin)

	// Statische Dateien (Bilder, CSS, JS) bereitstellen
	r.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Server starten
	log.Println("Server läuft auf http://localhost:8080")
	http.ListenAndServe(":8080", r)
}

// Render Index-Seite mit Light/Dark-Mode und Login/Signup-Option
func renderIndex(w http.ResponseWriter, r *http.Request) {
	// URL-Parameter (z. B. Theme oder Formular-Typ)
	theme := r.URL.Query().Get("theme")
	if theme == "" {
		theme = "light" // Standard: Light-Mode
	}

	form := r.URL.Query().Get("form")
	if form == "" {
		form = "signup" // Standard: Signup
	}

	// Daten an Template übergeben
	data := struct {
		Theme    string
		FormType string
	}{
		Theme:    theme,
		FormType: form,
	}

	// Template rendern
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

// Handle Signup-Formular
func handleSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Formulardaten auslesen
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Passwort hashen
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Benutzer speichern
	user := User{Username: username, Email: email, Password: string(hashedPassword)}
	if result := db.Create(&user); result.Error != nil {
		http.Error(w, "Could not save user: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Erfolg: Weiterleitung zur Login-Seite
	http.Redirect(w, r, "/?theme=light&form=login", http.StatusSeeOther)
}

// Handle Login-Formular (noch nicht implementiert)
func handleLogin(w http.ResponseWriter, r *http.Request) {
	// Noch nicht implementiert
	http.Error(w, "Login functionality not implemented yet", http.StatusNotImplemented)
}
