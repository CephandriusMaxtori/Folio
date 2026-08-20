package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/CephandriusMaxtori/Folio/internal/config"
	"github.com/CephandriusMaxtori/Folio/internal/models"
	"github.com/CephandriusMaxtori/Folio/internal/services"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc  *services.Service
	cfg  *config.Config
}

func New(svc *services.Service, cfg *config.Config) *Handler {
	return &Handler{svc: svc, cfg: cfg}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api", func(r chi.Router) {
		r.Post("/auth/register", h.Register)
		r.Post("/auth/login", h.Login)

		r.Group(func(r chi.Router) {
			r.Use(h.AuthMiddleware)
			r.Get("/auth/me", h.Me)

			r.Route("/libraries", func(r chi.Router) {
				r.Get("/", h.ListLibraries)
				r.Post("/", h.CreateLibrary)
				r.Put("/{id}", h.UpdateLibrary)
				r.Delete("/{id}", h.DeleteLibrary)
				r.Post("/{id}/scan", h.ScanLibrary)
			})

			r.Get("/series", h.ListSeries)
			r.Get("/series/{id}", h.GetSeries)

			r.Route("/reader", func(r chi.Router) {
				r.Get("/chapter/{id}/pages", h.GetChapterPages)
				r.Get("/chapter/{id}/page/{num}", h.GetPage)
				r.Post("/progress", h.SaveProgress)
				r.Get("/on-deck", h.OnDeck)
			})

			r.Get("/search", h.Search)
			r.Get("/settings", h.GetSettings)
			r.Put("/settings", h.UpdateSettings)
			r.Get("/admin/stats", h.Stats)
		})
	})

	r.Route("/opds", func(r chi.Router) {
		r.Get("/{apiKey}", h.OPDSCatalog)
		r.Get("/{apiKey}/libraries", h.OPDSLibraries)
		r.Get("/{apiKey}/library/{libraryId}", h.OPDSLibrary)
		r.Get("/{apiKey}/series/{seriesId}", h.OPDSSeries)
		r.Get("/{apiKey}/series/{seriesId}/volume/{volumeId}", h.OPDSVolume)
		r.Get("/{apiKey}/series/{seriesId}/volume/{volumeId}/chapter/{chapterId}", h.OPDSChapter)
		r.Get("/{apiKey}/series/{seriesId}/volume/{volumeId}/chapter/{chapterId}/download/{filename}", h.OPDSDownload)
		r.Get("/search/{apiKey}", h.OPDSSearchDescriptor)
	})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func uintParam(r *http.Request, name string) uint {
	s := chi.URLParam(r, name)
	v, _ := strconv.ParseUint(s, 10, 32)
	return uint(v)
}

func queryInt(r *http.Request, name string, def int) int {
	s := r.URL.Query().Get(name)
	if s == "" {
		return def
	}
	v, _ := strconv.Atoi(s)
	return v
}

func queryStr(r *http.Request, name string) string {
	return r.URL.Query().Get(name)
}

func (h *Handler) ListLibraries(w http.ResponseWriter, r *http.Request) {
	libs, err := h.svc.ListLibraries()
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, libs)
}

func (h *Handler) CreateLibrary(w http.ResponseWriter, r *http.Request) {
	var lib models.Library
	if err := json.NewDecoder(r.Body).Decode(&lib); err != nil {
		writeError(w, 400, "invalid request")
		return
	}
	if err := h.svc.CreateLibrary(&lib); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, lib)
}

func (h *Handler) UpdateLibrary(w http.ResponseWriter, r *http.Request) {
	id := uintParam(r, "id")
	existing, err := h.svc.GetLibrary(id)
	if err != nil {
		writeError(w, 404, "not found")
		return
	}
	if err := json.NewDecoder(r.Body).Decode(existing); err != nil {
		writeError(w, 400, "invalid request")
		return
	}
	if err := h.svc.UpdateLibrary(existing); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, existing)
}

func (h *Handler) DeleteLibrary(w http.ResponseWriter, r *http.Request) {
	id := uintParam(r, "id")
	if err := h.svc.DeleteLibrary(id); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, map[string]string{"status": "deleted"})
}

func (h *Handler) ScanLibrary(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 202, map[string]string{"status": "scan started"})
}

func (h *Handler) ListSeries(w http.ResponseWriter, r *http.Request) {
	libID := uint(queryInt(r, "library_id", 0))
	sort := queryStr(r, "sort")
	series, err := h.svc.ListSeries(libID, sort)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, series)
}

func (h *Handler) GetSeries(w http.ResponseWriter, r *http.Request) {
	id := uintParam(r, "id")
	ser, err := h.svc.GetSeries(id)
	if err != nil {
		writeError(w, 404, "not found")
		return
	}
	writeJSON(w, 200, ser)
}

func (h *Handler) GetChapterPages(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]int{"pages": 0})
}

func (h *Handler) GetPage(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", 501)
}

func (h *Handler) SaveProgress(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]string{"status": "ok"})
}

func (h *Handler) OnDeck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, []interface{}{})
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	q := queryStr(r, "q")
	libID := uint(queryInt(r, "library_id", 0))
	series, err := h.svc.Search(q, libID)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, series)
}

func (h *Handler) GetSettings(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]string{"theme": "dark"})
}

func (h *Handler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]string{"status": "ok"})
}

func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.svc.GetStats()
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, stats)
}
