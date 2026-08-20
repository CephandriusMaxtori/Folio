package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/CephandriusMaxtori/Folio/internal/models"
)

func (h *Handler) ListCollections(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	cols, err := h.svc.ListCollections(user.ID)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, cols)
}

func (h *Handler) CreateCollection(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	var col models.Collection
	if err := readJSON(r, &col); err != nil {
		writeError(w, 400, "invalid request")
		return
	}
	col.UserID = user.ID
	if err := h.svc.CreateCollection(&col); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, col)
}

func (h *Handler) AddToCollection(w http.ResponseWriter, r *http.Request) {
	id := uintParam(r, "id")
	var req struct {
		SeriesID uint `json:"seriesId"`
		Position int  `json:"position"`
	}
	if err := readJSON(r, &req); err != nil {
		writeError(w, 400, "invalid request")
		return
	}
	if err := h.svc.AddToCollection(id, req.SeriesID, req.Position); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, map[string]string{"status": "added"})
}

func (h *Handler) ListReadingLists(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	lists, err := h.svc.ListReadingLists(user.ID)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, lists)
}

func (h *Handler) CreateReadingList(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	var rl models.ReadingList
	if err := readJSON(r, &rl); err != nil {
		writeError(w, 400, "invalid request")
		return
	}
	rl.UserID = user.ID
	if err := h.svc.CreateReadingList(&rl); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, rl)
}

func (h *Handler) AddToReadingList(w http.ResponseWriter, r *http.Request) {
	id := uintParam(r, "id")
	var req struct {
		ChapterID uint `json:"chapterId"`
		Position  int  `json:"position"`
	}
	if err := readJSON(r, &req); err != nil {
		writeError(w, 400, "invalid request")
		return
	}
	if err := h.svc.AddToReadingList(id, req.ChapterID, req.Position); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, map[string]string{"status": "added"})
}

func (h *Handler) ListAnnotations(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	chapterID := uint(queryInt(r, "chapter_id", 0))
	anns, err := h.svc.ListAnnotations(user.ID, chapterID)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, anns)
}

func (h *Handler) CreateAnnotation(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	var ann models.Annotation
	if err := readJSON(r, &ann); err != nil {
		writeError(w, 400, "invalid request")
		return
	}
	ann.UserID = user.ID
	if err := h.svc.CreateAnnotation(&ann); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, ann)
}

func (h *Handler) UpdateAnnotation(w http.ResponseWriter, r *http.Request) {
	id := uintParam(r, "id")
	var ann models.Annotation
	if err := readJSON(r, &ann); err != nil {
		writeError(w, 400, "invalid request")
		return
	}
	ann.ID = id
	if err := h.svc.UpdateAnnotation(&ann); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, ann)
}

func (h *Handler) DeleteAnnotation(w http.ResponseWriter, r *http.Request) {
	id := uintParam(r, "id")
	if err := h.svc.DeleteAnnotation(id); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, map[string]string{"status": "deleted"})
}

func (h *Handler) GetSettings(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	st, err := h.svc.GetSettings(user.ID)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, st)
}

func (h *Handler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	st, err := h.svc.GetSettings(user.ID)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	if err := readJSON(r, st); err != nil {
		writeError(w, 400, "invalid request")
		return
	}
	if err := h.svc.UpdateSettings(st); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, st)
}

func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.svc.GetStats()
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, stats)
}

func (h *Handler) ListAPIKeys(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	keys, err := h.svc.ListAPIKeys(user.ID)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, keys)
}

func (h *Handler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	var req struct {
		Name string `json:"name"`
	}
	if err := readJSON(r, &req); err != nil {
		writeError(w, 400, "invalid request")
		return
	}
	if req.Name == "" {
		req.Name = "API Key"
	}

	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		writeError(w, 500, "failed to generate key")
		return
	}
	keyStr := hex.EncodeToString(b)

	key := &models.APIKey{
		UserID: user.ID,
		Key:    keyStr,
		Name:   req.Name,
	}
	if err := h.svc.CreateAPIKey(key); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, key)
}

func (h *Handler) DeleteAPIKey(w http.ResponseWriter, r *http.Request) {
	id := uintParam(r, "id")
	if err := h.svc.DeleteAPIKey(id); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, map[string]string{"status": "deleted"})
}
