package http

import (
	"net/http"

	"github.com/flosch/pongo2/v4"
	"github.com/netauth/netauth/pkg/netauth"
)

func (s *Server) viewEntityCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Chuck the form back and bail.
		s.doTemplate(w, r, "view/entity_create.p2", nil)
		return
	}
	if err := r.ParseForm(); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	fdata := struct {
		ID     string
		Number int
		Secret string
	}{}

	if err := s.f.Decode(&fdata, r.Form); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	ctx := netauth.Authorize(r.Context(), s.getToken(r))
	if err := s.c.EntityCreate(ctx, fdata.ID, fdata.Secret, fdata.Number); err != nil {
		s.doTemplate(w, r, "errors/netauth.p2", pongo2.Context{"error": err.Error()})
		return
	}
	http.Redirect(w, r, "/info/entity/"+fdata.ID, http.StatusSeeOther)
}
