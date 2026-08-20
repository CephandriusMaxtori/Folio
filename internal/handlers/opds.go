package handlers

import (
	"net/http"
)

func (h *Handler) OPDSCatalog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:opds="http://opds-spec.org/2013/catalog">
  <title>Folio</title>
  <link rel="self" href="/opds/catalog" type="application/atom+xml;profile=opds-catalog"/>
  <link rel="start" href="/opds/catalog" type="application/atom+xml;profile=opds-catalog"/>
  <id>urn:folio:catalog</id>
  <updated>2026-01-01T00:00:00Z</updated>
</feed>`))
}

func (h *Handler) OPDSLibraries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>Libraries</title>
  <id>urn:folio:libraries</id>
  <updated>2026-01-01T00:00:00Z</updated>
</feed>`))
}

func (h *Handler) OPDSLibrary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>Library</title>
  <id>urn:folio:library</id>
  <updated>2026-01-01T00:00:00Z</updated>
</feed>`))
}

func (h *Handler) OPDSSeries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>Series</title>
  <id>urn:folio:series</id>
  <updated>2026-01-01T00:00:00Z</updated>
</feed>`))
}

func (h *Handler) OPDSVolume(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>Volume</title>
  <id>urn:folio:volume</id>
  <updated>2026-01-01T00:00:00Z</updated>
</feed>`))
}

func (h *Handler) OPDSChapter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>Chapter</title>
  <id>urn:folio:chapter</id>
  <updated>2026-01-01T00:00:00Z</updated>
</feed>`))
}

func (h *Handler) OPDSDownload(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", 501)
}

func (h *Handler) OPDSSearchDescriptor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<OpenSearchDescription xmlns="http://a9.com/-/spec/opensearch/1.1/">
  <ShortName>Folio Search</ShortName>
  <Description>Search Folio library</Description>
  <Url type="application/atom+xml" template="/opds/search/{searchTerms}"/>
</OpenSearchDescription>`))
}
