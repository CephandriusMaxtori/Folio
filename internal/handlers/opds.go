package handlers

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type opdsFeed struct {
	XMLName xml.Name   `xml:"feed"`
	Xmlns   string     `xml:"xmlns,attr"`
	Title   string     `xml:"title"`
	ID      string     `xml:"id"`
	Updated string     `xml:"updated"`
	Links   []opdsLink `xml:"link"`
	Entries []opdsEntry `xml:"entry"`
}

type opdsEntry struct {
	Title   string     `xml:"title"`
	ID      string     `xml:"id"`
	Updated string     `xml:"updated"`
	Content string     `xml:"content,omitempty"`
	Links   []opdsLink `xml:"link"`
}

type opdsLink struct {
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
	Type string `xml:"type,attr,omitempty"`
}

func opdsNow() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func opdsBase(apiKey, prefix string) string {
	if prefix == "" {
		return fmt.Sprintf("/opds/%s", apiKey)
	}
	return fmt.Sprintf("/opds/%s/%s", apiKey, prefix)
}

func writeOPDS(w http.ResponseWriter, feed opdsFeed) {
	w.Header().Set("Content-Type", "application/atom+xml;profile=opds-catalog")
	w.Write([]byte(xml.Header))
	xml.NewEncoder(w).Encode(feed)
}

func (h *Handler) OPDSAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := chi.URLParam(r, "apiKey")
		if apiKey == "" {
			writeError(w, 401, "missing API key")
			return
		}

		k, err := h.svc.GetAPIKey(apiKey)
		if err != nil {
			writeError(w, 401, "invalid API key")
			return
		}
		if k.ExpiresAt != nil && k.ExpiresAt.Before(time.Now()) {
			writeError(w, 401, "API key expired")
			return
		}

		user, err := h.svc.GetUserByID(k.UserID)
		if err != nil {
			writeError(w, 401, "user not found")
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) OPDSCatalog(w http.ResponseWriter, r *http.Request) {
	apiKey := chi.URLParam(r, "apiKey")

	feed := opdsFeed{
		Xmlns:   "http://www.w3.org/2005/Atom",
		Title:   "Folio",
		ID:      "urn:folio:catalog",
		Updated: opdsNow(),
		Links: []opdsLink{
			{Rel: "self", Href: opdsBase(apiKey, ""), Type: "application/atom+xml;profile=opds-catalog"},
			{Rel: "start", Href: opdsBase(apiKey, ""), Type: "application/atom+xml;profile=opds-catalog"},
			{Rel: "search", Href: fmt.Sprintf("/opds/search/%s?searchTerms={searchTerms}", apiKey), Type: "application/opensearchdescription+xml"},
		},
		Entries: []opdsEntry{
			{
				Title:   "Libraries",
				ID:      "urn:folio:libraries",
				Updated: opdsNow(),
				Links: []opdsLink{
					{Rel: "subsection", Href: opdsBase(apiKey, "libraries"), Type: "application/atom+xml;profile=opds-catalog"},
				},
			},
			{
				Title:   "Collections",
				ID:      "urn:folio:collections",
				Updated: opdsNow(),
				Links: []opdsLink{
					{Rel: "subsection", Href: opdsBase(apiKey, "collections"), Type: "application/atom+xml;profile=opds-catalog"},
				},
			},
			{
				Title:   "Reading Lists",
				ID:      "urn:folio:reading-lists",
				Updated: opdsNow(),
				Links: []opdsLink{
					{Rel: "subsection", Href: opdsBase(apiKey, "reading-lists"), Type: "application/atom+xml;profile=opds-catalog"},
				},
			},
		},
	}
	writeOPDS(w, feed)
}

func (h *Handler) OPDSLibraries(w http.ResponseWriter, r *http.Request) {
	apiKey := chi.URLParam(r, "apiKey")

	libs, err := h.svc.ListLibraries()
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}

	feed := opdsFeed{
		Xmlns:   "http://www.w3.org/2005/Atom",
		Title:   "Libraries",
		ID:      "urn:folio:libraries",
		Updated: opdsNow(),
		Links: []opdsLink{
			{Rel: "self", Href: opdsBase(apiKey, "libraries"), Type: "application/atom+xml;profile=opds-catalog"},
			{Rel: "start", Href: opdsBase(apiKey, ""), Type: "application/atom+xml;profile=opds-catalog"},
		},
	}

	for _, lib := range libs {
		feed.Entries = append(feed.Entries, opdsEntry{
			Title:   lib.Name,
			ID:      fmt.Sprintf("urn:folio:library:%d", lib.ID),
			Updated: lib.CreatedAt.Format(time.RFC3339),
			Links: []opdsLink{
				{Rel: "subsection", Href: fmt.Sprintf("/opds/%s/library/%d", apiKey, lib.ID), Type: "application/atom+xml;profile=opds-catalog"},
			},
		})
	}
	writeOPDS(w, feed)
}

