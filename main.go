package main

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	appdb "glossary/internal/db"
	"glossary/internal/terms"
	"html/template"
	"log"
	"net/http"
	"os"
)

//go:embed scripts/database.sql
var schemaFS embed.FS

//go:embed web/templates/*.tmpl.html
var templatesFS embed.FS

func main() {
	_ = godotenv.Load()
	db, err := appdb.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	//	ensureSchema(db)

	tmpl := template.Must(template.ParseFS(templatesFS, "web/templates/*.tmpl.html"))
	repo := terms.NewRepository(db)
	h := terms.NewHandler(repo, tmpl)

	r := chi.NewRouter()
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	r.Mount("/terms", h.Routes())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s", port)
	http.ListenAndServe(":"+port, r)
}

func ensureSchema(db *sql.DB) {
	b, _ := schemaFS.ReadFile("db/schema.sql")
	db.Exec(string(b))
}
