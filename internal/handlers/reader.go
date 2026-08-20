package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/CephandriusMaxtori/Folio/internal/models"
	"github.com/CephandriusMaxtori/Folio/internal/readers"
)

func (h *Handler) GetChapterPages(w http.ResponseWriter, r *http.Request) {
	id := uintParam(r, "id")
	ch, err := h.svc.GetChapter(id)
	if err != nil {
		writeError(w, 404, "chapter not found")
		return
	}

	pages := ch.PageCount
	if pages == 0 {
		pages = 1
	}
	writeJSON(w, 200, map[string]int{
		"pages":     pages,
		"chapterId": int(ch.ID),
	})
}

func (h *Handler) GetPage(w http.ResponseWriter, r *http.Request) {
	chID := uintParam(r, "id")
	num := queryInt(r, "num", 0)

	ch, err := h.svc.GetChapter(chID)
	if err != nil {
		writeError(w, 404, "chapter not found")
		return
	}

	switch ch.FileType {
	case "cbz", "zip", "cbr", "rar":
		h.serveArchivePage(w, r, ch, num)
	case "epub":
		h.serveEpubPage(w, r, ch, num)
	default:
		writeError(w, 400, "unsupported file type")
	}
}

func (h *Handler) serveArchivePage(w http.ResponseWriter, r *http.Request, ch *models.Chapter, num int) {
	archive, err := readers.Open(ch.FilePath)
	if err != nil {
		log.Printf("Failed to open archive %s: %v", ch.FilePath, err)
		writeError(w, 500, "failed to open archive")
		return
	}
	defer archive.Close()

	data, mime, err := archive.GetPage(num)
	if err != nil {
		writeError(w, 404, "page not found")
		return
	}

	w.Header().Set("Content-Type", mime)
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write(data)
}

func (h *Handler) serveEpubPage(w http.ResponseWriter, r *http.Request, ch *models.Chapter, num int) {
	epub, err := readers.OpenEpub(ch.FilePath)
	if err != nil {
		log.Printf("Failed to open epub %s: %v", ch.FilePath, err)
		writeError(w, 500, "failed to open epub")
		return
	}
	defer epub.Close()

	text, err := epub.ChapterText(num)
	if err != nil {
		writeError(w, 404, "chapter not found")
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(text))
}

func (h *Handler) SaveProgress(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	var req struct {
		ChapterID uint    `json:"chapterId"`
		Page      int     `json:"page"`
		ScrollPct float64 `json:"scrollPct"`
	}
	if err := readJSON(r, &req); err != nil {
		writeError(w, 400, "invalid request")
		return
	}

	p, err := h.svc.GetProgress(user.ID, req.ChapterID)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}

	p.Page = req.Page
	p.ScrollPct = req.ScrollPct
	p.LastRead = time.Now()
	p.ReadCount++

	if err := h.svc.SaveProgress(p); err != nil {
		writeError(w, 500, err.Error())
		return
	}

	writeJSON(w, 200, map[string]string{"status": "ok"})
}

func (h *Handler) OnDeck(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	progress, err := h.svc.GetOnDeck(user.ID)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}

	type OnDeckItem struct {
		ChapterID   uint      `json:"chapterId"`
		ChapterName string    `json:"chapterName"`
		SeriesName  string    `json:"seriesName"`
		SeriesID    uint      `json:"seriesId"`
		Page        int       `json:"page"`
		TotalPages  int       `json:"totalPages"`
		LastRead    time.Time `json:"lastRead"`
	}

	var items []OnDeckItem
	for _, p := range progress {
		ch, err := h.svc.GetChapter(p.ChapterID)
		if err != nil {
			continue
		}
		ser, err := h.svc.GetSeries(ch.SeriesID)
		if err != nil {
			continue
		}
		items = append(items, OnDeckItem{
			ChapterID:   ch.ID,
			ChapterName: ch.Title,
			SeriesName:  ser.Name,
			SeriesID:    ser.ID,
			Page:        p.Page,
			TotalPages:  ch.PageCount,
			LastRead:    p.LastRead,
		})
	}

	if items == nil {
		items = []OnDeckItem{}
	}
	writeJSON(w, 200, items)
}

func (h *Handler) SeriesCover(w http.ResponseWriter, r *http.Request) {
	id := uintParam(r, "id")
	ser, err := h.svc.GetSeries(id)
	if err != nil {
		writeError(w, 404, "series not found")
		return
	}

	if ser.CoverPath == "" {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, ser.CoverPath)
}

func (h *Handler) VolumeCover(w http.ResponseWriter, r *http.Request) {
	id := uintParam(r, "id")
	vol, err := h.svc.GetVolume(id)
	if err != nil {
		writeError(w, 404, "volume not found")
		return
	}

	if vol.CoverPath == "" {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, vol.CoverPath)
}

func (h *Handler) ChapterCover(w http.ResponseWriter, r *http.Request) {
	id := uintParam(r, "id")
	ch, err := h.svc.GetChapter(id)
	if err != nil {
		writeError(w, 404, "chapter not found")
		return
	}

	switch ch.FileType {
	case "cbz", "zip", "cbr", "rar":
		archive, err := readers.Open(ch.FilePath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer archive.Close()

		data, mime, err := archive.DetectCover()
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", mime)
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Write(data)

	case "epub":
		epub, err := readers.OpenEpub(ch.FilePath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer epub.Close()

		data, mime, err := epub.CoverImage()
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", mime)
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Write(data)

	default:
		http.NotFound(w, r)
	}
}

func (h *Handler) GetSeriesVolumes(w http.ResponseWriter, r *http.Request) {
	id := uintParam(r, "id")
	ser, err := h.svc.GetSeries(id)
	if err != nil {
		writeError(w, 404, "series not found")
		return
	}
	writeJSON(w, 200, ser.Volumes)
}

func (h *Handler) ScanLibrary(w http.ResponseWriter, r *http.Request) {
	id := uintParam(r, "id")
	go func() {
		result, err := h.svc.ScanLibrary(id)
		if err != nil {
			log.Printf("Scan failed for library %d: %v", id, err)
			return
		}
		log.Printf("Scan complete for library %d: %d series, %d chapters, %d errors",
			id, result.SeriesFound, result.ChaptersFound, len(result.Errors))
	}()
	writeJSON(w, 202, map[string]string{"status": "scan started"})
}
