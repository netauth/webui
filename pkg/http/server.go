package http

import (
	"net/http"

	"github.com/flosch/pongo2/v4"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/form/v4"
	"github.com/hashicorp/go-hclog"
	"github.com/netauth/netauth/pkg/netauth"
	"github.com/netauth/netauth/pkg/token"
	"github.com/netauth/netauth/pkg/token/keyprovider"
	"github.com/spf13/viper"
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
	nc.SetServiceName("webui")

	kp, err := keyprovider.New(viper.GetString("token.keyprovider"))
	if err != nil {
		l.Error("Error initializing token service", "error", err)
		return nil, err
	}

	tsvc, err := token.New(viper.GetString("token.backend"), kp)
	if err != nil {
		l.Error("Error initializing token service", "error", err)
		return nil, err
	}

	s := Server{
		l:     l,
		r:     chi.NewRouter(),
		n:     &http.Server{},
		c:     nc,
		f:     form.NewDecoder(),
		t:     tsvc,
		tmpls: pongo2.NewSet("html", sbl),
	}

	s.tmpls.Debug = true

	s.r.Use(middleware.Logger)
	s.r.Use(middleware.Heartbeat("/healthz"))
	s.r.Use(s.authMiddleware)

	s.fileServer(s.r, "/static", http.Dir("theme/static"))
	s.r.Get("/", s.rootIndex)
	s.r.Get("/login", s.viewLoginPage)
	s.r.Post("/login", s.viewLoginSubmit)
	s.r.Get("/logout", s.viewLogout)
	s.r.Get("/info/group/{name}", s.viewGroupInfo)
	s.r.Get("/info/entity/{id}", s.viewEntityInfo)

	return &s, nil
}

func (s *Server) rootIndex(w http.ResponseWriter, r *http.Request) {
	s.doTemplate(w, r, "index.p2", nil)
}

// Serve binds, initializes the mux, and serves forever.
func (s *Server) Serve(bind string) error {
	s.l.Info("HTTP is starting")
	s.n.Addr = bind
	s.n.Handler = s.r
	return s.n.ListenAndServe()
}
