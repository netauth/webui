package http

import (
	"context"
	"net/http"

	"github.com/flosch/pongo2/v4"
	"github.com/go-chi/chi/v5"
)

func (s *Server) viewEntityInfo(w http.ResponseWriter, r *http.Request) {
	entity, err := s.c.EntityInfo(context.Background(), chi.URLParam(r, "id"))
	if err != nil {
		s.doTemplate(w, r, "errors/netauth.p2", pongo2.Context{"error": err.Error()})
		return
	}

	ctx := pongo2.Context{
		"entity":   entity,
	}
	s.doTemplate(w, r, "view/entity_info.p2", ctx)
}
