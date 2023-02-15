package app

import (
	"OceanID/app/impls"
	"OceanID/app/ocean_id"
	"OceanID/config"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net/http"
)

func RegisterOceanIDHttpServer(mux *http.ServeMux, idPool oceanID.IDPool) {
	service := oceanID.Mount[*impls.OceanIDHttp](idPool, &impls.OceanIDHttp{})

	handlerWithErr := func(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if err := handler(w, r); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
		}
	}

	mux.HandleFunc("/generate_id", handlerWithErr(service.GenerateID))
	mux.HandleFunc("/bulk_generate_id", handlerWithErr(service.BulkGenerateID))
}

type HttpServer struct {
	enable  bool
	ctx     context.Context
	oceanID oceanID.IDPool
	*http.Server
}

func (h *HttpServer) GetName() string {
	return "http"
}

func (h *HttpServer) Setup() error {
	args, err := config.AssertArgs(h.ctx)
	if err != nil {
		return errors.Wrap(err, "app.app_http")
	}
	serverAddr := args.Get("HTTP_SERVER_ADDR").(string)
	if serverAddr == "" {
		return nil
	}

	mux := &http.ServeMux{}
	h.Server = &http.Server{
		Addr:    serverAddr,
		Handler: mux,
	}

	RegisterOceanIDHttpServer(mux, h.oceanID)

	h.enable = true

	return nil
}

func (h *HttpServer) Start() error {
	if !h.enable {
		return nil
	}
	return h.Server.ListenAndServe()
}

func (h *HttpServer) Shutdown() error {
	if !h.enable {
		return nil
	}
	return h.Server.Shutdown(h.ctx)
}

func NewHttpServer(ctx context.Context, oi oceanID.IDPool) Application {
	return &HttpServer{
		ctx:     ctx,
		oceanID: oi,
	}
}
