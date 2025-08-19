package terms

import (
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
	"strconv"
)

type Handler struct {
	repo *Repository
	tmpl *template.Template
}

func NewHandler(repo *Repository, tmpl *template.Template) *Handler {
	return &Handler{repo: repo, tmpl: tmpl}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.List)
	r.Get("/new", h.NewForm)
	r.Post("/create", h.Create)
	r.Get("/{id}", h.View)
	r.Get("/{id}/edit", h.EditForm)
	r.Post("/{id}/update", h.Update)
	r.Post("/{id}/delete", h.Delete)
	return r
}
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	terms, err := h.repo.All(ctx)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	h.tmpl.ExecuteTemplate(w, "terms_list.tmpl.html", terms)
}

func (h *Handler) View(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	t, err := h.repo.ByID(r.Context(), id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	h.tmpl.ExecuteTemplate(w, "terms_view.tmpl.html", t)
}

func (h *Handler) NewForm(w http.ResponseWriter, r *http.Request) {
	h.tmpl.ExecuteTemplate(w, "terms_form.tmpl.html", nil)
}
func (h *Handler) EditForm(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	t, err := h.repo.ByID(r.Context(), id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	h.tmpl.ExecuteTemplate(w, "terms_form.tmpl.html", t)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	t := &Term{
		TermEn:       r.FormValue("term_en"),
		DefinitionEn: r.FormValue("definition_en"),
		Context:      r.FormValue("context"),
	}
	h.repo.Create(r.Context(), t)
	http.Redirect(w, r, "/terms/", http.StatusSeeOther)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	t := &Term{
		ID:           id,
		TermEn:       r.FormValue("term_en"),
		DefinitionEn: r.FormValue("definition_en"),
		Context:      r.FormValue("context"),
	}
	h.repo.Update(r.Context(), t)
	http.Redirect(w, r, "/terms/", http.StatusSeeOther)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	h.repo.Delete(r.Context(), id)
	http.Redirect(w, r, "/terms/", http.StatusSeeOther)
}
