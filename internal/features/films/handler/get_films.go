package film_handler

import "net/http"

func (h *filmHandler) GetFilms(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get Films"))
}
