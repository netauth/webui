package http

import (
	"net/http"

	"github.com/flosch/pongo2/v4"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hashicorp/go-hclog"
	"github.com/netauth/netauth/pkg/netauth"
)

// New initializes the server with its default routers.
func New(l hclog.Logger) (*Server, error) {
	sbl, err := pongo2.NewSandboxedFilesystemLoader("theme/p2")
	if err != nil {
		return nil, err
	}

	l = l.Named("http")

	nc, err := netauth.NewWithLog(l)
	if err != nil {
		return nil, err
	}

	s := Server{
		l:     l,
		r:     chi.NewRouter(),
		n:     &http.Server{},
		c:     nc,
		tmpls: pongo2.NewSet("html", sbl),
	}

	s.tmpls.Debug = true

	s.r.Use(middleware.Logger)
	s.r.Use(middleware.Heartbeat("/healthz"))

	s.fileServer(s.r, "/static", http.Dir("theme/static"))
	s.r.Get("/", s.rootIndex)
	s.r.Get("/info/group/{name}", s.viewGroupInfo)

	return &s, nil
}

func (s *Server) rootIndex(w http.ResponseWriter, r *http.Request) {
	t, err := s.tmpls.FromCache("index.p2")
	if err != nil {
		s.templateErrorHandler(w, err)
		return
	}
	if err := t.ExecuteWriter(nil, w); err != nil {
		s.templateErrorHandler(w, err)
	}
}

// Serve binds, initializes the mux, and serves forever.
func (s *Server) Serve(bind string) error {
	s.l.Info("HTTP is starting")
	s.n.Addr = bind
	s.n.Handler = s.r
	return s.n.ListenAndServe()
}

func (s *Server) doTemplate(w http.ResponseWriter, r *http.Request, tmpl string, ctx pongo2.Context) {
	t, err := s.tmpls.FromCache(tmpl)
	if err != nil {
		s.templateErrorHandler(w, err)
		return
	}
	if err := t.ExecuteWriter(ctx, w); err != nil {
		s.templateErrorHandler(w, err)
	}
}
