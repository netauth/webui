package http

import (
	"net/http"

	"github.com/flosch/pongo2/v4"
	"github.com/go-chi/chi/v5"
	"github.com/netauth/netauth/pkg/netauth"
)

func (s *Server) viewGroupInfo(w http.ResponseWriter, r *http.Request) {
	grp, mgd, err := s.c.GroupInfo(r.Context(), chi.URLParam(r, "name"))
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

func (s *Server) viewGroupCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Chuck the form back and bail.
		s.doTemplate(w, r, "view/group_create.p2", nil)
		return
	}
	if err := r.ParseForm(); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	fdata := struct {
		Name        string
		DisplayName string
		ManagedBy   string
		Number      int
	}{}

	if err := s.f.Decode(&fdata, r.Form); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	ctx := netauth.Authorize(r.Context(), s.getToken(r))
	if err := s.c.GroupCreate(ctx, fdata.Name, fdata.DisplayName, fdata.ManagedBy, fdata.Number); err != nil {
		s.doTemplate(w, r, "errors/netauth.p2", pongo2.Context{"error": err.Error()})
		return
	}
	http.Redirect(w, r, "/group/info/"+fdata.Name, http.StatusSeeOther)
}

func (s *Server) viewGroupEditForm(w http.ResponseWriter, r *http.Request) {
	grp, _, err := s.c.GroupInfo(r.Context(), chi.URLParam(r, "name"))
	if err != nil {
		s.doTemplate(w, r, "errors/netauth.p2", pongo2.Context{"error": err.Error()})
		return
	}

	ctx := pongo2.Context{
		"group": grp,
	}
	s.doTemplate(w, r, "view/group_edit.p2", ctx)
}

func (s *Server) viewGroupEditSubmit(w http.ResponseWriter, r *http.Request) {
	gname := chi.URLParam(r, "name")

	grp, _, err := s.c.GroupInfo(r.Context(), gname)
	if err != nil {
		s.doTemplate(w, r, "errors/netauth.p2", pongo2.Context{"error": err.Error()})
		return
	}

	if err := r.ParseForm(); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	if err := s.f.Decode(&grp, r.Form); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}
	// Prevent any funny business...
	grp.Name = &gname

	ctx := netauth.Authorize(r.Context(), s.getToken(r))
	if err := s.c.GroupUpdate(ctx, grp); err != nil {
		s.doTemplate(w, r, "errors/netauth.p2", pongo2.Context{"error": err.Error()})
		return
	}
	http.Redirect(w, r, "/group/info/"+gname, http.StatusSeeOther)
}
