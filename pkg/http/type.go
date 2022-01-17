package http

import (
	"net/http"

	"github.com/flosch/pongo2/v4"
	"github.com/go-chi/chi/v5"
	"github.com/hashicorp/go-hclog"
	"github.com/netauth/netauth/pkg/netauth"
)

// Server wraps up all the request routers and associated components
// that serve various parts of the nbuild stack.
type Server struct {
	l hclog.Logger
	r chi.Router

	n     *http.Server
	c     *netauth.Client
	tmpls *pongo2.TemplateSet
}
