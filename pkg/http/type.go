package http

import (
	"net/http"

	"github.com/flosch/pongo2/v4"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/form/v4"
	"github.com/hashicorp/go-hclog"
	"github.com/netauth/netauth/pkg/netauth"
	"github.com/netauth/netauth/pkg/token"
)

// Server wraps up all the request routers and associated components
// that serve various parts of the nbuild stack.
type Server struct {
	l hclog.Logger
	r chi.Router

	n     *http.Server
	c     *netauth.Client
	t     token.Service
	f     *form.Decoder
	tmpls *pongo2.TemplateSet
}

type ctxToken struct{}