func (h *Handler) OPDSLibrary(w http.ResponseWriter, r *http.Request) {
	apiKey := chi.URLParam(r, "apiKey")
	libID := uintParam(r, "libraryId")

	lib, err := h.svc.GetLibrary(libID)
	if err != nil {
		writeError(w, 404, "library not found")
		return
	}

	series, err := h.svc.ListSeries(libID, "")
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}

	feed := opdsFeed{
		Xmlns:   "http://www.w3.org/2005/Atom",
		Title:   lib.Name,
		ID:      fmt.Sprintf("urn:folio:library:%d", libID),
		Updated: opdsNow(),
		Links: []opdsLink{
			{Rel: "self", Href: fmt.Sprintf("/opds/%s/library/%d", apiKey, libID), Type: "application/atom+xml;profile=opds-catalog"},
			{Rel: "start", Href: opdsBase(apiKey, ""), Type: "application/atom+xml;profile=opds-catalog"},
		},
	}

	for _, s := range series {
		entry := opdsEntry{
			Title:   s.Name,
			ID:      fmt.Sprintf("urn:folio:series:%d", s.ID),
			Updated: s.CreatedAt.Format(time.RFC3339),
			Links: []opdsLink{
				{Rel: "subsection", Href: fmt.Sprintf("/opds/%s/series/%d", apiKey, s.ID), Type: "application/atom+xml;profile=opds-catalog"},
			},
		}
		if s.CoverPath != "" {
			entry.Links = append(entry.Links, opdsLink{
				Rel:  "http://opds-spec.org/image",
				Href: fmt.Sprintf("/api/image/series-cover?seriesId=%d", s.ID),
				Type: "image/jpeg",
			})
		}
		feed.Entries = append(feed.Entries, entry)
	}
	writeOPDS(w, feed)
}

func (h *Handler) OPDSSeries(w http.ResponseWriter, r *http.Request) {
	apiKey := chi.URLParam(r, "apiKey")
	seriesID := uintParam(r, "seriesId")

	ser, err := h.svc.GetSeries(seriesID)
	if err != nil {
		writeError(w, 404, "series not found")
		return
	}

	feed := opdsFeed{
		Xmlns:   "http://www.w3.org/2005/Atom",
		Title:   ser.Name,
		ID:      fmt.Sprintf("urn:folio:series:%d", seriesID),
		Updated: opdsNow(),
		Links: []opdsLink{
			{Rel: "self", Href: fmt.Sprintf("/opds/%s/series/%d", apiKey, seriesID), Type: "application/atom+xml;profile=opds-catalog"},
			{Rel: "start", Href: opdsBase(apiKey, ""), Type: "application/atom+xml;profile=opds-catalog"},
		},
	}

	for _, vol := range ser.Volumes {
		entry := opdsEntry{
			Title:   vol.Name,
			ID:      fmt.Sprintf("urn:folio:volume:%d", vol.ID),
			Updated: vol.CreatedAt.Format(time.RFC3339),
			Links: []opdsLink{
				{Rel: "subsection", Href: fmt.Sprintf("/opds/%s/series/%d/volume/%d", apiKey, seriesID, vol.ID), Type: "application/atom+xml;profile=opds-catalog"},
			},
		}
		if vol.CoverPath != "" {
			entry.Links = append(entry.Links, opdsLink{
				Rel:  "http://opds-spec.org/image",
				Href: fmt.Sprintf("/api/image/volume-cover?volumeId=%d", vol.ID),
				Type: "image/jpeg",
			})
		}
		feed.Entries = append(feed.Entries, entry)
	}
	writeOPDS(w, feed)
}

