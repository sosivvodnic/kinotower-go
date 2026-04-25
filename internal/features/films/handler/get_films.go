package film_handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/sosivvodnic/kinotower-go/internal/core/httpjson"
	"github.com/sosivvodnic/kinotower-go/internal/features/films/domain"
)

func (h *filmHandler) GetFilms(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	page := parseIntDefault(q.Get("page"), 1)
	size := parseIntDefault(q.Get("size"), 10)
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	sortBy := q.Get("sortBy")
	if strings.TrimSpace(sortBy) == "" {
		sortBy = "name"
	}
	sortDir := q.Get("sortDir")
	if strings.TrimSpace(sortDir) == "" {
		sortDir = "asc"
	}

	category := parseIntDefault(q.Get("category"), 0)
	country := parseIntDefault(q.Get("country"), 0)
	search := strings.TrimSpace(q.Get("search"))
	var searchPtr *string
	if search != "" {
		searchPtr = &search
	}

	filter := domain.FilmFilter{
		Page:     page,
		Size:     size,
		SortBy:   sortBy,
		SortDir:  sortDir,
		Category: category,
		Country:  country,
		Search:   searchPtr,
	}

	films, total, err := h.filmService.GetFilms(r.Context(), filter)
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, domain.FilmListResponse{
		Page:  page,
		Size:  size,
		Total: total,
		Films: films,
	})
}

func (h *filmHandler) GetFilmByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		httpjson.WriteError(w, http.StatusNotFound, "Film not found")
		return
	}

	film, err := h.filmService.GetFilmByID(r.Context(), id)
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if film == nil {
		httpjson.WriteError(w, http.StatusNotFound, "Film not found")
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, film)
}

func (h *filmHandler) GetFilmReviews(w http.ResponseWriter, r *http.Request) {
	filmID, err := strconv.Atoi(chi.URLParam(r, "film-id"))
	if err != nil || filmID <= 0 {
		httpjson.WriteError(w, http.StatusNotFound, "Film not found")
		return
	}

	exists, err := h.filmService.FilmExists(r.Context(), filmID)
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if !exists {
		httpjson.WriteError(w, http.StatusNotFound, "Film not found")
		return
	}

	reviews, err := h.filmService.GetApprovedReviewsByFilmID(r.Context(), filmID)
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, map[string]any{
		"reviews": reviews,
	})
}

func parseIntDefault(v string, def int) int {
	v = strings.TrimSpace(v)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}
