package http

import (
	"net/http"

	"github.com/flosch/pongo2/v4"
)

func (s *Server) viewGroupSearch(w http.ResponseWriter, r *http.Request) {
	q := r.FormValue("q")
	ctx := pongo2.Context{"q": q}
	if q != "" {
		s.l.Debug("Performing group search", "query", q)
		res, err := s.c.GroupSearch(r.Context(), q)
		if err != nil {
			s.doTemplate(w, r, "errors/netauth.p2", pongo2.Context{"error": err})
			return
		}
		ctx.Update(pongo2.Context{"groups": res})
	}
	s.doTemplate(w, r, "view/group_search.p2", ctx)
}

func (s *Server) viewEntitySearch(w http.ResponseWriter, r *http.Request) {
	q := r.FormValue("q")
	ctx := pongo2.Context{"q": q}
	if q != "" {
		s.l.Debug("Performing entity search", "query", q)
		res, err := s.c.EntitySearch(r.Context(), q)
		if err != nil {
			s.doTemplate(w, r, "errors/netauth.p2", pongo2.Context{"error": err})
			return
		}
		ctx.Update(pongo2.Context{"entities": res})
	}
	s.doTemplate(w, r, "view/entity_search.p2", ctx)
}