func (h *Handler) OPDSVolume(w http.ResponseWriter, r *http.Request) {
	apiKey := chi.URLParam(r, "apiKey")
	seriesID := uintParam(r, "seriesId")
	volID := uintParam(r, "volumeId")

	vol, err := h.svc.GetVolume(volID)
	if err != nil {
		writeError(w, 404, "volume not found")
		return
	}

	feed := opdsFeed{
		Xmlns:   "http://www.w3.org/2005/Atom",
		Title:   vol.Name,
		ID:      fmt.Sprintf("urn:folio:volume:%d", volID),
		Updated: opdsNow(),
		Links: []opdsLink{
			{Rel: "self", Href: fmt.Sprintf("/opds/%s/series/%d/volume/%d", apiKey, seriesID, volID), Type: "application/atom+xml;profile=opds-catalog"},
			{Rel: "start", Href: opdsBase(apiKey, ""), Type: "application/atom+xml;profile=opds-catalog"},
		},
	}

	for _, ch := range vol.Chapters {
		entry := opdsEntry{
			Title:   ch.Title,
			ID:      fmt.Sprintf("urn:folio:chapter:%d", ch.ID),
			Updated: ch.CreatedAt.Format(time.RFC3339),
			Links: []opdsLink{
				{Rel: "http://opds-spec.org/acquisition/open-access", Href: fmt.Sprintf("/opds/%s/series/%d/volume/%d/chapter/%d/download/%s", apiKey, seriesID, volID, ch.ID, ch.FilePath), Type: contentTypeFor(ch.FileType)},
			},
		}
		feed.Entries = append(feed.Entries, entry)
	}
	writeOPDS(w, feed)
}

func (h *Handler) OPDSChapter(w http.ResponseWriter, r *http.Request) {
	apiKey := chi.URLParam(r, "apiKey")
	seriesID := uintParam(r, "seriesId")
	volID := uintParam(r, "volumeId")
	chID := uintParam(r, "chapterId")

	ch, err := h.svc.GetChapter(chID)
	if err != nil {
		writeError(w, 404, "chapter not found")
		return
	}

	feed := opdsFeed{
		Xmlns:   "http://www.w3.org/2005/Atom",
		Title:   ch.Title,
		ID:      fmt.Sprintf("urn:folio:chapter:%d", chID),
		Updated: opdsNow(),
		Links: []opdsLink{
			{Rel: "self", Href: fmt.Sprintf("/opds/%s/series/%d/volume/%d/chapter/%d", apiKey, seriesID, volID, chID), Type: "application/atom+xml"},
			{Rel: "start", Href: opdsBase(apiKey, ""), Type: "application/atom+xml;profile=opds-catalog"},
		},
		Entries: []opdsEntry{
			{
				Title:   ch.Title,
				ID:      fmt.Sprintf("urn:folio:chapter:%d", chID),
				Updated: ch.CreatedAt.Format(time.RFC3339),
				Links: []opdsLink{
					{Rel: "http://opds-spec.org/acquisition/open-access", Href: fmt.Sprintf("/opds/%s/series/%d/volume/%d/chapter/%d/download/%s", apiKey, seriesID, volID, chID, ch.FilePath), Type: contentTypeFor(ch.FileType)},
				},
			},
		},
	}
	writeOPDS(w, feed)
}

func (h *Handler) OPDSDownload(w http.ResponseWriter, r *http.Request) {
	chID := uintParam(r, "chapterId")
	ch, err := h.svc.GetChapter(chID)
	if err != nil {
		writeError(w, 404, "chapter not found")
		return
	}

	w.Header().Set("Content-Type", contentTypeFor(ch.FileType))
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, ch.FilePath))
	http.ServeFile(w, r, ch.FilePath)
}

func (h *Handler) OPDSSearchDescriptor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/opensearchdescription+xml")
	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<OpenSearchDescription xmlns="http://a9.com/-/spec/opensearch/1.1/">
  <ShortName>Folio Search</ShortName>
  <Description>Search Folio library</Description>
  <Url type="application/atom+xml" template="/opds/search/{searchTerms}"/>
</OpenSearchDescription>`))
}

func contentTypeFor(fileType string) string {
	switch fileType {
	case "epub":
		return "application/epub+zip"
	case "cbz", "zip":
		return "application/x-cbz"
	case "cbr", "rar":
		return "application/x-cbr"
	default:
		return "application/octet-stream"
	}
}
