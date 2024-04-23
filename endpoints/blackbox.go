package endpoints

import (
	"context"
	"net/http"

	"github.com/giantswarm/microerror"
	kitendpoint "github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
)

const (
	// BlackboxMethod is the HTTP method this endpoint is register for.
	BlackboxMethod = "GET"
	// BlackboxName identifies the endpoint. It is aligned to the package path.
	BlackboxName = "blackbox"
	// BlackboxPath is the HTTP request path this endpoint is registered for.
	BlackboxPath = "/blackbox"
)

type BlackboxConfig struct {
}

func NewBlackbox(config BlackboxConfig) (*Blackbox, error) {
	return &Blackbox{}, nil
}

type Blackbox struct {
}

func (b *Blackbox) Decoder() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		return nil, nil
	}
}

func (b *Blackbox) Encoder() kithttp.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("ok"))
		return microerror.Mask(err)
	}
}

func (b *Blackbox) Endpoint() kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, nil
	}
}

func (b *Blackbox) Method() string {
	return BlackboxMethod
}

func (b *Blackbox) Middlewares() []kitendpoint.Middleware {
	return []kitendpoint.Middleware{}
}

func (b *Blackbox) Name() string {
	return BlackboxName
}

func (b *Blackbox) Path() string {
	return BlackboxPath
}
