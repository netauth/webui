package http

import (
	"context"
	"net/http"

	"github.com/flosch/pongo2/v4"
	"github.com/spf13/viper"
)

func (s *Server) viewLoginPage(w http.ResponseWriter, r *http.Request) {
	s.doTemplate(w, r, "view/login.p2", nil)
}

func (s *Server) viewLoginSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	fdata := struct {
		ID     string
		Secret string
	}{}
	if err := s.f.Decode(&fdata, r.Form); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	t, err := s.c.AuthGetToken(r.Context(), fdata.ID, fdata.Secret)
	if err != nil {
		s.doTemplate(w, r, "errors/netauth.p2", pongo2.Context{"error": err.Error()})
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "netauth_tkn",
		Value:    t,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(viper.GetDuration("token.lifetime").Seconds()),
	})

	s.doTemplate(w, r, "index.p2", nil)
}

func (s *Server) viewLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "netauth_tkn", MaxAge: 0})
	s.doTemplate(w, r, "index.p2", nil)
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.l.Trace("In authMiddleware")
		c, err := r.Cookie("netauth_tkn")
		if err != nil {
			s.l.Trace("Exiting authMiddleware", "reason", "notoken", "error", err)
			next.ServeHTTP(w, r)
			return
		}
		s.l.Trace("Extracted cookie", "cookie", c)

		tkn, err := s.t.Validate(c.Value)
		if err != nil {
			s.l.Trace("Exiting authMiddleware", "reason", "invalidtoken", "error", err)
			next.ServeHTTP(w, r)
			return
		}
		s.l.Trace("Validated Token", "token", tkn)
		ctx := context.WithValue(r.Context(), ctxToken{}, tkn)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) getToken(r *http.Request) string {
	c, err := r.Cookie("netauth_tkn")
	if err != nil {
		return ""
	}
	return c.Value
}
