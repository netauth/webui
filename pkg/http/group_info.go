package http

import (
	"context"
	"net/http"

	"github.com/flosch/pongo2/v4"
	"github.com/go-chi/chi/v5"
)

func (s *Server) viewGroupInfo(w http.ResponseWriter, r *http.Request) {
	grp, mgd, err := s.c.GroupInfo(context.Background(), chi.URLParam(r, "name"))
	if err != nil {
		s.doTemplate(w, r, "errors/netauth.p2", pongo2.Context{"error": err.Error()})
		return
	}

	ctx := pongo2.Context{
		"group":   grp,
		"managed": mgd,
	}
	s.doTemplate(w, r, "view/group_info.p2", ctx)
}
